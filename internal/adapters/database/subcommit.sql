CREATE TABLE subcommit (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    title TEXT NOT NULL,
    idea TEXT NOT NULL,
    description TEXT NOT NULL,
    commit_sha TEXT NOT NULL,
    type TEXT NOT NULL,
    epic TEXT NOT NULL,
    files TEXT [] NOT NULL,
    FOREIGN KEY(commit_sha) REFERENCES commit(sha)
)