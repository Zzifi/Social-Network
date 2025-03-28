CREATE TABLE IF NOT EXISTS USERS (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS SESSIONS (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES USERS(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS USER_PROFILE (
    user_id INT PRIMARY KEY REFERENCES USERS(id) ON DELETE CASCADE,
    phone_number VARCHAR(20),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    birthday VARCHAR(10),
    bio TEXT,
    avatar_url TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_profile_update
BEFORE UPDATE ON USER_PROFILE
FOR EACH ROW
EXECUTE FUNCTION set_timestamp();
