package ports

import "github.com/octokerbs/chronocode-backend/internal/app"

type contextKey string

const githubTokenKey contextKey = "github_token"

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{app: application}
}

func (h HttpServer) RunHTTPServer() {

}
