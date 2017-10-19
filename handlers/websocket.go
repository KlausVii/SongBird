package handlers

import (
	"net/http"

	"github.com/KlausVii/SongBird/agent"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	up := websocket.Upgrader{}
	c, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade: %v", err)
		return
	}
	defer c.Close()
	if err = c.WriteMessage(websocket.TextMessage, []byte("welcome!")); err != nil {
		log.Error("welcome: %v", err)
		return
	}
	a := agent.NewAgent()
	go a.Listen(c)
	go a.Talk(c)
	for a.Alive() {
	}
}
