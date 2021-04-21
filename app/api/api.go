package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

type APIRoute http.Handler

type API struct {
	r *chi.Mux
}

type APIResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (api *API) AddRouter(route string, router APIRoute) {
	api.r.Mount(route, router)
}

func (api *API) Start(port string) error {
	log.Infof("[API] Listening for REST API Requests on port %v\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), api.r)
}

func NewAPI() *API {
	r := chi.NewRouter()
	r.Use(
		middleware.AllowContentType("application/json"),
		middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.StandardLogger()}),
		middleware.Timeout(60*time.Second),
	)

	return &API{r: r}
}
