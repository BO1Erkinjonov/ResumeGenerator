package usecase

import (
	"context"
	"resume-generator/internal/entity"
)

type Resume interface {
	CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
	GetResumeById(ctx context.Context, req *entity.FieldValueReq) (*entity.Resume, error)
	CheckUniques(ctx context.Context, req *entity.FieldValueReq) (*entity.Result, error)
	GetAllResumes(ctx context.Context, req *entity.GetAllReq) ([]*entity.Resume, error)
	DeleteResume(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error)
	UpdateResumeById(ctx context.Context, req *entity.UpdateResumeReq) (*entity.Resume, error)
}

type resumeUserCase struct {
	repo Resume
}

func (uc *resumeUserCase) CreateResume(ctx context.Context,
	resume *entity.Resume) (*entity.Resume, error) {
	return uc.repo.CreateResume(ctx, resume)
}
func (uc *resumeUserCase) GetResumeById(ctx context.Context,
	req *entity.FieldValueReq) (*entity.Resume, error) {
	return uc.repo.GetResumeById(ctx, req)
}
func (uc *resumeUserCase) CheckUniques(ctx context.Context,
	req *entity.FieldValueReq) (*entity.Result, error) {
	return uc.repo.CheckUniques(ctx, req)
}
func (uc *resumeUserCase) GetAllResumes(ctx context.Context,
	req *entity.GetAllReq) ([]*entity.Resume, error) {
	return uc.repo.GetAllResumes(ctx, req)
}
func (uc *resumeUserCase) DeleteResume(ctx context.Context,
	req *entity.DeleteReq) (*entity.Result, error) {
	return uc.repo.DeleteResume(ctx, req)
}
func (uc *resumeUserCase) UpdateResumeById(ctx context.Context,
	req *entity.UpdateResumeReq) (*entity.Resume, error) {
	return uc.repo.UpdateResumeById(ctx, req)
}
