CREATE TABLE IF NOT EXISTS educations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    education_type VARCHAR(200) CHECK ( education_type IN ('course', 'university', 'institute')) NOT NULL,
    education_name VARCHAR(200),
    start_year VARCHAR(200) NOT NULL,
    finish_year VARCHAR(200),
    serteficat_link TEXT,
    resume_id UUID REFERENCES resumes(id)
)