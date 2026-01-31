# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Chronocode is an intelligent GitHub repository analyzer that summarizes and categorizes commits using Gemini AI, GitHub API, and PostgreSQL. Built with Go 1.23+, it uses worker pools for concurrent commit downloading, AI processing, and database operations.

This is a port of a project originally built with Python for the ShipBA Hackaton 2025.

## Development Commands

### Running the Application

```bash
docker-compose up --build
```

This starts:
- PostgreSQL database (port 5460:5432)
- Database migrator (runs migrations automatically)
- Go backend application (port 8080)

### Local Development

```bash
# Build the application
go build -o main ./internal/main.go

# Run tests
go test ./...

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### Database Migrations

Migrations are in the `migrations/` directory and run automatically via the `migrator` service in docker-compose. The project uses [migrate/migrate](https://github.com/golang-migrate/migrate) for database migrations.

Migration files follow the pattern: `{version}_{description}.{up|down}.sql`

## Architecture

The project follows **Domain-Driven Design (DDD)** with **Hexagonal Architecture** and **CQRS** patterns.

### Directory Structure

```
internal/
├── domain/repository/    # Domain entities (Repo, Commit, Subcommit, CommitAnalysis)
├── app/                  # Application layer (CQRS)
│   ├── command/         # Write operations (AnalyzeRepo)
│   └── query/           # Read operations (IsRepoAnalyzed, RepoSubcommits)
├── adapters/            # Infrastructure implementations
│   ├── repo_postgresql_repository.go  # PostgreSQL adapter
│   ├── codehost_github.go            # GitHub API adapter
│   ├── agent_gemini.go               # Gemini AI adapter
│   └── auth_github.go                # GitHub OAuth
├── ports/               # Entry points (HTTP handlers)
└── service/             # Dependency injection & app initialization

common/                  # Shared utilities
├── decorator/          # Command/Query decorators (logging, etc.)
├── logs/              # Zap logger setup
└── server/            # HTTP server utilities
```

### Key Architectural Patterns

#### 1. Hexagonal Architecture (Ports & Adapters)

**Domain** (core business logic) is isolated from external concerns through interfaces:
- **Ports**: Interfaces defined in [internal/app/command/analyze_repo.go](internal/app/command/analyze_repo.go) (e.g., `Agent`, `CodeHost`, `CodeHostFactory`)
- **Adapters**: Concrete implementations in [internal/adapters/](internal/adapters/) (e.g., `GeminiAgent`, `GitHubCodeHost`)

#### 2. CQRS (Command Query Responsibility Segregation)

Operations are split into:
- **Commands**: Write operations that change state (e.g., `AnalyzeRepo`)
- **Queries**: Read operations that return data (e.g., `IsRepoAnalyzed`, `RepoSubcommits`)

Both use the **Handler** pattern:
```go
type CommandHandler[C any] interface {
    Handle(ctx context.Context, cmd C) error
}

type QueryHandler[Q any, R any] interface {
    Handle(ctx context.Context, q Q) (R, error)
}
```

#### 3. Decorator Pattern

All handlers are wrapped with decorators for cross-cutting concerns (logging, metrics, etc.) via [common/decorator/](common/decorator/). See `ApplyCommandDecorators` and `ApplyQueryDecorators`.

### Worker Pool Architecture

The `AnalyzeRepo` command ([internal/app/command/analyze_repo.go](internal/app/command/analyze_repo.go)) uses **concurrent worker pools**:

1. **Commit Producer**: Fetches commit SHAs from GitHub and sends to channel
2. **Analyzer Workers** (5 concurrent):
   - Fetch commit diffs from GitHub
   - Analyze with Gemini AI
   - Fetch full commit data
   - Send analyzed commits to output channel

Uses Go channels for communication between workers and `sync.WaitGroup` for coordination.

## Domain Model

Core entities in [internal/domain/repository/](internal/domain/repository/):
- **Repo**: GitHub repository metadata
- **Commit**: Git commit with analysis metadata
- **Subcommit**: Categorized chunks of a commit (feature, bugfix, refactor, etc.)
- **CommitAnalysis**: AI-generated analysis result

## Dependencies & Integration

### External Services
- **GitHub API**: Repository and commit data via [google/go-github](https://github.com/google/go-github)
- **Gemini AI**: Commit diff analysis via [google/generative-ai-go](https://github.com/google/generative-ai-go)
- **PostgreSQL**: Data persistence via [lib/pq](https://github.com/lib/pq)
- **GitHub OAuth**: Authentication flow

### Required Environment Variables

```env
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
DATABASE_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
GEMINI_API_KEY=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
REDIRECT_URL=http://localhost:8080/auth/github/callback
```

## Application Entry Point

The application bootstraps in [internal/main.go](internal/main.go):
1. Initialize logger ([common/logs/zap.go](common/logs/zap.go))
2. Create application context
3. Wire dependencies via `service.NewApplication` (dependency injection)
4. Start HTTP server with routes

## Important Notes

- The Dockerfile references `./cmd/server/main.go` but the actual main is in `./internal/main.go` - this may need correction
- HTTP handlers in [internal/ports/http.go](internal/ports/http.go) are currently stub implementations
- Worker pool size is hardcoded to 5 workers in [internal/app/command/analyze_repo.go:45](internal/app/command/analyze_repo.go#L45)
