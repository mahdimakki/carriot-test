package main

import (
	"context"

	cfg "git.carriot.ir/warning-detector/cmd/config"
	"git.carriot.ir/warning-detector/internal/adapters/mqtt"
	"git.carriot.ir/warning-detector/internal/application"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	config, err := cfg.New()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	qu, err := mqtt.NewBroker(config.MQTT)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	app, err := application.New(qu, nil, nil, nil, nil)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	app.Do(ctx)
	w := make(chan bool)
	<-w
}
