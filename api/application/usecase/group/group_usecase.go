package group

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/CommitUpp/backend/api/domain/repository"
)

// グループ目標月間視聴数（作成時は一人なので3）
const moviesPerMember = 3

var (
	// 認証済みユーザーIDがusecaseへ渡されなかった場合に返すエラー
	ErrInvalidUserID = errors.New("invalid user id")
	// グループ名が空、または空白だけだった場合に返すエラー
	ErrInvalidName = errors.New("invalid group name")
	// グループIDが空、または空白だけだった場合に返すエラー
	ErrInvalidGroupID     = errors.New("invalid group id")
	ErrInvalidAccessToken = errors.New("invalid access token")
	// グループ保存用repositoryが未設定の場合に返すエラー
	ErrGroupRepositoryNotConfigured = errors.New("group repository is not configured")
	ErrGroupNotFound                = repository.ErrGroupNotFound
	ErrGroupMemberAlreadyExists     = repository.ErrGroupMemberAlreadyExists
)

// group_membersに登録するためのカラム
type CreateGroupInput struct {
	UserID      string
	Name        string
	AccessToken string
}

// グループ作成の結果をhandlerへ返すための構造体
type CreateGroupOutput struct {
	ID          string
	Name        string
	MonthlyGoal int
	CreatedAt   time.Time
}

type JoinGroupInput struct {
	GroupID     string
	UserID      string
	AccessToken string
}

type JoinGroupOutput struct {
	ID          string
	Name        string
	MonthlyGoal int
	CreatedAt   time.Time
}

// グループ関連のアプリケーションロジックを定義
type GroupUsecase interface {
	CreateGroup(ctx context.Context, input CreateGroupInput) (CreateGroupOutput, error)
	JoinGroup(ctx context.Context, input JoinGroupInput) (JoinGroupOutput, error)
}

type groupUsecaseImpl struct {
	groupRepository repository.GroupRepository
}

// グループ関連usecaseの実装を生成
func NewGroupUsecase(groupRepository repository.GroupRepository) GroupUsecase {
	return &groupUsecaseImpl{
		groupRepository: groupRepository,
	}
}

// handlerから渡された値をグループ作成用に検証・整形
func (u *groupUsecaseImpl) CreateGroup(ctx context.Context, input CreateGroupInput) (CreateGroupOutput, error) {
	if input.UserID == "" {
		return CreateGroupOutput{}, ErrInvalidUserID
	}

	if input.AccessToken == "" {
		return CreateGroupOutput{}, ErrInvalidAccessToken
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return CreateGroupOutput{}, ErrInvalidName
	}

	if u.groupRepository == nil {
		return CreateGroupOutput{}, ErrGroupRepositoryNotConfigured
	}

	createdGroup, err := u.groupRepository.CreateGroupWithOwner(ctx, repository.CreateGroupWithOwnerInput{
		Name:        name,
		MonthlyGoal: moviesPerMember,
		AccessToken: input.AccessToken,
	})
	if err != nil {
		return CreateGroupOutput{}, err
	}

	return CreateGroupOutput{
		ID:          createdGroup.ID,
		Name:        createdGroup.Name,
		MonthlyGoal: createdGroup.MonthlyGoal,
		CreatedAt:   createdGroup.CreatedAt,
	}, nil
}

func (u *groupUsecaseImpl) JoinGroup(ctx context.Context, input JoinGroupInput) (JoinGroupOutput, error) {
	if input.UserID == "" {
		return JoinGroupOutput{}, ErrInvalidUserID
	}

	if input.AccessToken == "" {
		return JoinGroupOutput{}, ErrInvalidAccessToken
	}

	groupID := strings.TrimSpace(input.GroupID)
	if groupID == "" {
		return JoinGroupOutput{}, ErrInvalidGroupID
	}

	if u.groupRepository == nil {
		return JoinGroupOutput{}, ErrGroupRepositoryNotConfigured
	}

	joinedGroup, err := u.groupRepository.JoinGroup(ctx, repository.JoinGroupInput{
		GroupID:     groupID,
		AccessToken: input.AccessToken,
	})
	if err != nil {
		return JoinGroupOutput{}, err
	}

	return JoinGroupOutput{
		ID:          joinedGroup.ID,
		Name:        joinedGroup.Name,
		MonthlyGoal: joinedGroup.MonthlyGoal,
		CreatedAt:   joinedGroup.CreatedAt,
	}, nil
}
