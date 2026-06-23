package repository

import (
	"context"
)

type MovieRepository interface {
	GetMovies(ctx context.Context, keyword string) ([]Movie, error)
}

type Movie struct {
	MovieID     string
	Title       string
	PosterURL   string
}
