package entity

import "time"

type Resume struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Category    string    `json:"category_id"`
	BirthData   time.Time `json:"birth_date"`
	Salary      string    `json:"salary"`
	Description string    `json:"description"`
	WorkType    string    `json:"work_type"`
}

type UpdateResumeReq struct {
	ResumeID    string    `json:"resume_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Category    string    `json:"category"`
	BirthDate   time.Time `json:"birth_date"`
	Salary      string    `json:"salary"`
	Description string    `json:"description"`
	WorkType    string    `json:"work_type"`
}
