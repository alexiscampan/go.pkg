// Package mqtt to define utils function for mqtt protocol
package mqtt

import (
	"context"
	"fmt"
	"time"

	"github.com/alexiscampan/go.pkg/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Consumer defines the consumer
type consumer struct {
	MQTTClientOptions *mqtt.ClientOptions
	MQTTClient        mqtt.Client
	QOS               byte
	Quiesce           uint
}

type brokerOptions struct {
	broker   string
	clientID string
	username string
	password string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	msg.Ack()
}

// Consumer interface
type Consumer interface {
	Connect(ctx context.Context) error
	Sub(context.Context) error
	Pub(ctx context.Context, message []byte) error
	Validate(ctx context.Context) error
	Close(ctx context.Context)
}

// NewConsumer to create a new consumer instance
func NewConsumer(params ...brokerOptions) (Consumer, error) {
	options := defaultOptions()
	if params != nil {
		options = setOptions(params[0])
	}
	return &consumer{
		MQTTClientOptions: options,
		QOS:               0,
		Quiesce:           0,
	}, nil
}

func setOptions(params brokerOptions) *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	options.AddBroker(params.broker)
	options.SetClientID(params.clientID)
	options.SetUsername(params.username)
	options.SetPassword(params.password)
	options.SetKeepAlive(7200 * time.Second)
	options.SetPingTimeout(10 * time.Second)
	options.SetDefaultPublishHandler(messagePubHandler)
	return options
}

func defaultOptions() *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	options.AddBroker("mqtt://test.mosquitto.org:1883")
	options.SetClientID("test1234")
	options.SetKeepAlive(7200 * time.Second)
	options.SetPingTimeout(10 * time.Second)
	options.SetDefaultPublishHandler(messagePubHandler)
	return options
}

func (c *consumer) Connect(ctx context.Context) error {
	c.MQTTClient = mqtt.NewClient(c.MQTTClientOptions)
	if token := c.MQTTClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect: %w", token.Error())
	}
	log.For(ctx).Info("Subscriber connected")
	return nil
}

func (c *consumer) Pub(ctx context.Context, message []byte) error {
	token := c.MQTTClient.Publish("topic/test", 0, false, message)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	time.Sleep(time.Second)
	return nil
}

func (c *consumer) Sub(ctx context.Context) error {
	topic := "topic/test"
	token := c.MQTTClient.Subscribe(topic, 1, nil)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Subscribed to topic: %s", topic)
	return nil
}

func (c *consumer) Close(ctx context.Context) {
	c.MQTTClient.Disconnect(c.Quiesce)
	log.For(ctx).Info("Consumer closed")
}

func (c *consumer) Validate(ctx context.Context) error {
	err := c.Connect(ctx)
	if err != nil {
		return nil
	}
	err = c.Sub(ctx)
	if err != nil {
		return nil
	}
	for i := 0; i < 10; i++ {
		err = c.Pub(ctx, []byte("test"))
		if err != nil {
			return err
		}
	}
	c.Close(ctx)
	return nil
}
