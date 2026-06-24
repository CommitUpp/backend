package postgres

import (
	"context"

	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSQLに対してグループ関連テーブルを操作する実装です。
type GroupRepository struct {
	pool *pgxpool.Pool
}

// PostgreSQL接続プールを使うグループrepositoryを生成します。
func NewGroupRepository(pool *pgxpool.Pool) repository.GroupRepository {
	return &GroupRepository{
		pool: pool,
	}
}

// groupsを作成し、その作成者をgroup_membersへadminとして登録
func (r *GroupRepository) CreateGroupWithOwner(ctx context.Context, input repository.CreateGroupWithOwnerInput) (repository.CreatedGroup, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return repository.CreatedGroup{}, err
	}
	defer tx.Rollback(ctx)

	var createdGroup repository.CreatedGroup
	if err := tx.QueryRow(ctx, `
		INSERT INTO groups (name, monthly_goal)
		VALUES ($1, $2)
		RETURNING id, name, monthly_goal, created_at
	`, input.Name, input.MonthlyGoal).Scan(
		&createdGroup.ID,
		&createdGroup.Name,
		&createdGroup.MonthlyGoal,
		&createdGroup.CreatedAt,
	); err != nil {
		return repository.CreatedGroup{}, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO group_members (group_id, user_id, role)
		VALUES ($1, $2, $3)
	`, createdGroup.ID, input.OwnerUserID, string(input.OwnerRole)); err != nil {
		return repository.CreatedGroup{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.CreatedGroup{}, err
	}

	return createdGroup, nil
}

// 指定されたユーザーが指定されたグループのメンバーであるかを確認
func (r *GroupRepository) IsGroupMember(ctx context.Context, userID string, groupID string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM group_members
			WHERE group_id = $1
			  AND user_id = $2
		)
	`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, groupID, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
