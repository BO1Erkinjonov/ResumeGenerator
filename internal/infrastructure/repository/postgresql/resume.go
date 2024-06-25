package postgresql

import (
	"context"
	"fmt"
	"resume-generator/internal/entity"
	"resume-generator/internal/pkg/postgres"
)

type ResumeRepository struct {
	db        *postgres.PostgresDB
	tableName string
}

func NewResumeRepo(pg *postgres.PostgresDB) *ResumeRepository {
	return &ResumeRepository{
		db:        pg,
		tableName: "resumes",
	}
}

func (u *ResumeRepository) resumeSelectQueryPrefix() string {
	return ` id,
            user_id,
            first_name,
			last_name,
			category,
			birth_date,
			salary,
			description,
			work_type `
}

func (r *ResumeRepository) CreateResume(ctx context.Context, req *entity.Resume) (*entity.Resume, error) {
	data := map[string]interface{}{
		"id":          req.ID,
		"user_id":     req.UserID,
		"first_name":  req.FirstName,
		"last_name":   req.LastName,
		"category":    req.Category,
		"birth_date":  req.BirthData,
		"salary":      req.Salary,
		"description": req.Description,
		"work_type":   req.WorkType,
	}

	query, args, err := r.db.Sq.Builder.Insert(r.tableName).
		SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", r.resumeSelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " create"))
	}
	var user entity.Resume
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Category,
		&req.BirthData,
		&user.Salary,
		&user.Description,
		&user.WorkType)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " create "))
	}
	return &user, nil
}

func (r *ResumeRepository) GetResumeById(ctx context.Context, req *entity.FieldValueReq) (*entity.Resume, error) {
	query, argc, err := r.db.Sq.Builder.Select(r.resumeSelectQueryPrefix()).
		From(r.tableName).Where(r.db.Sq.Equal(req.Field, req.Value)).ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " get "))
	}
	var user entity.Resume
	err = r.db.QueryRow(ctx, query, argc...).Scan(
		&user.ID,
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Category,
		&user.BirthData,
		&user.Salary,
		&user.Description,
		&user.WorkType)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " get "))
	}
	return &user, nil
}

func (r *ResumeRepository) CheckUniques(ctx context.Context, req *entity.FieldValueReq) (*entity.Result, error) {
	u := UserRepo{
		db:        r.db,
		tableName: r.tableName,
	}

	return u.CheckUniques(ctx, req)
}

func (r *ResumeRepository) GetAllResumes(ctx context.Context, req *entity.GetAllReq) ([]*entity.Resume, error) {
	toSql := r.db.Sq.Builder.Select(r.resumeSelectQueryPrefix()).From(r.tableName)
	if req.Values != "" && req.Field != "" {
		toSql = toSql.Where(r.db.Sq.Equal(req.Field, req.Values))
	}
	if req.Offset != 0 {
		toSql = toSql.Offset(req.Offset)
	}
	if req.Limit != 0 {
		toSql = toSql.Limit(req.Limit)
	}
	query, argc, err := toSql.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " get all"))
	}
	rows, err := r.db.Query(ctx, query, argc...)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " get all"))
	}
	defer rows.Close()
	res := make([]*entity.Resume, 0)
	for rows.Next() {
		var resume entity.Resume
		err = rows.Scan(
			&resume.ID,
			&resume.UserID,
			&resume.FirstName,
			&resume.LastName,
			&resume.Category,
			&resume.BirthData,
			&resume.Salary,
			&resume.Description,
			&resume.WorkType)

		if err != nil {
			return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " get all"))
		}
		res = append(res, &resume)
	}
	if err := rows.Err(); err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " get all"))
	}
	return res, nil
}

func (r *ResumeRepository) DeleteResume(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error) {
	u := UserRepo{
		db:        r.db,
		tableName: r.tableName,
	}
	req.IsHardDeleted = true
	return u.DeleteUserById(ctx, req)
}

func (r *ResumeRepository) queryUserUpdate(req *entity.UpdateResumeReq) map[string]interface{} {
	data := map[string]interface{}{}
	if req.FirstName != "" {
		data["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		data["last_name"] = req.LastName
	}
	if req.Description != "" {
		data["description"] = req.Description

	}
	if req.WorkType != "" {
		data["work_type"] = req.WorkType
	}
	if req.Salary != "" {
		data["salary"] = req.Salary
	}
	if req.Category != "" {
		data["category"] = req.Category
	}
	if !req.BirthDate.IsZero() {
		data["birth_date"] = req.BirthDate
	}
	return data
}

func (r *ResumeRepository) UpdateResumeById(ctx context.Context, req *entity.UpdateResumeReq) (*entity.Resume, error) {

	data := r.queryUserUpdate(req)
	query, argc, err := r.db.Sq.Builder.Update(r.tableName).
		SetMap(data).
		Where(r.db.Sq.Equal("id", req.ResumeID)).
		Suffix(fmt.Sprintf("RETURNING %s", r.resumeSelectQueryPrefix())).
		ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "update"))
	}
	var user entity.Resume

	err = r.db.QueryRow(ctx, query, argc...).Scan(
		&user.ID,
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Category,
		&user.BirthData,
		&user.Salary,
		&user.Description,
		&user.WorkType)

	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, " update"))

	}
	return &user, nil
}
