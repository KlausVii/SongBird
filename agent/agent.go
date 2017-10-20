package agent

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/KlausVii/SongBird/consts"
	"github.com/KlausVii/SongBird/proto"
	ps "github.com/KlausVii/SongBird/pubsub"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/murlokswarm/log"
	"github.com/pkg/errors"
)

type Agent interface {
	Start(c *websocket.Conn) (chan struct{}, error)
	Stop()
}

const maxBuffer = 512

type agentImpl struct {
	conn       *websocket.Conn
	psClient   *pubsub.Client
	auth       bool
	killCh     chan struct{}
	alive      chan struct{}
	listenBuff chan []byte
	talkBuff   chan []byte
}

func NewAgent() Agent {
	a := &agentImpl{
		auth:       false,
		listenBuff: make(chan []byte, maxBuffer),
		talkBuff:   make(chan []byte, maxBuffer),
	}
	return a
}

func (a *agentImpl) Start(c *websocket.Conn) (chan struct{}, error) {
	if c == nil {
		return nil, errors.New("nil connection")
	}
	a.killCh = make(chan struct{})
	a.alive = make(chan struct{})
	a.conn = c
	var err error
	a.psClient, err = pubsub.NewClient(context.Background(), consts.ProjectID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	go a.handle()
	go a.listen()
	go a.talk()
	return a.alive, nil
}

func (a *agentImpl) Stop() {
	if err := a.conn.Close(); err != nil {
		log.Error(err)
	}
	return
}

func (a *agentImpl) shutdown() {
	a.talkBuff <- []byte("shutting down")
	close(a.talkBuff)
	close(a.listenBuff)
	close(a.killCh)
	close(a.alive)

}
func (a *agentImpl) handle() {
	for {
		select {
		case msg, ok := <-a.listenBuff:
			if !ok {
				log.Info("stopped handling")
				return
			}
			switch {
			case a.isAuth():
				a.do(string(msg))
			case string(msg) == "login":
				a.login()
			default:
				a.talkBuff <- []byte("Please login")
			}
		}
	}
}

func (a *agentImpl) listen() {
	log.Info("listening")
	for {
		_, message, err := a.conn.ReadMessage()
		if err != nil {
			log.Errorf("read: %v", err)
			return
		}
		log.Infof("got msg: %s", string(message))
		a.listenBuff <- message
	}
}
func (a *agentImpl) talk() {
	log.Info("talking")
	for {
		select {
		case msg, ok := <-a.talkBuff:
			if !ok {
				log.Info("stopped talking")
				return
			}
			if err := a.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error(err)
				return
			}
		}
	}
}

func (a *agentImpl) login() {
	a.auth = true
	a.talkBuff <- []byte("successful login")
	return
}

func (a *agentImpl) isAuth() bool {
	return a.auth
}

func (a *agentImpl) do(req string) {
	switch req {
	case "listen":
		a.startListen()
	case "kill":
		a.shutdown()
	case "req":
		a.makeReq()
	default:
		a.talkBuff <- append([]byte("I'm sorry I cannot do: "), []byte(req)...)
	}
}

func (a *agentImpl) makeReq() {
	log.Info("making req")
	reqTopicName := "req"
	msgCh, kCh, err := ps.ListenTopicWithClient(context.Background(), reqTopicName, a.psClient)
	if err != nil {
		log.Error(err)
		a.talkBuff <- []byte(err.Error())
		return
	}
	go func() {
		select {
		case m := <-msgCh:
			log.Info("received response")
			a.talkBuff <- m.Read()
			m.Ack()
			close(kCh)
			return
		}
	}()
	req := songbird.Request{
		Topic: reqTopicName,
		Req:   "hello world",
	}
	reqByte, err := proto.Marshal(&req)
	if err != nil {
		log.Error(err)
		a.talkBuff <- []byte(err.Error())
		return
	}
	if err := ps.PublishToExistingTopic(context.Background(), a.psClient, consts.APITopic, reqByte); err != nil {
		log.Error(err)
		a.talkBuff <- []byte(err.Error())
		return
	}
	a.talkBuff <- []byte("made request, awaiting reply")
}

func (a *agentImpl) startListen() {
	go func() {
		msgCh, killCh, err := ps.ListenTopicWithClient(context.Background(), consts.APITopic, a.psClient)
		if err != nil {
			a.talkBuff <- []byte(err.Error())
			return
		}
		for {
			select {
			case <-a.killCh:
				log.Info("stopped pub sub listen")
				close(killCh)
				return
			case m := <-msgCh:
				a.talkBuff <- m.Read()
				m.Ack()
			}
		}
	}()
	a.talkBuff <- []byte("started listening api topic")
	return
}
