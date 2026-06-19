package router

import (
	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/CommitUpp/backend/api/interfaces/handler"
	"github.com/labstack/echo/v4"
)

type RouterConfig struct {
	AuthUsecase user.AuthUsecase
	Server      *handler.Server
}

func NewRouter(cfg RouterConfig) *echo.Echo {
	e := echo.New()

	e.Use(handler.CORSMiddleware())

	authMiddleware := handler.SupabaseAuthMiddleware(cfg.AuthUsecase)
	handler.RegisterHandlersWithOptions(e, cfg.Server, handler.RegisterHandlersOptions{
		BaseURL: "/api/v1",
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"loginCallback": {authMiddleware},
			"logout":        {authMiddleware},
			"watchStatus":   {authMiddleware},
			"createGroup":   {authMiddleware},
		},
	})

	return e
}
