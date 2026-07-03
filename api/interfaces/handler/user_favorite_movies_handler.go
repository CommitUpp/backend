package handler

import (
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
)

type FavoriteMoviesHandler struct {
	favoriteMovieUsecase user.FavoriteMovieUsecase
}

func (h *FavoriteMoviesHandler) GetFavoriteMovies(c echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "認証情報が見つかりません",
		})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "認証情報が見つかりません",
		})
	}

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
					PosterUrl: m.PosterURL,
			})
	}

	return c.JSON(http.StatusOK, GetFavoriteMoviesResponse{
		Movies: res,
	})
}

func (h *FavoriteMoviesHandler) PutFavoriteMovies(c echo.Context) error {
	ctx := c.Request().Context()

	var req FavoriteMoviesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError{
			Message: "リクエストの形式が不正です",
		})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "認証情報が見つかりません",
		})
	}

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{
			Message: "認証情報が見つかりません",
		})
	}

	// リクエストのMovieIds(uuid)を文字列のスライスに変換
	movieIDs := make([]string, 0, len(req.MovieIds))
	for _, id := range req.MovieIds {
		movieIDs = append(movieIDs, id.String())
	}

	err := h.favoriteMovieUsecase.UpdateFavoriteMovies(ctx, userID, movieIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: "お気に入り映画の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, FavoriteMoviesResponse{Status: "success"})
}

func NewFavoriteMoviesHandler(u user.FavoriteMovieUsecase) *FavoriteMoviesHandler {
	return &FavoriteMoviesHandler{
		favoriteMovieUsecase: u,
	}
}
