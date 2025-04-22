package domain

import "context"

type AuthorizationService interface {
	Register(ctx context.Context, login, password string) (string, error)
	Authenticate(ctx context.Context, login, password string) (string, error)
}

type AuthorizationRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByLogin(ctx context.Context, login string) (*User, error)
}
