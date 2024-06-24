package usecase

import (
	"context"
	"resume-generator/internal/entity"
)

type User interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserById(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error)
	CheckUniques(ctx context.Context, user *entity.FieldValueReq) (*entity.Result, error)
}

type userUseCase struct {
	repo User
}

func (u userUseCase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u userUseCase) GetUserById(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error) {
	return u.repo.GetUserById(ctx, req)
}

func (u userUseCase) CheckUniques(ctx context.Context, user *entity.FieldValueReq) (*entity.Result, error) {
	return u.repo.CheckUniques(ctx, user)
}

func NewUserUseCase(u User) *userUseCase {
	return &userUseCase{u}
}
