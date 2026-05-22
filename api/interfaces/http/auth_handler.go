package http

import (
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	apiHttp "github.com/CommitUpp/backend/api/interfaces/http"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	au user.AuthUseCase
}

func NewAuthHandler(au user.AuthUseCase) *AuthHandler {
	return &AuthHandler{au: au}
}

func (h *AuthHandler) LoginCallback(c echo.Context) error {
	ctx := c.Request().Context()

	var req apiHttp.LoginCallbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
}
