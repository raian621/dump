package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(tf *TokenFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessTokenStr, found := strings.CutPrefix(c.Request().Header.Get("Authorization"), "Bearer: ")
			if !found {
				return c.String(http.StatusUnauthorized, "No access token provided")
			}

			accessToken, err := tf.parseToken(accessTokenStr, &AccessTokenClaims{})
			if err != nil {
				c.Logger().Error("error authenticating user:", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token")
			}

			c.Set("user_id", accessToken.Claims.(*AccessTokenClaims).UserId)

			return next(c)
		}
	}
}
