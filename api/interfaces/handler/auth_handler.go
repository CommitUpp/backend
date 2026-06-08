package handler

import (
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase user.AuthUsecase
}

func NewAuthHandler(au user.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: au,
	}
}

func (h *AuthHandler) LoginCallback(c echo.Context) error {
	ctx := c.Request().Context()

	accessToken := bearerToken(c.Request().Header.Get("Authorization"))
	if accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "missing token",
		})
	}

	_, err := h.authUsecase.LoginCallback(ctx, accessToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: err.Error(),
		})
	}

	res := LoginCallbackResponse{
		Status: "success",
	}

	return c.JSON(http.StatusOK, res)
}
