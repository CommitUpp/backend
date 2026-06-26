package repository

import (
	"context"
)

type MovieDetailRepository interface {
	GetMovieDetail(ctx context.Context, movieID string) (*MovieDetail, error)
}

type MovieDetail struct {
	MovieID           string
	TMDBID            string
	Title             string
	PosterURL         string
	TrailerURL        string
	Overview          string
	ReleaseDate       string
	WatchedUser       []WatchedUser
	StreamingServices []string
}

type WatchedUser struct {
	UserID    string
	AvatarURL string
}
