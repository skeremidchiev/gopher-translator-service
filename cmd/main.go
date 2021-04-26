package main

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/skeremidchiev/gopher-translator-service/app/api"
	"github.com/skeremidchiev/gopher-translator-service/app/api/routers"
	"github.com/skeremidchiev/gopher-translator-service/app/storage"
	"github.com/skeremidchiev/gopher-translator-service/app/translater"
)

func main() {
	godotenv.Load()
	setupLogger()

	fs := NewFlagSet()
	err := fs.Init(os.Args[1:])
	if err != nil {
		log.Fatal("[Main] Error parsing comandline options: ", err)
	}

	config := getConfiguration()
	storage := storage.NewStorage()
	translator := translater.NewTranslater(config)

	historyRouter := routers.NewHistoryRouter(storage)
	sentenceRouter := routers.NewSentenceRouter(translator, storage)
	wordRouter := routers.NewWordRouter(translator, storage)

	a := api.NewAPI()
	a.AddRouter("/history", historyRouter)
	a.AddRouter("/sentence", sentenceRouter)
	a.AddRouter("/word", wordRouter)
	a.Start(fs.PortNumber())
}
