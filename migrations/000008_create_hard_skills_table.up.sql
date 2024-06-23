CREATE TABLE IF NOT EXISTS hard_skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    skill_name VARCHAR(200),
    skill_level VARCHAR(200) CHECK (skill_level IN ('beginner', 'medium', 'expert')) NOT NULL,
    resume_id UUID REFERENCES resumes(id)
    )