package main

import (
	"fmt"
	"os"
	"template/models"

	"github.com/go-kit/kit/log"
	"template/application"
	"template/prometheus"
	"template/provider"
	"template/repository"
	"template/service"
)

var (
	appConfig models.Config
	pr        *prometheus.Prometheus
	logger    log.Logger
)

func init() {
	models.LoadConfig(&appConfig)
	fmt.Println(appConfig.NatsServer)
}

func main() {
	pr = prometheus.New("new-test")
	logger = log.With(
		log.NewJSONLogger(os.Stderr),
		"caller", log.DefaultCaller,
	)
	logger = prometheus.NewLogger(logger, pr)

	natsStreaming := provider.NewNATS()
	err := natsStreaming.ConnectNatsStreaming(appConfig.NatsServer.Address)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	repNatsStreaming := repository.New(natsStreaming)

	svc := service.New(repNatsStreaming)

	app := application.New(&application.Options{
		Serv:    appConfig.ServerOpt,
		HashSum: appConfig.HashSum,
		Svc:     svc,
		Pr:      pr,
		Logger:  logger,
	})

	app.Start()
}
