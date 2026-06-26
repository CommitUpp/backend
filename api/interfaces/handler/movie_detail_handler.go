package handler

import (
	"net/http"

	movieusecase "github.com/CommitUpp/backend/api/application/usecase/movie"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/labstack/echo/v4"
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

	movieDetail, err := h.usecase.GetMovieDetail(
		c.Request().Context(),
		movieId.String(),
		params.GroupId.String(),
		userIDStr,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, movieDetail)
}
