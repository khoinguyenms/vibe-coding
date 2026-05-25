package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	sqlcdb "github.com/vibe-be/db/sqlc"
	"github.com/vibe-be/internal/model"
	"github.com/vibe-be/internal/repository"
	"github.com/vibe-be/pkg/response"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, in model.CreateUserRequest) (sqlcdb.User, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.FullName = strings.TrimSpace(in.FullName)

	if in.Email == "" || in.FullName == "" || len(in.Password) < 6 {
		return sqlcdb.User{}, response.ErrInvalidInput
	}

	data, err := s.repo.Create(ctx, sqlcdb.CreateUserParams{
		Email:    in.Email,
		FullName: in.FullName,
		Password: in.Password,
	})
	if err != nil {
		return sqlcdb.User{}, err
	}
	return data, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (sqlcdb.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context, limit, offset int32) ([]sqlcdb.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}
