package appserver

import (
	"fmt"
	"net/http"

	"github.com/patrickisgreat/go-api-starter/internal/generated/go_api_starter/api"
)

func NewRouter(s api.TwirpServer) http.Handler {
	router := http.NewServeMux()
	router.Handle(s.PathPrefix(), s)
	return addHealthCheckRoute(router)
}

func addHealthCheckRoute(router *http.ServeMux) *http.ServeMux {
	router.Handle("/-/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "OK")
		if err != nil {
			http.Error(w, "health check failed", 500)
		}
	}))
	return router
}
