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
