package repository

import (
	"context"

	"github.com/google/uuid"
	sqlc "github.com/vibe-be/db/sqlc"
)

type UserRepository interface {
	Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	GetByEmail(ctx context.Context, email string) (sqlc.User, error)
	List(ctx context.Context, limit, offset int32) ([]sqlc.User, error)
	Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
