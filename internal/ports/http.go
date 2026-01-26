package ports

import (
	"net/http"

	"github.com/octokerbs/chronocode-backend/internal/app"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) GetSubcommits(w http.ResponseWriter, r *http.Request) {

}
