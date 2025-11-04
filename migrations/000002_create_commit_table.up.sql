CREATE TABLE IF NOT EXISTS commit (
    sha VARCHAR(40) PRIMARY KEY, 
    created_at TIMESTAMP WITH TIME ZONE,
    author VARCHAR(100),
    date VARCHAR(50),
    message TEXT,
    url TEXT,
    author_email VARCHAR(100),
    description TEXT, 
    author_url TEXT,
    files TEXT[], 
    repo_id BIGINT NOT NULL,
    
    CONSTRAINT fk_repo
        FOREIGN KEY(repo_id) 
        REFERENCES repository(id)
        ON DELETE CASCADE
);