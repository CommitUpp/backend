package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/CommitUpp/backend/api/domain/repository"
)

type GroupRepository struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewGroupRepository(baseURL string, apiKey string) repository.GroupRepository {
	return &GroupRepository{
		baseURL: strings.TrimRight(baseURL, "/") + "/rest/v1",
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *GroupRepository) CreateGroupWithOwner(ctx context.Context, input repository.CreateGroupWithOwnerInput) (repository.CreatedGroup, error) {
	body := map[string]interface{}{
		"name":         input.Name,
		"monthly_goal": input.MonthlyGoal,
	}

	return r.callGroupRPC(ctx, "create_group", body, input.AccessToken)
}

func (r *GroupRepository) JoinGroup(ctx context.Context, input repository.JoinGroupInput) (repository.CreatedGroup, error) {
	body := map[string]interface{}{
		"group_id": input.GroupID,
	}

	return r.callGroupRPC(ctx, "join_group", body, input.AccessToken)
}

func (r *GroupRepository) callGroupRPC(ctx context.Context, name string, body interface{}, accessToken string) (repository.CreatedGroup, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return repository.CreatedGroup{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/rpc/"+name, bytes.NewReader(reqBody))
	if err != nil {
		return repository.CreatedGroup{}, err
	}

	req.Header.Set("apikey", r.apiKey)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := r.client.Do(req)
	if err != nil {
		return repository.CreatedGroup{}, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return repository.CreatedGroup{}, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return repository.CreatedGroup{}, groupRPCError(resBody)
	}

	return decodeCreatedGroup(resBody)
}

func groupRPCError(body []byte) error {
	var errRes struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
		Hint    string `json:"hint"`
	}
	if err := json.Unmarshal(body, &errRes); err != nil {
		return fmt.Errorf("group rpc failed: %s", string(body))
	}

	message := strings.ToLower(errRes.Message)
	switch {
	case errRes.Code == "23505":
		return repository.ErrGroupMemberAlreadyExists
	case errRes.Code == "23503", strings.Contains(message, "not found"):
		return repository.ErrGroupNotFound
	default:
		return fmt.Errorf("group rpc failed: code=%s message=%s", errRes.Code, errRes.Message)
	}
}

func decodeCreatedGroup(body []byte) (repository.CreatedGroup, error) {
	var rows []groupRPCRow
	if err := json.Unmarshal(body, &rows); err == nil {
		if len(rows) == 0 {
			return repository.CreatedGroup{}, repository.ErrGroupNotFound
		}
		return rows[0].createdGroup()
	}

	var row groupRPCRow
	if err := json.Unmarshal(body, &row); err != nil {
		return repository.CreatedGroup{}, err
	}
	return row.createdGroup()
}

type groupRPCRow struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	MonthlyGoal      int       `json:"monthlyGoal"`
	MonthlyGoalSnake int       `json:"monthly_goal"`
	CreatedAt        time.Time `json:"createdAt"`
	CreatedAtSnake   time.Time `json:"created_at"`
}

func (r groupRPCRow) createdGroup() (repository.CreatedGroup, error) {
	monthlyGoal := r.MonthlyGoal
	if monthlyGoal == 0 {
		monthlyGoal = r.MonthlyGoalSnake
	}

	createdAt := r.CreatedAt
	if createdAt.IsZero() {
		createdAt = r.CreatedAtSnake
	}

	if r.ID == "" {
		return repository.CreatedGroup{}, errors.New("group rpc response id is empty")
	}

	return repository.CreatedGroup{
		ID:          r.ID,
		Name:        r.Name,
		MonthlyGoal: monthlyGoal,
		CreatedAt:   createdAt,
	}, nil
}
