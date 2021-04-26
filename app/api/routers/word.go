package routers

import (
	"encoding/json"
	"net/http"

	"github.com/apsystole/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/skeremidchiev/gopher-translator-service/app/api"
	"github.com/skeremidchiev/gopher-translator-service/app/storage"
	"github.com/skeremidchiev/gopher-translator-service/app/translater"
)

type wordRequest struct {
	word string `json:"english-word"`
}

type wordResponse struct {
	api.APIResponse
	word string `json:"gopher-word"`
}

func handleWordRequest(tr translater.Translater, s storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wordRequest *wordRequest

		err := json.NewDecoder(r.Body).Decode(&wordRequest)
		if err != nil {
			log.Errorln("[Routers] Error occurred while parsing request body: ", err)
			render.JSON(w, r, wordResponse{api.APIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		// TODO:
		// goWord, err :=

		// if err != nil {
		// 	log.Errorln(err.Error())
		// 	render.JSON(w, r, wordResponse{api.APIResponse{Status: false, Error: err.Error()}, nil})
		// 	return
		// }

		// render.JSON(w, r, wordResponse{api.APIResponse{Status: true, Error: ""}, goWord})
	}
}

func NewWordRouter(tr translater.Translater, s storage.Storage) http.Handler {
	r := chi.NewRouter()
	r.Post("/", handleWordRequest(tr, s))
	return r
}
