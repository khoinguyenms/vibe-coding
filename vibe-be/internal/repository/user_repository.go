package repository

import (
	"context"

	"github.com/google/uuid"
	sqlc "github.com/vibe-be/db/sqlc"
)

type userRepository struct {
	db sqlc.Querier
}

func NewUserRepository(db sqlc.Querier) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.db.CreateUser(ctx, arg)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
	return r.db.GetUserByID(ctx, id)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return r.db.GetUserByEmail(ctx, email)
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]sqlc.User, error) {
	return r.db.ListUsers(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
}

func (r *userRepository) Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	return r.db.UpdateUser(ctx, arg)
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteUser(ctx, id)
}
