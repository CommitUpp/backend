package handler

import (
	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
)

type FavoriteMoviesHandler struct {
	favoriteMovieUsecase user.FavoriteMovieUsecase
}

func (h *FavoriteMoviesHandler) GetFavoriteMovies(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}

func (h *FavoriteMoviesHandler) PutFavoriteMovies(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}

func NewFavoriteMoviesHandler(u user.FavoriteMovieUsecase) *FavoriteMoviesHandler {
	return &FavoriteMoviesHandler{
		favoriteMovieUsecase: u,
	}
}
