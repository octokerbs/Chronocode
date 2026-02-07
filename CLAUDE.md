# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Chronocode is a Go backend that analyzes GitHub repositories using AI (Google Gemini / Gemma LLMs) to break commits into logical "subcommits". It uses concurrent worker pools for parallel analysis and persistence.

## Build & Run Commands

```bash
# Full rebuild (stops containers, removes DB volume, rebuilds, starts fresh)
./scripts/rebuild.sh

# Standard Docker Compose operations
docker compose up --build -d
docker compose down

# Build Go binary locally
go build -o main ./cmd/http-server/main.go

# Run locally (requires .env with DATABASE_URL, GEMINI_API_KEY, GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET)
go run ./cmd/http-server/main.go
```

There are no tests in this codebase currently. `testify` is imported but unused.

## Architecture

The project follows **hexagonal architecture** (ports & adapters) with clean separation:

```
cmd/http-server/main.go          Entry point, dependency wiring (buildDependencies)
internal/
  api/http/                      Gin HTTP server, handlers, error mapping
  application/                   Business logic services (Analyzer, Auth, Querier, etc.)
  domain/                        Entities & port interfaces (no external dependencies)
    analysis/                    Commit, Subcommit, Repository entities + AgentPort
    auth/                        AuthPort interface
    codehost/                    CodeHostPort, CodeHostFactoryPort interfaces
    database/                    DatabasePort interface
  infrastructure/                Adapter implementations
    agent/gemini/                Google Gemini adapter (primary AI)
    agent/gemma/                 Local Gemma LLM adapter (alternative)
    auth/githubauth/             GitHub OAuth2 adapter
    codehost/githubapi/          GitHub API adapter (repos, commits, diffs)
    database/postgres/           PostgreSQL adapter
  errors/                       Custom error types with category (BadRequest, InternalFailure, etc.)
config/                          Environment-based configuration
migrations/                      PostgreSQL migrations (3 tables: repository, commit, subcommit)
```

## Key Patterns

- **Dependency injection**: All wired in `main.go:buildDependencies()`, interfaces defined as domain ports
- **Worker pools**: Analyzer uses 5 goroutines for AI analysis, 100 goroutines for DB persistence, connected via buffered channels (capacity 100)
- **Port interfaces**: `AgentPort`, `DatabasePort`, `CodeHostPort`, `CodeHostFactoryPort`, `AuthPort` — add new adapters by implementing these
- **Error flow**: Domain errors (`errors.Error` with category) are mapped to HTTP status codes in `handler/error.go`

## API Endpoints

- `GET /auth/github/login` — OAuth redirect
- `GET /auth/github/callback` — OAuth callback, sets `access_token` cookie
- `POST /analyze?repo_url=...&github_token=...` — Start async repo analysis (returns 202 with repo ID)
- `GET /subcommits-timeline?repo_id=...` — Fetch analyzed subcommits

## Go Module

Module path: `github.com/octokerbs/chronocode-backend`
Go version: 1.23.7
HTTP framework: Gin
Logging: Uber Zap
Database driver: lib/pq (PostgreSQL)

## Docker Setup

Three services in `docker-compose.yml`:
1. **postgres** (port 5468 external → 5432 internal)
2. **migrator** — runs `migrate` against postgres, then exits
3. **app** — Go binary (port 8080), depends on migrator completing
