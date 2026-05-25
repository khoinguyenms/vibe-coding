package service

import (
	"context"

	"github.com/google/uuid"
	sqlc "github.com/vibe-be/db/sqlc"
	"github.com/vibe-be/internal/model"
)

type UserService interface {
	Create(ctx context.Context, in model.CreateUserRequest) (sqlc.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	List(ctx context.Context, limit, offset int32) ([]sqlc.User, error)
}
