package mqtt

import (
	"context"
	"fmt"

	cfg "git.carriot.ir/warning-detector/cmd/config"
	"git.carriot.ir/warning-detector/internal/ports"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type LogBroker struct {
	client     mqtt.Client
	topic      string
	handlerFun func(context.Context, []byte) error
}

var _ ports.LogBroker = &LogBroker{}

func NewBroker(c cfg.MQTT) (*LogBroker, error) {
	client, err := queueClient(c)
	if err != nil {
		return nil, err
	}
	res := &LogBroker{
		client: client,
		topic:  c.Topic,
	}
	return res, nil
}

func queueClient(cfg cfg.MQTT) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", cfg.Host, cfg.Port))
	opts.SetUsername(cfg.User)
	opts.SetPassword(cfg.Pass)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func (l *LogBroker) Subscribe(ctx context.Context) {
	token := l.client.Subscribe(l.topic, 1, l.handleMessage)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", l.topic)
}

func (l *LogBroker) SetHandler(handler func(context.Context, []byte) error) {
	l.handlerFun = handler
}

func (l *LogBroker) handleMessage(client mqtt.Client, msg mqtt.Message) {
	ctx := context.Background() // TODO dirty
	err := l.handlerFun(ctx, msg.Payload())
	if err != nil {
		zap.L().Error(err.Error())
	}
}
