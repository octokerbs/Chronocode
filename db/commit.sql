CREATE TABLE commit (
    sha VARCHAR(40) NOT NULL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    author VARCHAR(100) NOT NULL,
    date VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    url TEXT NOT NULL,
    author_email VARCHAR(100) NOT NULL,
    description TEXT,
    author_url TEXT,
    files TEXT [],
    repo_id INT NOT NULL,
    
    FOREIGN KEY (repo_id) REFERENCES repository(id)
)