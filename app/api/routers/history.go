package routers

import (
	"fmt"
	"net/http"

	"github.com/apsystole/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/skeremidchiev/gopher-translator-service/app/api"
	"github.com/skeremidchiev/gopher-translator-service/app/storage"
)

type historyRequest struct {
}

type historyResponse struct {
	api.APIResponse
	History string `json:"history"`
}

func handleHistoryRequest(s storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ret := fmt.Sprintf("Paniced, recovered value: %v\n", err)
				log.Errorln("[Routers] ", ret)

				render.Status(r, 500)
				render.JSON(w, r, historyResponse{api.APIResponse{Status: false, Error: ret}, ""})
			}
		}()

		history := s.GetAll()

		render.JSON(w, r, historyResponse{api.APIResponse{Status: true, Error: ""}, history})
	}
}

func NewHistoryRouter(s storage.Storage) http.Handler {
	r := chi.NewRouter()
	r.Get("/", handleHistoryRequest(s))
	return r
}
