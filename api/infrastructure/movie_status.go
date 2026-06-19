package infrastructure

import (
	"context"
	"strings"

	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/supabase-community/postgrest-go"
)

type MovieStatusRepository struct {
	baseURL string
	apiKey  string
}

func NewMovieStatusRepository(baseURL string, apiKey string) repository.MovieStatusRepository {
	return &MovieStatusRepository{
		baseURL: strings.TrimRight(baseURL, "/") + "/rest/v1",
		apiKey:  apiKey,
	}
}

func (r *MovieStatusRepository) WatchStatus(ctx context.Context, movieID string, userID string, status string, accessToken string) error {
	// データを作成
	row := map[string]interface{}{
		"user_id":  userID,
		"movie_id": movieID,
		"status":   status,
	}

	client := postgrest.NewClient(
		r.baseURL,
		"public",
		map[string]string{
			"apikey":        r.apiKey,
			"Authorization": "Bearer " + accessToken,
		},
	)
	if client.ClientError != nil {
		return client.ClientError
	}

	// watch_statusesに登録
	_, _, err := client.
		From("watch_statuses").
		Upsert(row, "user_id,movie_id", "minimal", "merge").
		ExecuteWithContext(ctx)

	if err != nil {
		return err
	}

	return nil
}
