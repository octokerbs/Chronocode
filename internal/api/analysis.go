package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/octokerbs/chronocode-go/internal"
	"github.com/octokerbs/chronocode-go/internal/domain/agent/gemini"
	"github.com/octokerbs/chronocode-go/internal/domain/sourcecodehost/githubapi"
	"github.com/octokerbs/chronocode-go/internal/repository/supabase"
)

func AnalyzeRepositoryHandler(c *gin.Context) {
	accessToken := c.Query("access_token")
	repoURL := c.Query("repo_url")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	githubService, err := githubapi.NewGithubClient(ctx, accessToken, repoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	geminiGenerativeService, err := gemini.NewGeminiGenerativeService(ctx, os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer geminiGenerativeService.Close()

	supabaseService, err := supabase.NewSupabaseService(ctx, os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repositoryAnalyzer, err := internal.NewRepositoryAnalyzer(ctx, geminiGenerativeService, githubService, supabaseService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commits, errors := repositoryAnalyzer.AnalyzeRepository(ctx)

	advisory := ""
	if len(commits) > 20 {
		advisory = "Not all commits were analyzed due to repository analysis limit reached"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"message":          fmt.Sprintf("Successfully analyzed and stored %d commits", len(commits)),
		"analyses_count":   len(commits),
		"subcommits_count": 0,
		"time_taken":       0,
		"advisory":         advisory,
		"errors":           errors,
	})
}
