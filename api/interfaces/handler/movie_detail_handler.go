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
) error {
	movieDetail, err := h.usecase.GetMovieDetail(
		c.Request().Context(),
		movieId.String(),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, movieDetail)
}
