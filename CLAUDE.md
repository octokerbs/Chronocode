# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Chronocode is a GitHub repository analyzer that uses Gemini AI to analyze commits, break them into logical "subcommits" (units of work), and store results in PostgreSQL. It authenticates via GitHub OAuth and exposes an HTTP API.

## Build & Run Commands

```bash
# Full rebuild (stops containers, removes DB volume, rebuilds and starts)
./scripts/rebuild.sh

# Or manually:
docker compose down
docker compose build
docker compose up -d

# Run locally (requires .env file)
go run ./cmd/http-server/main.go

# Run tests
go test ./...

# Run single test
go test -run TestName ./path/to/package
```

## Required Environment Variables

Create a `.env` file with:
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` - Database credentials
- `GEMINI_API_KEY` - Google Gemini API key
- `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET` - GitHub OAuth app credentials
- `REDIRECT_URL` - OAuth callback URL (e.g., `http://localhost:8080/auth/github/callback`)
- `DATABASE_URL` - PostgreSQL connection string (auto-set in Docker)

## Architecture

### Hexagonal/Ports-and-Adapters Structure

```
cmd/http-server/main.go     # Entry point, dependency wiring
config/                     # Environment configuration
internal/
  api/http/                 # HTTP handlers and routing (Gin framework)
  application/              # Use cases/services (Analyzer, Querier, Auth)
  domain/                   # Core domain models and port interfaces
    analysis/               # Commit, Subcommit, Repository models; Agent port
    codehost/               # CodeHost port for fetching commits
    database/               # Database port
  infrastructure/           # Implementations of ports
    agent/gemini/           # Gemini AI implementation
    auth/githubauth/        # GitHub OAuth implementation
    codehost/githubapi/     # GitHub API implementation
    database/postgres/      # PostgreSQL implementation
  errors/                   # Domain error types
migrations/                 # SQL migration files
```

### Key Domain Concepts

- **Commit**: A Git commit with metadata and AI-generated analysis
- **Subcommit**: A logical unit of work within a commit (identified by AI), with title, description, type (FEATURE/BUG/REFACTOR/DOCS/CHORE/MILESTONE/WARNING), and epic
- **Repository**: Tracks analyzed repos and their last analyzed commit SHA

### Core Flow

1. User authenticates via GitHub OAuth
2. User requests analysis of a repository via `POST /analyze`
3. `Analyzer` spawns worker pool (5 workers) that:
   - Fetches commit SHAs from GitHub API
   - Fetches diffs and sends to Gemini for analysis
   - Emits analyzed commits via channel
4. `PersistCommits` stores results in PostgreSQL
5. `GET /subcommits-timeline` queries stored subcommits for a repo

### Port Interfaces

- `analysis.Agent` - AI agent for commit analysis (`internal/domain/analysis/agent_port.go`)
- `codehost.CodeHost` - Git hosting provider (`internal/domain/codehost/codehost_port.go`)
- `database.Database` - Data persistence (`internal/domain/database/database_port.go`)

### Error Handling

Use `internal/errors` package with category-based errors:
```go
errors.NewError(errors.ErrBadRequest, specificErr)
errors.NewError(errors.ErrInternalFailure, specificErr)
errors.NewError(errors.ErrNotFound, specificErr)
```

## API Endpoints

- `GET /` - Home (returns login status)
- `GET /auth/github/login` - Start OAuth flow
- `GET /auth/github/callback` - OAuth callback
- `POST /analyze` - Analyze a repository (authenticated)
- `GET /subcommits-timeline` - Get subcommits for a repo (authenticated)
