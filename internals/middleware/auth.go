package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"wakuwaku_nihongo/config"
	res "wakuwaku_nihongo/internals/utils/response"

	"github.com/golang-jwt/jwt"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	jwtKey := config.Get().JWT.Key
	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		if authToken == "" {
			return res.ErrorWrap(res.ErrUnauthorized, fmt.Errorf("invalid token")).Send(c)
		}

		splitToken := strings.Split(authToken, "Bearer ")

		if len(splitToken) < 2 {
			return res.ErrorMessageFrom(res.ErrUnauthorized, fmt.Errorf("invalid token")).Send(c)
		}

		token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
			}

			return []byte(jwtKey), nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "Token is expired") {
				return res.ErrorWrap(res.ErrExpiredAccessToken, err).Send(c)
			}
			return res.ErrorWrap(res.ErrUnauthorized, err).Send(c)
		}

		if !token.Valid {
			return res.ErrorWrap(res.ErrUnauthorized, err).Send(c)
		}

		var user_id string
		destructName := token.Claims.(jwt.MapClaims)["user_id"]
		if destructName != nil {
			user_id = destructName.(string)
		} else {
			user_id = ""
		}

		c.Set("user_id", user_id)

		var email string
		destructName = token.Claims.(jwt.MapClaims)["email"]
		if destructName != nil {
			email = destructName.(string)
		} else {
			email = ""
		}

		c.Set("user_id", user_id)
		c.Set("email", email)
		return next(c)
	}
}
