package usecase

import (
	"context"
	"resume-generator/internal/entity"
)

type User interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserById(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error)
	CheckUniques(ctx context.Context, user *entity.FieldValueReq) (*entity.Result, error)
	GetAllUsers(ctx context.Context, req *entity.GetAllReq) ([]*entity.User, error)
	DeleteUserById(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error)
	UpdateUserById(ctx context.Context, req *entity.UpdateUserReq) (*entity.Result, error)
}

type userUseCase struct {
	repo User
}

func (u userUseCase) GetAllUsers(ctx context.Context, req *entity.GetAllReq) ([]*entity.User, error) {
	return u.repo.GetAllUsers(ctx, req)
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

func (u userUseCase) DeleteUserById(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error) {
	return u.repo.DeleteUserById(ctx, req)
}

func (u userUseCase) UpdateUserById(ctx context.Context, req *entity.UpdateUserReq) (*entity.Result, error) {
	return u.repo.UpdateUserById(ctx, req)
}
