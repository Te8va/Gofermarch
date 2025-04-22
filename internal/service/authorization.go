package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
	"github.com/Te8va/Gofermarch/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Authorization struct {
	repo   domain.AuthorizationRepository
	JWTKey string
}

func NewAuthorization(repo domain.AuthorizationRepository, jwtKey string) *Authorization {
	return &Authorization{repo: repo, JWTKey: jwtKey}
}

func (s *Authorization) Register(ctx context.Context, login, password string) (string, error) {
	_, err := s.repo.GetUserByLogin(ctx, login)
	if err == nil {
		return "", appErrors.ErrAlreadyRegistered
	}
	if !errors.Is(err, appErrors.ErrUserNotFound) {
		return "", fmt.Errorf("service.Register: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("service.Register: %w", err)
	}

	tokenStr, err := jwt.CreateJWT(login, []byte(s.JWTKey), time.Now().Add(24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("service.Register: %w", err)
	}

	user := domain.User{
		Login:    login,
		Password: string(passwordHash),
		Token:    tokenStr,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("service.Register: %w", err)
	}

	return tokenStr, nil
}

func (s *Authorization) Authenticate(ctx context.Context, login, password string) (string, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", appErrors.ErrWrongPassword
	}

	return user.Token, nil
}
