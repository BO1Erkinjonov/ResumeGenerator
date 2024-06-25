CREATE TABLE IF NOT EXISTS resumes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    first_name VARCHAR(200),
    last_name VARCHAR(200),
    category VARCHAR(200) NOT NULL,
    birth_date TIMESTAMP NOT NULL,
    salary VARCHAR(200),
    description TEXT,
    work_type VARCHAR(200) CHECK (work_type IN ('offline', 'online', 'does not matter')) NOT NULL
);
