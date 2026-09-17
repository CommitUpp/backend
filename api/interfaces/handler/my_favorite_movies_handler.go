package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/user"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type MyFavoriteMoviesHandler struct {
	favoriteMovieUsecase user.FavoriteMovieUsecase
}

func (h *MyFavoriteMoviesHandler) GetMyFavoriteMovies(c echo.Context) error {
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
			MovieId:    m.MovieID,
			TmdbId:     m.TMDBID,
			Title:      m.Title,
			TrailerUrl: m.TrailerURL,
		})
	}

	return c.JSON(http.StatusOK, GetFavoriteMoviesResponse{
		Movies: res,
	})
}

func (h *MyFavoriteMoviesHandler) PutMyFavoriteMovies(c echo.Context) error {
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
		if errors.Is(err, user.ErrFavoriteMovieLimitExceeded) {
			return c.JSON(http.StatusConflict, ConflictError{
				Message: "お気に入り映画は3件までです",
			})
		}

		return c.JSON(http.StatusInternalServerError, InternalServerError{
			Message: "お気に入り映画の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, FavoriteMoviesResponse{Status: "success"})
}

func (h *MyFavoriteMoviesHandler) AddMyFavoriteMovie(c echo.Context, movieId openapi_types.UUID) error {
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

	err := h.favoriteMovieUsecase.AddFavoriteMovie(ctx, userID, movieId.String())
	if err != nil {
		switch {
		case errors.Is(err, user.ErrMovieIDRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "映画IDが指定されていません"})
		case errors.Is(err, user.ErrMovieNotFound):
			return c.JSON(http.StatusNotFound, NotFoundError{Message: "映画が見つかりません"})
		case errors.Is(err, user.ErrFavoriteMovieAlreadyExists):
			return c.JSON(http.StatusConflict, ConflictError{Message: "既にお気に入り登録済みです"})
		case errors.Is(err, user.ErrFavoriteMovieLimitExceeded):
			return c.JSON(http.StatusConflict, ConflictError{Message: "お気に入り映画は3件までです"})
		default:
			log.Printf("failed to add favorite movie: user_id=%s movie_id=%s err=%v", userID, movieId.String(), err)
			return c.JSON(http.StatusInternalServerError, InternalServerError{
				Message: "お気に入り映画の追加に失敗しました",
			})
		}
	}

	return c.JSON(http.StatusCreated, FavoriteMoviesResponse{Status: "success"})
}

func (h *MyFavoriteMoviesHandler) DeleteMyFavoriteMovie(c echo.Context, movieId openapi_types.UUID) error {
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

	err := h.favoriteMovieUsecase.DeleteFavoriteMovie(ctx, userID, movieId.String())
	if err != nil {
		switch {
		case errors.Is(err, user.ErrMovieIDRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "映画IDが指定されていません"})
		case errors.Is(err, user.ErrFavoriteMovieNotFound):
			return c.JSON(http.StatusNotFound, NotFoundError{Message: "お気に入り登録が見つかりません"})
		default:
			log.Printf("failed to delete favorite movie: user_id=%s movie_id=%s err=%v", userID, movieId.String(), err)
			return c.JSON(http.StatusInternalServerError, InternalServerError{
				Message: "お気に入り映画の解除に失敗しました",
			})
		}
	}

	return c.JSON(http.StatusOK, FavoriteMoviesResponse{Status: "success"})
}

func NewMyFavoriteMoviesHandler(u user.FavoriteMovieUsecase) *MyFavoriteMoviesHandler {
	return &MyFavoriteMoviesHandler{
		favoriteMovieUsecase: u,
	}
}
