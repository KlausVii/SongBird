package main

import (
	"context"

	"time"

	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/KlausVii/SongBird/agent"
	log "github.com/cihub/seelog"
)

func main() {
	cli, top, err := agent.StartPubSub(context.Background(), agent.ProjectID, agent.Topic)
	if err != nil {
		panic(err)
	}
	publish(cli, top)
}

func publish(client *pubsub.Client, topic *pubsub.Topic) {
	ctx := context.Background()
	i := 0
	for {
		var msg string
		if i%2 == 0 {
			msg = "tweet " + strconv.Itoa(i)
		} else {
			msg = "toot " + strconv.Itoa(i)
		}
		_ = topic.Publish(ctx, &pubsub.Message{
			Data: []byte(msg),
		})
		log.Infof("tried to publish %v", msg)
		time.Sleep(5 * time.Second)
		i++
	}
}
