package handler

import (
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
)

type UserFavoriteMoviesHandler struct {
	favoriteMovieUsecase user.FavoriteMovieUsecase
}

func (h *UserFavoriteMoviesHandler) GetUserFavoriteMovies(
	c echo.Context,
	userID string,
) error {
	ctx := c.Request().Context()

	movies, err := h.favoriteMovieUsecase.GetFavoriteMovies(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: "お気に入り映画の取得に失敗しました",
		})
	}

	var res []Movie
	for _, m := range movies {
		res = append(res, Movie{
			MovieId:   m.MovieID,
			TmdbId:    m.TMDBID,
			Title:     m.Title,
			TrailerUrl: m.TrailerURL,
		})
	}

	return c.JSON(http.StatusOK, GetFavoriteMoviesResponse{
		Movies: res,
	})
}

func NewUserFavoriteMoviesHandler(u user.FavoriteMovieUsecase) *UserFavoriteMoviesHandler {
	return &UserFavoriteMoviesHandler{
		favoriteMovieUsecase: u,
	}
}
