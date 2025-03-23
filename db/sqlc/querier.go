// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateCopy(ctx context.Context, arg CreateCopyParams) (Copywriting, error)
	CreateMultipleCopy(ctx context.Context, arg CreateMultipleCopyParams) error
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteCopy(ctx context.Context, arg DeleteCopyParams) error
	GetCopy(ctx context.Context, id int64) (Copywriting, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id int32) (User, error)
	GetUserByName(ctx context.Context, username string) (User, error)
	ListCopies(ctx context.Context, arg ListCopiesParams) ([]Copywriting, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
