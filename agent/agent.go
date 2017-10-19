package agent

import (
	"github.com/gorilla/websocket"
	"github.com/murlokswarm/log"
)

type Agent interface {
	Login() ([]byte, error)
	IsAuth() bool
	Do(req string) ([]byte, error)
	Listen(c *websocket.Conn)
	Talk(c *websocket.Conn)
	Alive() bool
}

const maxBuffer = 512

type agentImpl struct {
	auth       bool
	alive      bool
	listenBuff chan []byte
	talkBuff   chan []byte
	msgCh      chan Msg
	killCh     chan struct{}
}

func NewAgent() Agent {
	a := &agentImpl{
		auth:       false,
		alive:      true,
		listenBuff: make(chan []byte, maxBuffer),
		talkBuff:   make(chan []byte, maxBuffer),
	}
	go a.handle()
	return a
}

func (a *agentImpl) Alive() bool {
	return a.alive
}

func (a *agentImpl) handle() {
	for {
		select {
		case msg := <-a.msgCh:
			a.talkBuff <- msg.Read()
			msg.Ack()
		case msg := <-a.listenBuff:
			switch {
			case a.IsAuth():
				ret, err := a.Do(string(msg))
				if err != nil {
					log.Error(err)
					return
				}
				a.talkBuff <- ret
			case string(msg) == "login":
				ret, err := a.Login()
				if err != nil {
					log.Error(err)
					return
				}
				a.talkBuff <- ret
			default:
				a.talkBuff <- []byte("Please login")
			}
		}
	}
}

func (a *agentImpl) Listen(c *websocket.Conn) {
	log.Info("listening")
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Errorf("read: %v", err)
			return
		}
		log.Infof("got msg: %s", string(message))
		a.listenBuff <- message
	}
}
func (a *agentImpl) Talk(c *websocket.Conn) {
	log.Info("talking")
	for {
		select {
		case msg := <-a.talkBuff:
			log.Infof("saying: %s", string(msg))
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error(err)
				return
			}
		}
	}
}

func (a *agentImpl) Login() ([]byte, error) {
	a.auth = true
	return []byte("succesful login"), nil
}

func (a *agentImpl) IsAuth() bool {
	return a.auth
}

func (a *agentImpl) Do(req string) ([]byte, error) {
	if req == "listen" {
		return a.startListen()
	}
	if req == "kill" {
		close(a.killCh)
		a.alive = false
		return []byte("dying"), nil
	}
	return append([]byte("I'm sorry I cannot do: "), []byte(req)...), nil
}

func (a *agentImpl) startListen() ([]byte, error) {
	var err error
	a.msgCh, a.killCh, err = ListenToTopic()
	if err != nil {
		return nil, err
	}
	return []byte("started listening"), nil
}
