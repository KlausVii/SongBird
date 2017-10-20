package main

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/KlausVii/SongBird/consts"
	"github.com/KlausVii/SongBird/proto"
	ps "github.com/KlausVii/SongBird/pubsub"
	"github.com/golang/protobuf/proto"
	"github.com/murlokswarm/log"
)

func main() {
	c, err := pubsub.NewClient(context.Background(), consts.ProjectID)
	if err != nil {
		log.Error(err)
		return
	}
	msgCh, killCh, err := ps.ListenTopicWithClient(context.Background(), consts.APITopic, c)
	if err != nil {
		log.Error(err)
		return
	}
	handle(msgCh, killCh, c)
	log.Info("closing server")
}

func handle(msgCh chan ps.Msg, killCh chan struct{}, c *pubsub.Client) {
	for {
		select {
		case msg := <-msgCh:
			req := songbird.Request{}
			if err := proto.Unmarshal(msg.Read(), &req); err != nil {
				log.Errorf("failed to unmarshal req: %v", err)
			}
			if err := handler(req, c); err != nil {
				log.Error(err)
			}
		}
	}
}

func handler(req songbird.Request, c *pubsub.Client) error {
	log.Infof("received req: %v", req.GetReq())
	log.Infof("responding on topic: %v", req.GetTopic())
	rsp := songbird.Response{
		Rsp: "message handled",
	}
	m, err := proto.Marshal(&rsp)
	if err != nil {
		return err
	}
	return ps.PublishToExistingTopic(context.Background(), c, req.GetTopic(), m)
}
