package routers

import (
	"encoding/json"
	"net/http"

	"github.com/apsystole/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/skeremidchiev/gopher-translator-service/app/api"
)

type sentenceRequest struct {
	sentence string `json:"english-sentence"`
}

type sentenceResponse struct {
	api.APIResponse
	sentence string `json:"gopher-sentence"`
}

func handleSentenceRequest() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var sentenceRequest *sentenceRequest

		err := json.NewDecoder(r.Body).Decode(&sentenceRequest)
		if err != nil {
			log.Errorln("[Routers] Error occurred while parsing request body: ", err)
			render.JSON(w, r, sentenceResponse{api.APIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		// TODO:
		// goSentence, err :=

		// if err != nil {
		// 	log.Errorln(err.Error())
		// 	render.JSON(w, r, sentenceResponse{api.APIResponse{Status: false, Error: err.Error()}, nil})
		// 	return
		// }

		// render.JSON(w, r, sentenceResponse{api.APIResponse{Status: true, Error: ""}, goSentence})
	}
}

func NewSentenceRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", handleSentenceRequest())
	return r
}
