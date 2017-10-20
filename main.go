package main

import (
	"net/http"

	"github.com/KlausVii/SongBird/agent"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", agent.HandleWebSocket)
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
