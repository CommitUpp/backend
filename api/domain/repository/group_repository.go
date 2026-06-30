package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrGroupNotFound            = errors.New("group not found")
	ErrGroupMemberAlreadyExists = errors.New("group member already exists")
)

// role に保存する権限種別
type GroupMemberRole string

const (
	RoleAdmin  GroupMemberRole = "admin"
	RoleEditor GroupMemberRole = "editor"
	RoleUser   GroupMemberRole = "user"
)

// groupsとgroup_members を同時に作成するための入力値
type CreateGroupWithOwnerInput struct {
	Name        string
	MonthlyGoal int
	AccessToken string
}

type JoinGroupInput struct {
	GroupID     string
	AccessToken string
}

// groups作成後にDBから返してほしい値です。
type CreatedGroup struct {
	ID          string
	Name        string
	MonthlyGoal int
	CreatedAt   time.Time
}

type GroupRepository interface {
	CreateGroupWithOwner(ctx context.Context, input CreateGroupWithOwnerInput) (CreatedGroup, error)
	JoinGroup(ctx context.Context, input JoinGroupInput) (CreatedGroup, error)
}
