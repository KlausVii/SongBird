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
	alive      chan struct{}
	auth       bool
	listenBuff chan []byte
	talkBuff   chan []byte
	msgCh      chan ps.Msg
	killCh     chan struct{}
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
	a.alive = make(chan struct{})
	a.killCh = make(chan struct{})
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
	// todo
	return
}

func (a *agentImpl) handle() {
	for {
		select {
		case msg := <-a.msgCh:
			a.talkBuff <- msg.Read()
			msg.Ack()
		case msg := <-a.listenBuff:
			switch {
			case a.isAuth():
				ret, err := a.do(string(msg))
				if err != nil {
					log.Error(err)
					return
				}
				a.talkBuff <- ret
			case string(msg) == "login":
				ret, err := a.login()
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
		case msg := <-a.talkBuff:
			log.Infof("saying: %s", string(msg))
			if err := a.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error(err)
				return
			}
		}
	}
}

func (a *agentImpl) login() ([]byte, error) {
	a.auth = true
	return []byte("successful login"), nil
}

func (a *agentImpl) isAuth() bool {
	return a.auth
}

func (a *agentImpl) do(req string) ([]byte, error) {
	switch req {
	case "listen":
		return a.startListen()
	case "kill":
		close(a.killCh)
		close(a.alive)
		return []byte("dying"), nil
	case "req":
		return a.makeReq()
	default:
	}
	return append([]byte("I'm sorry I cannot do: "), []byte(req)...), nil
}

func (a *agentImpl) makeReq() ([]byte, error) {
	log.Info("making req")
	reqTopicName := "req"
	msgCh, kCh, err := ps.ListenTopicWithClient(context.Background(), reqTopicName, a.psClient)
	if err != nil {
		log.Error(err)
		return nil, err
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
		return nil, err
	}
	if err := ps.PublishToExistingTopic(context.Background(), a.psClient, consts.APITopic, reqByte); err != nil {
		log.Error(err)
		return nil, err
	}
	return []byte("made request, awaiting reply"), nil
}

func (a *agentImpl) startListen() ([]byte, error) {
	var err error
	a.msgCh, a.killCh, err = ps.ListenTopicWithClient(context.Background(), consts.APITopic, a.psClient)
	if err != nil {
		return nil, err
	}
	return []byte("started listening api topic"), nil
}
