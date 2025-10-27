CREATE TABLE IF NOT EXISTS subcommits (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    title VARCHAR(255) NOT NULL,
    idea TEXT,
    description TEXT,
    commit_sha VARCHAR(40) NOT NULL, 
    type VARCHAR(50),
    epic VARCHAR(100),
    files TEXT[], 

    CONSTRAINT fk_commit
        FOREIGN KEY(commit_sha)
        REFERENCES commit(sha)
        ON DELETE CASCADE
);