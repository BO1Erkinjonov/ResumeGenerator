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

func (l *LinkRepository) updateLinkQuery(req *entity.LinksUpdateReq) map[string]interface{} {
	data := map[string]interface{}{}
	if req.LinkURL != "" {
		data["link_url"] = req.LinkURL
	}
	if req.LinkName != "" {
		data["link_name"] = req.LinkName
	}

	return data
}

func (l *LinkRepository) linkSelectQueryPrefix() string {
	return `
	id,
	resume_id ,
    link_name ,
    link_url 
`
}

func (l *LinkRepository) UpdateLinkById(ctx context.Context, req *entity.LinksUpdateReq) (*entity.Link, error) {
	data := l.updateLinkQuery(req)
	query, argc, err := l.db.Sq.Builder.Update(l.tableName).
		SetMap(data).
		Where(l.db.Sq.Equal("id", req.LinID)).
		Suffix(fmt.Sprintf("RETURNING %s", l.linkSelectQueryPrefix())).
		ToSql()
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, "update"))
	}
	var link entity.Link

	err = l.db.QueryRow(ctx, query, argc...).Scan(
		&link.ID,
		&link.ResumeID,
		&link.LinkName,
		&link.LinkURL)

	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " update"))

	}
	return &link, nil
}

func (l *LinkRepository) GetAllLinks(ctx context.Context, req *entity.GetAllReq) ([]*entity.Link, error) {
	toSql := l.db.Sq.Builder.Select(l.linkSelectQueryPrefix()).From(l.tableName)
	if req.Values != "" && req.Field != "" {
		toSql = toSql.Where(l.db.Sq.Equal(req.Field, req.Values))
	}
	if req.Offset != 0 {
		toSql = toSql.Offset(req.Offset)
	}
	if req.Limit != 0 {
		toSql = toSql.Limit(req.Limit)
	}
	query, argc, err := toSql.ToSql()
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get all"))
	}
	rows, err := l.db.Query(ctx, query, argc...)
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get all"))
	}
	defer rows.Close()
	res := make([]*entity.Link, 0)
	for rows.Next() {
		var resume entity.Link
		err = rows.Scan(
			&resume.ID,
			&resume.ResumeID,
			&resume.LinkName,
			&resume.LinkURL)

		if err != nil {
			return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get all"))
		}
		res = append(res, &resume)
	}
	if err := rows.Err(); err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get all"))
	}
	return res, nil
}

func (l *LinkRepository) GetLinkById(ctx context.Context, req *entity.FieldValueReq) (*entity.Link, error) {
	query, argc, err := l.db.Sq.Builder.Select(l.linkSelectQueryPrefix()).
		From(l.tableName).Where(l.db.Sq.Equal(req.Field, req.Value)).ToSql()
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get "))
	}
	var user entity.Link
	err = l.db.QueryRow(ctx, query, argc...).Scan(
		&user.ID,
		&user.ResumeID,
		&user.LinkName,
		&user.LinkURL)
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get "))
	}
	return &user, nil
}
