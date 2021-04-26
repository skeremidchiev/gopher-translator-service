package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apsystole/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/skeremidchiev/gopher-translator-service/app/api"
	"github.com/skeremidchiev/gopher-translator-service/app/storage"
	"github.com/skeremidchiev/gopher-translator-service/app/translater"
)

type wordRequest struct {
	Word string `json:"english-word"`
}

type wordResponse struct {
	api.APIResponse
	Word string `json:"gopher-word,omitempty"`
}

func handleWordRequest(tr translater.Translater, s storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				ret := fmt.Sprintf("Paniced, recovered value: %v\n", err)
				log.Errorln("[Routers] ", ret)

				render.Status(r, 500)
				render.JSON(w, r, wordResponse{api.APIResponse{Status: false, Error: ret}, ""})
			}
		}()

		var wordRequest *wordRequest

		err := json.NewDecoder(r.Body).Decode(&wordRequest)
		if err != nil {
			log.Errorln("[Routers] Error occurred while parsing request body: ", err)

			render.Status(r, 400)
			render.JSON(w, r, wordResponse{api.APIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		goWord, err := s.GetTranslation(wordRequest.Word) // get from db
		if err == nil {                                   // if no error
			render.JSON(w, r, wordResponse{api.APIResponse{Status: true}, goWord})
			return
		}

		goWord, err = tr.ParseWord(wordRequest.Word)
		if err != nil {
			log.Errorln(err.Error())
			render.Status(r, 400)
			render.JSON(w, r, wordResponse{api.APIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		s.Save(wordRequest.Word, goWord) // save in db

		render.JSON(w, r, wordResponse{api.APIResponse{Status: true}, goWord})
	}
}

func NewWordRouter(tr translater.Translater, s storage.Storage) http.Handler {
	r := chi.NewRouter()
	r.Post("/", handleWordRequest(tr, s))
	return r
}
