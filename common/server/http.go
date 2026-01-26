package server

import (
	"net/http"

	"go.uber.org/zap"
)

func RunHTTPServer(createHandler func(mux *http.ServeMux) http.Handler, logger *zap.Logger) {
	apiMux := http.NewServeMux()

	apiHandler := createHandler(apiMux)

	rootMux := http.NewServeMux()
	rootMux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	logger.Info("Starting HTTP server")

	err := http.ListenAndServe(":8080", rootMux)
	if err != nil {
		logger.Error("Unable to start http server")
	}
}
