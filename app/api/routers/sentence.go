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

type sentenceRequest struct {
	Sentence string `json:"english-sentence"`
}

type sentenceResponse struct {
	api.APIResponse
	Sentence string `json:"gopher-sentence"`
}

func handleSentenceRequest(tr translater.Translater, s storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ret := fmt.Sprintf("Paniced, recovered value: %v\n", err)
				log.Errorln("[Routers] ", ret)

				render.Status(r, 500)
				render.JSON(w, r, sentenceResponse{api.APIResponse{Status: false, Error: ret}, ""})
			}
		}()

		var sentenceRequest *sentenceRequest

		err := json.NewDecoder(r.Body).Decode(&sentenceRequest)
		if err != nil {
			log.Errorln("[Routers] Error occurred while parsing request body: ", err)

			render.Status(r, 400)
			render.JSON(w, r, sentenceResponse{api.APIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		goSentence, err := s.GetTranslation(sentenceRequest.Sentence) // get from db
		if err == nil {                                               // if no error
			render.JSON(w, r, sentenceResponse{api.APIResponse{Status: true, Error: ""}, goSentence})
			return
		}

		goSentence, err = tr.ParseSentence(sentenceRequest.Sentence) // in case it isn't in base
		if err != nil {
			log.Errorln(err.Error())

			render.Status(r, 400)
			render.JSON(w, r, sentenceResponse{api.APIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		s.Save(sentenceRequest.Sentence, goSentence) // save in db

		render.JSON(w, r, sentenceResponse{api.APIResponse{Status: true, Error: ""}, goSentence})
	}
}

func NewSentenceRouter(tr translater.Translater, s storage.Storage) http.Handler {
	r := chi.NewRouter()
	r.Post("/", handleSentenceRequest(tr, s))
	return r
}
