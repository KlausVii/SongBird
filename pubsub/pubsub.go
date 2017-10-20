package pubsub

import (
	"context"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/murlokswarm/log"
	"github.com/pkg/errors"
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

func ListenTopicWithClient(ctx context.Context, topicName string, c *pubsub.Client) (chan Msg, chan struct{}, error) {
	kill := make(chan struct{})
	log.Info("starting listen")
	topic, err := createTopic(ctx, topicName, c)
	if err != nil {
		return nil, nil, err
	}
	ch, err := listenToTopic(ctx, c, topic, kill)
	if err != nil {
		return nil, nil, err
	}
	log.Info("listen ready")
	return ch, kill, nil
}

func CreatePublishToTopic(ctx context.Context, c *pubsub.Client, t string, msg []byte) error {
	topic, err := createTopic(ctx, t, c)
	if err != nil {
		return err
	}
	_ = topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	log.Infof("tried to publish %v", msg)
	return nil
}

func PublishToExistingTopic(ctx context.Context, c *pubsub.Client, t string, msg []byte) error {
	topic, err := getTopicIfExists(ctx, t, c)
	if err != nil {
		return err
	}
	_ = topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	log.Infof("tried to publish %v", msg)
	return nil
}

func createTopic(ctx context.Context, topicName string, c *pubsub.Client) (*pubsub.Topic, error) {
	topic, err := getTopicIfExists(ctx, topicName, c)
	if err != nil {
		log.Info("topic does not exist")
		// Creates the new topic.
		return c.CreateTopic(ctx, topicName)
	}
	return topic, nil
}

func getTopicIfExists(ctx context.Context, topicName string, c *pubsub.Client) (*pubsub.Topic, error) {
	topic := c.Topic(topicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		log.Errorf("Failed to check topic: %v", err)
		return nil, err
	}
	if !ok {
		log.Info("topic does not exist")
		return nil, errors.New("topic does not exist")
	}
	return topic, nil
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
	ch := make(chan Msg)
	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			select {
			case <-kill:
				log.Info(sub.Delete(ctx))
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
