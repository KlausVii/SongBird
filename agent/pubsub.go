package agent

import (
	"context"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	log "github.com/cihub/seelog"
)

const (
	ProjectID = "songbird"
	Topic     = "songbird"
	sub       = "subsub"
)

type Msg struct {
	msg *pubsub.Message
}

func (m *Msg) Ack() {
	m.msg.Ack()
}

func (m *Msg) Read() []byte {
	return m.msg.Data
}

func ListenToTopic() (chan Msg, chan struct{}, error) {
	ctx := context.Background()
	log.Info("starting pubsub")
	c, t, err := StartPubSub(ctx, ProjectID, Topic)
	if err != nil {
		return nil, nil, err
	}
	kill := make(chan struct{})
	log.Info("starting listen")
	ch, err := listenToTopic(ctx, c, t, kill)
	if err != nil {
		return nil, nil, err
	}
	log.Info("listen ready")
	return ch, kill, nil
}

func StartPubSub(ctx context.Context, projectID, topicName string) (*pubsub.Client, *pubsub.Topic, error) {
	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Errorf("Failed to create client: %v", err)
		return nil, nil, err
	}

	topic := client.Topic(topicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		log.Errorf("Failed to create topic: %v", err)
		return nil, nil, err
	}
	if ok {
		log.Info("topic exists")
		return client, topic, nil
	}
	log.Info("topic does not exist")
	// Creates the new topic.
	topic, err = client.CreateTopic(ctx, topicName)
	if err != nil {
		log.Errorf("Failed to create topic: %v", err)
		return nil, nil, err
	}

	log.Infof("Topic %v created.\n", topic)
	return client, topic, nil
}

func createSub(client *pubsub.Client, name string, topic *pubsub.Topic) error {
	ctx := context.Background()

	sub := client.Subscription(name)
	ok, err := sub.Exists(ctx)
	if err != nil {
		log.Errorf("Failed to create topic: %v", err)
		return err
	}
	if ok {
		return nil
	}
	// [START create_subscription]
	sub, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return err
	}
	log.Infof("Created subscription: %v\n", sub)
	// [END create_subscription]
	return nil
}

func listenToTopic(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, kill chan struct{}) (chan Msg, error) {
	subName := "sub" + strconv.Itoa(int(time.Now().Unix()))
	if err := createSub(client, subName, topic); err != nil {
		return nil, err
	}
	var mu sync.Mutex
	//received := 0
	sub := client.Subscription(subName)
	//cctx, cancel := context.WithCancel(ctx)
	ch := make(chan Msg, maxBuffer)
	go func() {
		select {
		case <-kill:
			log.Info(sub.Delete(ctx))
			return
		default:
		}
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			select {
			case <-kill:
				return
			default:
			}
			mu.Lock()
			defer mu.Unlock()
			log.Infof("Got message: %v", string(msg.Data))
			ch <- Msg{
				msg: msg,
			}
		})
		if err != nil {
			panic(err)
		}
	}()
	return ch, nil
}
