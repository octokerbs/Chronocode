CREATE TABLE repository_records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    last_analyzed_commit VARCHAR(40)
);

-- TODO: Escribir lo mismo para create commit y para subcommit