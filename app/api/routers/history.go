package routers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/skeremidchiev/gopher-translator-service/app/api"
)

type historyRequest struct {
}

type historyResponse struct {
	api.APIResponse
	history string `json:"history"`
}

func handleHistoryRequest() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

func NewHistoryRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", handleHistoryRequest())
	return r
}
