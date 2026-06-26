package repository

import (
	"context"
)

type MovieDetailRepository interface {
	IsGroupMember(ctx context.Context, groupID string, userID string) (bool, error)
	GetMovieDetail(ctx context.Context, movieID string, groupID string, userID string) (*MovieDetail, error)
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
	UserName  string
	AvatarURL string
}
