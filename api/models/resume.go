package models

type Resume struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Category    string `json:"category_id"`
	BirthData   string `json:"birth_date" example:"2020-01-01"`
	Salary      string `json:"salary"`
	Description string `json:"description"`
	WorkType    string `json:"work_type"`
}

type UpdateResume struct {
	ID          string `json:"id"`
	UserID      string `json:"-"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Category    string `json:"category_id"`
	BirthData   string `json:"birth_date" example:"2020-01-01"`
	Salary      string `json:"salary"`
	Description string `json:"description"`
	WorkType    string `json:"-"`
}

type ReqResume struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Category    string `json:"category_id"`
	BirthData   string `json:"birth_date" example:"2020-01-01"`
	Salary      string `json:"salary"`
	Description string `json:"description"`
	WorkType    string `json:"-"`
}

type ListResume struct {
	Resumes []Resume `json:"resumes"`
}
