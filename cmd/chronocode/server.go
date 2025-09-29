package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/octokerbs/chronocode-go/internal/api"
)

type Server struct {
	anEngine *gin.Engine
	aPort    string
}

func NewServer(aPort string) Server {
	engine := gin.Default()
	log.SetFlags(log.LstdFlags)

	engine.GET("/analyze-repository", api.AnalyzeRepositoryHandler)

	return Server{engine, aPort}
}

func (s Server) Run() {
	s.anEngine.Run(s.aPort)
}
