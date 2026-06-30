package repository

import (
	"context"
	"time"
)

type UserMovieStatusRepository interface {
	WatchStatus(ctx context.Context, movieID string, userID string, status string, accessToken string) error
	GetWatchStatuses(ctx context.Context, userID string, status *string, accessToken string) ([]UserMovieStatus, error)
}

type UserMovieStatus struct {
	MovieID     string
	TMDBID      string
	Title       string
	PosterURL   string
	TrailerURL  string
	Overview    string
	ReleaseDate string
	UpdatedAt   time.Time
}
