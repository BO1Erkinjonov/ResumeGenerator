package entity

import "time"

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"user_name"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type FieldValueReq struct {
	Field string `json:"field"`
	Value string `json:"value"`
}
type GetAllReq struct {
	Field  string `schema:"field"`
	Values string `schema:"values"`
	Limit  uint64 `schema:"limit"`
	Offset uint64 `schema:"offset"`
}

type Result struct {
	IsExists bool `json:"resp"`
}

type DeleteReq struct {
	ID            string `json:"id"`
	IsHardDeleted bool   `json:"is_hard_deleted"`
}

type UpdateUserReq struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}
