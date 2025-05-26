CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    user_like TEXT,
    user_like_embedding vector(768), -- 假设使用pgvector扩展
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_user_id ON users(user_id);