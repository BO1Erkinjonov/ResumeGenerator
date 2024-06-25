package models

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"user_name"`
	ImageUrl  string `json:"image_url"`
}

type UserBody struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"user_name"`
	ImageUrl  string `json:"image_url"`
}

type AccessToken struct {
	Access string `json:"access"`
}

type ListUsers struct {
	Users []UserBody
	Count int
}

type UpdateUserReq struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type DeleteReq struct {
	ID string `json:"id"`
}
