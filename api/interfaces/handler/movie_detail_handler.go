package handler

import (
	"net/http"

	movieusecase "github.com/CommitUpp/backend/api/application/usecase/movie"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type MovieDetailHandler struct {
	usecase *movieusecase.MovieDetailUsecase
}

func NewMovieDetailHandler(
	usecase *movieusecase.MovieDetailUsecase,
) *MovieDetailHandler {
	return &MovieDetailHandler{
		usecase: usecase,
	}
}

func (h *MovieDetailHandler) GetMovieDetail(
	c echo.Context,
	movieId openapi_types.UUID,
	params GetMovieDetailParams,
) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "認証情報が見つかりません",
		})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: "ユーザーIDの解析に失敗しました",
		})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "認証情報が見つかりません",
		})
	}

	movieDetail, err := h.usecase.GetMovieDetail(
		c.Request().Context(),
		movieId.String(),
		params.GroupId.String(),
		userIDStr,
		accessToken,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, movieDetail)
}
