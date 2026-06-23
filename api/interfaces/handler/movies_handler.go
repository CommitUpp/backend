package handler

import (
	"log"
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/movie"
	"github.com/labstack/echo/v4"
)

type MoviesHandler struct {
	moviesUsecase movie.MoviesUsecase
}

func NewMoviesHandler(u movie.MoviesUsecase) *MoviesHandler {
	return &MoviesHandler{
		moviesUsecase: u,
	}
}

func (h *MoviesHandler) GetMovies(
	c echo.Context,
	params GetMoviesParams,
) error {
	ctx := c.Request().Context()

	keyword := ""
	if params.Keyword != nil {
		keyword = *params.Keyword
	}

	movies, err := h.moviesUsecase.GetMovies(ctx, keyword)
	if err != nil {
		log.Printf("failed to get movies: keyword=%s err=%v", keyword, err)

		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: "映画一覧の取得に失敗しました",
		})
	}

	res := make([]Movie, 0, len(movies))
	for _, m := range movies {
		res = append(res, Movie{
			MovieId:   m.MovieID,
			Title:     m.Title,
			PosterUrl: m.PosterURL,
		})
	}

	return c.JSON(http.StatusOK, MoviesResponse{
		Movies: res,
	})
}
