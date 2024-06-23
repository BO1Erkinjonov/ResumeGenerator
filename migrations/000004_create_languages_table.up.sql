CREATE TABLE IF NOT EXISTS languages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    language_name VARCHAR(200) NOT NULL,
    language_level VARCHAR(200) CHECK (language_level IN ('A1', 'A2', 'B1', 'B2', 'C1', 'C2')) NOT NULL,
    resume_id UUID REFERENCES resumes(id)
)