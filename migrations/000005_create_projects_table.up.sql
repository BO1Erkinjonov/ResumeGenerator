CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_name VARCHAR(200) NOT NULL,
    description TEXT,
    link TEXT NOT NULL,
    resume_id UUID REFERENCES resumes(id)
)