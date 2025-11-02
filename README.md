<div align="center">
  
# 📅 Chronocode 
An intelligent GitHub repo analyzer that summarizes and categorizes commits using Gemini, GitHub API, and Supabase. This iteration is implemented in golang with worker pools for commit downloading, concurrent AI processing and database pushing. 

This is a port of a project made for the [ShipBA Hackaton 2025](https://www.shipba.dev/) with [Octavio Pavón](https://x.com/octaviopvn1) and [Tiago Prelato](https://x.com/SneyX_). Originally made with Python and Lovable.
</div>

## Frontend example
![Demo](assets/Demo.png)

## ENV file setup
```env
GEMINI_API_KEY=
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=
```

## Run
```bash
docker compose up
```

## Example
```
http://localhost:8080/analyze-repository?repo_url=https://github.com/octokerbs/50Cent-Dolar-Blue-Bot
Bearer token with your github token.
```

## Notes
- Google gives free gemini api keys for testing here [Google AI Studio](https://aistudio.google.com/app/apikey).
- The performance directory is just for testing the speed of the workers pool.
- The frontend is not currently available. 
- We use supabase to store commits and not re-analyze them when a repo is asked to again.

## Notes for me
```bash
docker compose up 
```

```bash
docker compose down
```

```bash
docker compose build
```