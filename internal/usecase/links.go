package usecase

import (
	"context"
	"resume-generator/internal/entity"
)

type Links interface {
	CreateLink(ctx context.Context, link *entity.Link) (*entity.Link, error)
	GetLinkById(ctx context.Context, req *entity.FieldValueReq) (*entity.Link, error)
	GetAllLinks(ctx context.Context, req *entity.GetAllReq) ([]*entity.Link, error)
	DeleteLink(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error)
	UpdateLinkById(ctx context.Context, req *entity.LinksUpdateReq) (*entity.Link, error)
}

type linkUserCase struct {
	links Links
}

func (l *linkUserCase) CreateLink(ctx context.Context, link *entity.Link) (*entity.Link, error) {
	return l.links.CreateLink(ctx, link)
}
func (l *linkUserCase) GetAllLinks(ctx context.Context, req *entity.GetAllReq) ([]*entity.Link, error) {
	return l.links.GetAllLinks(ctx, req)
}
func (l *linkUserCase) GetLinkById(ctx context.Context, req *entity.FieldValueReq) (*entity.Link, error) {
	return l.links.GetLinkById(ctx, req)
}

func (l *linkUserCase) DeleteLink(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error) {
	return l.links.DeleteLink(ctx, req)
}
func (l *linkUserCase) UpdateLinkById(ctx context.Context, req *entity.LinksUpdateReq) (*entity.Link, error) {
	return l.links.UpdateLinkById(ctx, req)
}

func NewLinkUseCase(links Links) Links {
	return &linkUserCase{links: links}
}
