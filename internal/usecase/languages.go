package usecase

import (
	"context"
	"resume-generator/internal/entity"
)

type Languages interface {
	CreateLanguage(ctx context.Context, link *entity.Language) (*entity.Language, error)
	GetLanguageById(ctx context.Context, req *entity.FieldValueReq) (*entity.Language, error)
	GetAllLanguage(ctx context.Context, req *entity.GetAllReq) ([]*entity.Language, error)
	DeleteLanguage(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error)
	UpdateLanguageById(ctx context.Context, req *entity.LanguagesUpdateReq) (*entity.Language, error)
}

type languagesUseCase struct {
	language Languages
}

func (l *languagesUseCase) CreateLanguage(ctx context.Context, link *entity.Language) (*entity.Language, error) {
	return l.language.CreateLanguage(ctx, link)
}
func (l *languagesUseCase) GetLanguageById(ctx context.Context, req *entity.FieldValueReq) (*entity.Language, error) {
	return l.language.GetLanguageById(ctx, req)
}
func (l *languagesUseCase) DeleteLanguage(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error) {
	return l.language.DeleteLanguage(ctx, req)
}

func (l *languagesUseCase) UpdateLanguageById(ctx context.Context, req *entity.LanguagesUpdateReq) (*entity.Language, error) {
	return l.language.UpdateLanguageById(ctx, req)
}

func (l *languagesUseCase) GetAllLanguage(ctx context.Context, req *entity.GetAllReq) ([]*entity.Language, error) {
	return l.language.GetAllLanguage(ctx, req)
}
