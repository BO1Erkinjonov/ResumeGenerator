package postgresql

import (
	"context"
	"fmt"
	"resume-generator/internal/entity"
	"resume-generator/internal/pkg/postgres"
)

//TODO::
//type Links interface {
//	GetLinkById(ctx context.Context, req *entity.FieldValueReq) (*entity.Link, error)
//	CheckUniques(ctx context.Context, req *entity.FieldValueReq) (*entity.Link, error)
//	GetAllLinks(ctx context.Context, req *entity.GetAllReq) ([]*entity.Link, error)
//	DeleteLink(ctx context.Context, req *entity.DeleteReq) (*entity.Link, error)
//	UpdateLinkById(ctx context.Context, req *entity.LinksUpdateReq) (*entity.Link, error)
//}

type LinkRepository struct {
	db        *postgres.PostgresDB
	tableName string
}

func NewLinRepository(db *postgres.PostgresDB) *LinkRepository {
	return &LinkRepository{
		db:        db,
		tableName: "links",
	}
}

func (l *LinkRepository) selectLinkQuerySuffix() string {
	return `id,
	resume_id,
	link_name,
	link_url
`
}

func (l *LinkRepository) CreateLink(ctx context.Context, req *entity.Link) (*entity.Link, error) {
	data := map[string]interface{}{
		"id":        req.ID,
		"resume_id": req.ResumeID,
		"link_name": req.LinkName,
		"link_url":  req.LinkURL,
	}
	query, argc, err := l.db.Sq.Builder.Insert(l.tableName).
		SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", l.
		selectLinkQuerySuffix())).ToSql()
	if err != nil {
		return nil, err
	}
	var link entity.Link
	err = l.db.QueryRow(ctx, query, argc...).Scan(
		&link.ID,
		&link.ResumeID,
		&link.LinkName,
		&link.LinkURL)

	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (l *LinkRepository) DeleteLink(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error) {
	user := UserRepo{
		db:        l.db,
		tableName: l.tableName,
	}
	req.IsHardDeleted = true
	return user.DeleteUserById(ctx, req)
}
