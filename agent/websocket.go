package agent

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/murlokswarm/log"
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
		log.Errorf("welcome: %v", err)
		return
	}
	a := NewAgent()
	ch, err := a.Start(c)
	if err != nil {
		log.Errorf("start failed: %v", err)
		return
	}
	<-ch
	a.Stop()
	log.Infof("session closed")
}
