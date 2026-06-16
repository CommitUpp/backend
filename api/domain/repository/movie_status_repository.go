package repository

import "context"

type MovieStatusRepository interface {
	UpsertWatchStatus(ctx context.Context, movieID string, userID string, status string, accessToken string) error
}
