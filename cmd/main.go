package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/skeremidchiev/gopher-translator-service/app/api"
	"github.com/skeremidchiev/gopher-translator-service/app/api/routers"
)

func main() {
	setupLogger()

	historyRouter := routers.NewHistoryRouter()
	sentenceRouter := routers.NewSentenceRouter()
	wordRouter := routers.NewWordRouter()

	a := api.NewAPI()
	a.AddRouter("/history", historyRouter)
	a.AddRouter("/sentence", sentenceRouter)
	a.AddRouter("/word", wordRouter)

	fs := NewFlagSet()
	err := fs.Init(os.Args[1:])
	if err != nil {
		log.Fatal("[Main] Error parsing comandline options: ", err)
	}

	a.Start(fs.PortNumber())
}
