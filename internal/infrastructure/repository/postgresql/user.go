package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"resume-generator/internal/entity"
	"resume-generator/internal/pkg/postgres"
)

type UserRepo struct {
	db        *postgres.PostgresDB
	tableName string
}

func NewUserRepo(pg *postgres.PostgresDB) *UserRepo {
	return &UserRepo{
		db:        pg,
		tableName: "users",
	}
}

func (u *UserRepo) userSelectQueryPrefix() string {
	return `id,
            first_name,
            last_name,
			email,
			password,
			user_name,
			image_url,
			created_at,
			updated_at,
			deleted_at`
}

func (u *UserRepo) CreateUser(ctx context.Context, req *entity.User) (*entity.User, error) {
	data := map[string]interface{}{
		"id":         req.ID,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"email":      req.Email,
		"password":   req.Password,
		"user_name":  req.Username,
		"image_url":  req.ImageUrl,
	}
	query, args, err := u.db.Sq.Builder.Insert(u.tableName).
		SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", u.userSelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " create"))
	}
	var user entity.User
	var trash sql.NullTime
	err = u.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.ImageUrl,
		&user.CreatedAt,
		&trash,
		&trash,
	)
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " create"))
	}
	return &user, nil
}

func (u *UserRepo) GetUserById(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error) {
	query, args, err := u.db.Sq.Builder.Select(u.userSelectQueryPrefix()).
		From(u.tableName).Where(u.db.Sq.Equal(req.Field, req.Value)).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " get"))
	}
	var user entity.User
	var updatedAt, deletedAt sql.NullTime
	err = u.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.ImageUrl,
		&user.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		user.DeletedAt = deletedAt.Time
	}
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " get"))
	}
	return &user, nil
}

func (u *UserRepo) CheckUniques(ctx context.Context, req *entity.FieldValueReq) (*entity.Result, error) {
	query, args, err := u.db.Sq.Builder.Select("count(1)").
		From(u.tableName).Where(u.db.Sq.Equal(req.Field, req.Value)).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " check"))
	}
	var count int
	err = u.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " check"))
	}
	if count != 0 {
		return &entity.Result{IsExists: true}, err
	}
	return &entity.Result{IsExists: false}, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, req *entity.GetAllUserReq) ([]*entity.User, error) {
	query, args, err := u.db.Sq.Builder.Select(u.userSelectQueryPrefix()).From(u.tableName).
		Where(u.db.Sq.Equal(req.Field, req.Values)).Limit(req.Limit).Offset(req.Offset).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
	}
	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
	}
	defer rows.Close()
	users := make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.Username,
			&user.ImageUrl,
			&user.CreatedAt,
			&user.UpdatedAt)
		if err != nil {
			return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
		}
		users = append(users, &user)

	}
	if err := rows.Err(); err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
	}

	return users, nil
}
