package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/CommitUpp/backend/api/application/usecase/group"
	"github.com/labstack/echo/v4"
)

type GroupWatchedMovieHandler struct {
	groupWatchedMovieUsecase group.GroupWatchedMovieUsecase
}

func NewGroupWatchedMovieHandler(u group.GroupWatchedMovieUsecase) *GroupWatchedMovieHandler {
	return &GroupWatchedMovieHandler{
		groupWatchedMovieUsecase: u,
	}
}

func (h *GroupWatchedMovieHandler) GetGroupWatchedMovies(c echo.Context, groupId string) error {
	ctx := c.Request().Context()

	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "ユーザーIDの解析に失敗しました"})
	}

	if groupId == "" {
		return c.JSON(http.StatusBadRequest, BadRequestError{Message: "グループIDが指定されていません"})
	}

	accessToken, ok := c.Get("access_token").(string)
	if !ok || accessToken == "" {
		return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
	}

	watchedMovies, err := h.groupWatchedMovieUsecase.GetWatchedMovies(ctx, userIDStr, groupId, accessToken)
	if err != nil {
		switch {
		case errors.Is(err, group.ErrUserIDRequired), errors.Is(err, group.ErrAccessTokenRequired):
			return c.JSON(http.StatusUnauthorized, UnauthorizedError{Message: "認証情報が見つかりません"})
		case errors.Is(err, group.ErrGroupIDRequired):
			return c.JSON(http.StatusBadRequest, BadRequestError{Message: "グループIDが指定されていません"})
		case errors.Is(err, group.ErrForbidden):
			return c.JSON(http.StatusForbidden, ForbiddenError{Message: "このグループの閲覧権限がありません"})
		default:
			log.Printf("failed to get group watched movies: user_id=%s group_id=%s err=%v", userIDStr, groupId, err)
			return c.JSON(http.StatusInternalServerError, InternalServerError{Message: "グループ内の視聴済み映画一覧の取得に失敗しました"})
		}
	}

	movieIndex := make(map[string]int)
	movies := make([]GetGroupWatchedMoviesMovie, 0)

	for _, row := range watchedMovies {
		index, exists := movieIndex[row.MovieID]
		if !exists {
			movies = append(movies, GetGroupWatchedMoviesMovie{
				MovieId:       row.MovieID,
				Title:         row.Title,
				TrailerUrl:     row.TrailerURL,
				WatchedMember: []WatchedMember{},
			})
			index = len(movies) - 1
			movieIndex[row.MovieID] = index
		}

		movies[index].WatchedMember = append(movies[index].WatchedMember, WatchedMember{
			UserId:    row.UserID,
			AvatarUrl: row.AvatarURL,
		})
	}

	return c.JSON(http.StatusOK, GetGroupWatchedMoviesResponse{
		GroupId: groupId,
		Movies:  movies,
	})
}
