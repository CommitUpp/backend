package handler

import (
	"log"
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/movie"
	"github.com/labstack/echo/v4"
)

type MovieStatusHandler struct {
	movieStatusUsecase movie.MovieStatusUsecase
}

func NewMovieStatusHandler(u movie.MovieStatusUsecase) *MovieStatusHandler {
	return &MovieStatusHandler{
		movieStatusUsecase: u,
	}
}

func (h *MovieStatusHandler) WatchStatus(c echo.Context) error {
	ctx := c.Request().Context()

	//	リクエストボディのバインド
	var req MovieStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError{Message: "リクエストの形式が不正です"})
	}

	// 空チェック
	if req.MovieId.String() == "" || req.Status == "" {
		return c.JSON(http.StatusBadRequest, BadRequestError{Message: "必須項目が不足しています"})
	}

	//	認証ミドルウェアからユーザーIDを取得
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "ユーザーIDの解析に失敗しました"})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	//	ユースケースの実行
	err := h.movieStatusUsecase.Execute(ctx, req.MovieId.String(), userIDStr, string(req.Status), accessToken)
	if err != nil {
		// Enum値が不正な場合は400、その他DBエラー等は500を返却
		if err.Error() == "invalid status value" {
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: err.Error()})
		}
		log.Printf("failed to update movie status: movie_id=%s user_id=%s status=%s err=%v", req.MovieId.String(), userIDStr, req.Status, err)
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "ステータスの更新に失敗しました"})
	}

	return c.JSON(http.StatusOK, MovieStatusResponse{Status: "success"})
}

func (h *MovieStatusHandler) GetMovieStatus(c echo.Context, params GetMovieStatusParams) error {
	ctx := c.Request().Context()

	var status *string
	if params.Status != nil {
		s := string(*params.Status)
		status = &s
	}

	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "ユーザーIDの解析に失敗しました"})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	movieStatuses, err := h.movieStatusUsecase.GetStatuses(ctx, userIDStr, status, accessToken)
	if err != nil {
		if err.Error() == "invalid status filter value" {
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: err.Error()})
		}
		if err.Error() == "user ID is required" || err.Error() == "access token is required" {
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
		}
		log.Printf("failed to get movie statuses: user_id=%s status=%v err=%v", userIDStr, status, err)
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "映画一覧の取得に失敗しました"})
	}

	movies := make([]GetMovieStatus, 0, len(movieStatuses))
	for _, movieStatus := range movieStatuses {
		movies = append(movies, GetMovieStatus{
			MovieId:     movieStatus.MovieID,
			TmdbId:      movieStatus.TMDBID,
			Title:       movieStatus.Title,
			PosterUrl:   movieStatus.PosterURL,
			TrailerUrl:  movieStatus.TrailerURL,
			Overview:    movieStatus.Overview,
			ReleaseDate: movieStatus.ReleaseDate,
			UpdatedAt:   movieStatus.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, GetMovieStatusResponse{Movies: movies})
}
