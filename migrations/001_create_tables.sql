CREATE TABLE IF NOT EXISTS repository (
    id                       BIGINT PRIMARY KEY,
    name                     TEXT NOT NULL,
    url                      TEXT NOT NULL UNIQUE,
    last_analyzed_commit_sha TEXT NOT NULL DEFAULT '',
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subcommit (
    id                BIGSERIAL PRIMARY KEY,
    title             TEXT NOT NULL,
    idea              TEXT NOT NULL DEFAULT '',
    description       TEXT NOT NULL,
    epic              TEXT NOT NULL DEFAULT '',
    modification_type TEXT NOT NULL,
    commit_sha        TEXT NOT NULL,
    files             TEXT[] NOT NULL DEFAULT '{}',
    repo_id           BIGINT NOT NULL REFERENCES repository(id),
    committed_at      TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_subcommit_repo_id ON subcommit(repo_id);
CREATE INDEX IF NOT EXISTS idx_subcommit_repo_commit ON subcommit(repo_id, commit_sha);
