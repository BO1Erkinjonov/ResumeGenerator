package postgresql

import (
	"context"
	"fmt"
	"resume-generator/internal/entity"
	"resume-generator/internal/pkg/postgres"
)

//type Languages interface {
//	GetLanguageById(ctx context.Context, req *entity.FieldValueReq) (*entity.Language, error)

//}

type LanguageRepository struct {
	db        *postgres.PostgresDB
	tableName string
}

func NewLanguageRepository(db *postgres.PostgresDB) *LanguageRepository {
	return &LanguageRepository{
		db:        db,
		tableName: "languages",
	}
}

func (l *LanguageRepository) selectLanguageQuerySuffix() string {
	return `id,
	language_name,
	language_level,
	resume_id
`
}

func (l *LanguageRepository) CreateLanguage(ctx context.Context, link *entity.Language) (*entity.Language, error) {
	data := map[string]interface{}{
		"id":             link.ID,
		"language_name":  link.Name,
		"language_level": link.Level,
		"resume_id":      link.ResumeID,
	}
	query, argc, err := l.db.Sq.Builder.Insert(l.tableName).
		SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", l.
		selectLanguageQuerySuffix())).ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(query, argc)
	var links entity.Language
	err = l.db.QueryRow(ctx, query, argc...).Scan(
		&links.ID,
		&links.Name,
		&links.Level,
		&links.ResumeID)

	if err != nil {
		return nil, err
	}
	return &links, nil
}

func (l *LanguageRepository) DeleteLanguage(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error) {
	user := UserRepo{
		db:        l.db,
		tableName: l.tableName,
	}
	req.IsHardDeleted = true
	return user.DeleteUserById(ctx, req)
}

func (l *LanguageRepository) updateLinkQuery(req *entity.LanguagesUpdateReq) map[string]interface{} {
	data := map[string]interface{}{}
	if req.Name != "" {
		data["language_name"] = req.Name
	}
	if req.Level != "" {
		data["language_level"] = req.Level
	}

	return data
}

func (l *LanguageRepository) UpdateLanguageById(ctx context.Context, req *entity.LanguagesUpdateReq) (*entity.Language, error) {
	data := l.updateLinkQuery(req)
	query, argc, err := l.db.Sq.Builder.Update(l.tableName).
		SetMap(data).
		Where(l.db.Sq.Equal("id", req.ID)).
		Suffix(fmt.Sprintf("RETURNING %s", l.selectLanguageQuerySuffix())).
		ToSql()
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, "update"))
	}
	var link entity.Language

	err = l.db.QueryRow(ctx, query, argc...).Scan(
		&link.ID,
		&link.Name,
		&link.Level,
		&link.ResumeID)

	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " update"))

	}
	return &link, nil
}

func (l *LanguageRepository) GetAllLanguage(ctx context.Context, req *entity.GetAllReq) ([]*entity.Language, error) {
	toSql := l.db.Sq.Builder.Select(l.selectLanguageQuerySuffix()).From(l.tableName)
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
	res := make([]*entity.Language, 0)
	for rows.Next() {
		var resume entity.Language
		err = rows.Scan(
			&resume.ID,
			&resume.Name,
			&resume.Level,
			&resume.ResumeID)

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

func (l *LanguageRepository) GetLanguageById(ctx context.Context, req *entity.FieldValueReq) (*entity.Language, error) {
	query, argc, err := l.db.Sq.Builder.Select(l.selectLanguageQuerySuffix()).
		From(l.tableName).Where(l.db.Sq.Equal(req.Field, req.Value)).ToSql()
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get "))
	}
	var user entity.Language
	err = l.db.QueryRow(ctx, query, argc...).Scan(
		&user.ID,
		&user.Name,
		&user.Level,
		&user.ResumeID)
	if err != nil {
		return nil, l.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", l.tableName, " get "))
	}
	return &user, nil
}
