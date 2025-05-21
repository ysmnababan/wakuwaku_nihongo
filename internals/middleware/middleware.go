package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"wakuwaku_nihongo/config"

	"wakuwaku_nihongo/internals/pkg/redisutil"

	"wakuwaku_nihongo/internals/utils/validator"
)

var redis *redisutil.Redis

func Init(e *echo.Echo, fRedis *redisutil.Redis) {
	redis = fRedis

	name := fmt.Sprintf("%s-%s", config.Get().App.Name, config.Env())

	e.Validator = validator.NewCustomValidator()

	e.Use(
		Recover,
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}),
		echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format: fmt.Sprintf(`{"time":"${time_custom}","remote_ip": "${remote_ip}",`+
				`"host":"${host}","method":"${method}","uri":"${uri}","status":${status},`+
				`"error":"${error}","user_agent":"${user_agent}","latency":${latency},"latency_human":"${latency_human}"`+
				`,"name":"%s"}`+"\n", name),
			CustomTimeFormat: time.RFC3339,
			Output:           os.Stdout,
		}),
	)
}

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func(c echo.Context) {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				log.Error().Any("error", r).RawJSON("stackTrace", stackTrace).Send()

				c.JSON(500, map[string]any{
					"message": "something went wrong",
				})
			}
		}(c)

		return next(c)
	}
}
