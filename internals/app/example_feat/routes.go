package example_feat

import (
	"github.com/labstack/echo/v4"
	"wakuwaku_nihongo/internals/middleware"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.GetUsers, middleware.Authentication)
	g.POST("", h.CreateUser)
	g.POST("/auth", h.Login)
}
