package repository

import (
	"context"
	"time"
)

type MovieStatusRepository interface {
	WatchStatus(ctx context.Context, movieID string, userID string, status string, accessToken string) error
	GetWatchStatuses(ctx context.Context, userID string, status *string, accessToken string) ([]MovieStatus, error)
}

type MovieStatus struct {
	MovieID     string    `json:"movie_id"`
	TMDBID      string    `json:"tmdb_id"`
	Title       string    `json:"title"`
	PosterURL   string    `json:"poster_url"`
	TrailerURL  string    `json:"trailer_url"`
	Overview    string    `json:"overview"`
	ReleaseDate string    `json:"release_date"`
	UpdatedAt   time.Time `json:"updated_at"`
}
