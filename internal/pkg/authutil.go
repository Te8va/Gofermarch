package pkg

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Te8va/Gofermarch/pkg/jwt"
)

func ExtractLoginFromRequest(r *http.Request, jwtKey string) (string, error) {
	tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := jwt.ParseJWT(tokenStr, []byte(jwtKey))
	if err != nil {
		return "", err
	}

	login := claims.Subject
	if login == "" {
		return "", errors.New("unauthorized user")
	}

	return login, nil
}
