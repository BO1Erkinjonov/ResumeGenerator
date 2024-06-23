CREATE TABLE IF NOT EXISTS work_xp (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_name VARCHAR(200) NOT NULL,
    description TEXT,
    work_position VARCHAR(200),
    company_link TEXT,
    start_year VARCHAR(200) NOT NULL,
    finish_year VARCHAR(200),
    resume_id UUID REFERENCES resumes(id)
)