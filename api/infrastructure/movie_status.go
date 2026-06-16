package infrastructure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
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
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (r *MovieStatusRepository) UpsertWatchStatus(ctx context.Context, movieID string, userID string, status string, accessToken string) error {
	claims, claimsErr := jwtClaims(accessToken)
	if claimsErr != nil {
		log.Printf("failed to decode access token claims for watch status upsert: err=%v", claimsErr)
	} else {
		log.Printf("watch status upsert auth claims: sub=%s role=%s matches_user_id=%t", claims.Sub, claims.Role, claims.Sub == userID)
	}

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

type accessTokenClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
}

func jwtClaims(accessToken string) (*accessTokenClaims, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid jwt format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims accessTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
