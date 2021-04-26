package routers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/skeremidchiev/gopher-translator-service/app/api"
	"github.com/skeremidchiev/gopher-translator-service/app/storage"
)

type historyRequest struct {
}

type historyResponse struct {
	api.APIResponse
	history string `json:"history"`
}

func handleHistoryRequest(s storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Handle panic

		// TODO:
		// history, err :=

		// if err != nil {
		// 	log.Errorln(err.Error())
		// 	render.JSON(w, r, historyResponse{api.APIResponse{Status: false, Error: err.Error()}, nil})
		// 	return
		// }

		// render.JSON(w, r, historyResponse{api.APIResponse{Status: true, Error: ""}, history})
	}
}

func NewHistoryRouter(s storage.Storage) http.Handler {
	r := chi.NewRouter()
	r.Get("/", handleHistoryRequest(s))
	return r
}
