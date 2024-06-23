CREATE TABLE IF NOT EXISTS soft_skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    skill_name VARCHAR(200),
    resume_id UUID REFERENCES resumes(id)
)