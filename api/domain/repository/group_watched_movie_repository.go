package repository

import "context"

type GroupWatchedMovieRepository interface {
	GetWatchedMovies(ctx context.Context, groupID string, excludeUserID string) ([]GroupWatchedMovieRow, error)
}

type GroupWatchedMovieRow struct {
	GroupID   string
	MovieID   string
	Title     string
	PosterURL string
	UserID    string
	AvatarURL string
}
