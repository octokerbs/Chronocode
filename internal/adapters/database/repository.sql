CREATE TABLE repository (
    id PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    last_analyzed_commit VARCHAR(255)
)