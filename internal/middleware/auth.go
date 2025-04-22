package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Te8va/Gofermarch/internal/domain"
	"github.com/Te8va/Gofermarch/pkg/jwt"
	"github.com/Te8va/Gofermarch/pkg/logger"
	"go.uber.org/zap"
)

func Auth(jwtKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			if cookie, err := r.Cookie("auth_token"); err == nil {
				token = cookie.Value
			} else {
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			if token == "" {
				logger.Logger().Warn("Missing auth token (cookie or header)")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ParseJWT(token, []byte(jwtKey))
			if err != nil {
				logger.Logger().Warn("Failed to parse JWT", zap.String("token", token), zap.Error(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			login := claims.Subject
			if login == "" {
				logger.Logger().Warn("Empty login (subject) in token")
				http.Error(w, "Invalid token subject", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), domain.UserCtxKey, login)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
