package entity

import "time"

// CREATE TABLE IF NOT EXISTS resumes (
// id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// user_id UUID REFERENCES users(id),
// first_name VARCHAR(200),
// last_name VARCHAR(200),
// category VARCHAR(200) NOT NULL,
// birth_date TIMESTAMP NOT NULL,
// salary VARCHAR(200),
// description TEXT,
// work_type VARCHAR(200) CHECK (work_type IN ('offline', 'online', 'does not matter')) NOT NULL
// );
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
