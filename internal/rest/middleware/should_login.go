package middleware

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/internal/repository"
	"github.com/usual2970/retell/pkg/config"
	pkgJwt "github.com/usual2970/retell/pkg/jwt"
)

var authorizationKey = "Authorization"
var bearerWord = "Bearer"

func ShouldLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := getAuthToken(c)
		if err != nil {
			return next(c)
		}

		repo := repository.NewAccessTokenRepository()

		_, err = repo.GetAccessToken(c.Request().Context(), token.Raw)
		if err != nil {
			return next(c)
		}

		ctx := pkgJwt.WithContext(c.Request().Context(), token)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func getAuthToken(c echo.Context) (*jwt.Token, error) {
	auths := strings.SplitN(c.Request().Header.Get(authorizationKey), " ", 2)
	if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
		return nil, jwt.ErrTokenMalformed
	}
	jwtToken := auths[1]

	conf := config.GetConfig().Auth

	signingKey := conf.JwtSecret
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(signingKey), nil
	}
	token, err := jwt.Parse(jwtToken, keyFunc)
	if err != nil {
		return nil, constant.ErrNotLogin
	}

	if !token.Valid {
		return nil, constant.ErrNotLogin
	}
	return token, nil
}
