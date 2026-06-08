package supabase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CommitUpp/backend/auth/domain/repository"
)

var ErrInvalidToken = errors.New("invalid token")

type tokenVerifierRepositoryImpl struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewTokenVerifierRepository(baseURL string, apiKey string) repository.TokenVerifierRepository {
	return &tokenVerifierRepositoryImpl{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *tokenVerifierRepositoryImpl) Verify(ctx context.Context, accessToken string) (*repository.VerifiedToken, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/auth/v1/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	if r.apiKey != "" {
		req.Header.Set("apikey", r.apiKey)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, ErrInvalidToken
	}

	var user struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, ErrInvalidToken
	}

	ttl, err := tokenTTL(accessToken, time.Now())
	if err != nil {
		return nil, err
	}

	return &repository.VerifiedToken{
		UserID: user.ID,
		TTL:    ttl,
	}, nil
}

func tokenTTL(accessToken string, now time.Time) (int64, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return 0, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, err
	}

	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, err
	}
	if claims.Exp <= 0 {
		return 0, fmt.Errorf("%w: exp claim is missing", ErrInvalidToken)
	}

	ttl := claims.Exp - now.Unix()
	if ttl <= 0 {
		return 0, ErrInvalidToken
	}

	return ttl, nil
}
