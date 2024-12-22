package main

import (
	"context"
	"log"
)

func main() {
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	webServerApp, err := wireLeaderBoardApiApplication()
	if err != nil {
		log.Fatal(err)
	}
	if err := webServerApp.Start(mainCtx); err != nil {
		log.Fatal(err)
	}
}
