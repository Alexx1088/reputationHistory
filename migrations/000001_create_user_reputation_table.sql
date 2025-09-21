CREATE TABLE IF NOT EXISTS user_reputation (
                                               user_id UUID PRIMARY KEY,
                                               score   INTEGER NOT NULL DEFAULT 0
);