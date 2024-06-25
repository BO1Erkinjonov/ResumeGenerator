package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"resume-generator/internal/entity"
	"resume-generator/internal/pkg/postgres"
	"time"
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

func (u *UserRepo) userUpdateQuery(req *entity.UpdateUserReq) map[string]interface{} {
	date := make(map[string]interface{})
	if req.Password != "" {
		date["password"] = req.Password
	}
	if req.FirstName != "" {
		date["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		date["last_name"] = req.LastName
	}
	if req.UserName != "" {
		date["user_name"] = req.UserName
	}
	date["updated_at"] = time.Now()

	return date
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
	if user.DeletedAt.IsZero() {
		if err != nil {
			return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " get"))
		}
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

func (u *UserRepo) GetAllUsers(ctx context.Context, req *entity.GetAllReq) ([]*entity.User, error) {

	tosql := u.db.Sq.Builder.Select(u.userSelectQueryPrefix()).From(u.tableName)
	if req.Field != "" && req.Values != "" {
		tosql = tosql.Where(u.db.Sq.Equal(req.Field, req.Values))
	}
	if req.Offset != 0 {
		tosql = tosql.Offset(req.Offset)
	}
	if req.Limit != 0 {
		tosql = tosql.Limit(req.Limit)
	}

	query, argc, err := tosql.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
	}
	rows, err := u.db.Query(ctx, query, argc...)
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
	}
	defer rows.Close()

	users := []*entity.User{}

	for rows.Next() {
		var user entity.User
		var updatedAt, deletedAt sql.NullTime
		err = rows.Scan(
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
		if err != nil {
			return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			user.DeletedAt = deletedAt.Time
		}

		if user.DeletedAt.IsZero() {
			users = append(users, &user)
		}

	}

	if err := rows.Err(); err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " all"))
	}

	return users, nil
}

func (u *UserRepo) DeleteUserById(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error) {
	now := time.Now()
	query, argc, err := u.db.Sq.Builder.Update(u.tableName).Set("deleted_at", now).
		Where(u.db.Sq.Equal("id", req.ID)).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " update"))
	}
	_, err = u.db.Exec(ctx, query, argc...)
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " update"))
	}
	return &entity.Result{IsExists: true}, nil
}

func (u *UserRepo) UpdateUserById(ctx context.Context, req *entity.UpdateUserReq) (*entity.User, error) {
	data := u.userUpdateQuery(req)
	query, argc, err := u.db.Sq.Builder.Update(u.tableName).SetMap(data).
		Where(u.db.Sq.Equal("id", req.UserId)).Suffix(fmt.Sprintf("RETURNING %s", u.userSelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " update"))
	}
	fmt.Println(query, argc)
	var user entity.User
	var updatedAt, deletedAt sql.NullTime
	err = u.db.QueryRow(ctx, query, argc...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.ImageUrl,
		&user.CreatedAt,
		&deletedAt,
		&updatedAt)
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, " update"))
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		user.DeletedAt = deletedAt.Time
	}

	if !user.DeletedAt.IsZero() {
		return &user, nil
	}

	return nil, errors.New("user delete failed")
}
