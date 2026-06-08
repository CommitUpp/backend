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

	var req LoginCallbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError{
			Message: "invalid request body",
		})
	}

	_, err := h.authUsecase.LoginCallback(ctx, req.AccessToken)
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
