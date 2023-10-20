package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Christian-Bull/slack-reactbot/handlers"
	"github.com/Christian-Bull/slack-reactbot/util"
)

func main() {

	l := log.New(os.Stdout, "walabot", log.LstdFlags)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	bindAddr := fmt.Sprintf(":%s", port)

	// post test connection message
	err := util.PostMessage(
		l,
		util.CreateMessage("Connected :wala:", os.Getenv("LOGCHANNEL")),
	)
	if err != "" {
		l.Fatal("Error posting slack message: ", err)
	}

	// Create and register handlers
	sh := handlers.NewSlack(l, os.Getenv("SLACKAPIKEY"))

	http.Handle("/slack", sh)

	l.Printf("Starting server on port %s", port)
	l.Fatal(http.ListenAndServe(bindAddr, nil))
}
