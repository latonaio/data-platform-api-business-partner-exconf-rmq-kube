package main

import (
	"context"
	"data-platform-business-partner-exconf/config"
	"data-platform-business-partner-exconf/database"
	"data-platform-business-partner-exconf/input_reader"
	"data-platform-business-partner-exconf/output_creator"

	"github.com/latonaio/golang-logging-library/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client"
)

func main() {
	ctx := context.Background()
	l := logger.NewLogger()
	c := config.NewConf()
	db, err := database.NewMySQL(c.DB)
	if err != nil {
		l.Error(err)
		return
	}

	rmq, err := rabbitmq.NewRabbitmqClient(c.RMQ.URL(), c.RMQ.QueueFrom(), c.RMQ.QueueTo())
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Close()
	iter, err := rmq.Iterator()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Stop()
	for msg := range iter {
		dataCheckProcess(ctx, c, rmq, db, msg, l)
	}
}

func dataCheckProcess(
	ctx context.Context,
	c *config.Conf,
	rmq *rabbitmq.RabbitmqClient,
	db *database.Mysql,
	rmqMsg rabbitmq.RabbitmqMessage,
	l *logger.Logger,
) {
	defer rmqMsg.Success()
	data := rmqMsg.Data()
	l.Info(data)
	input, err := input_reader.ConvertToInput(data)
	if err != nil {
		l.Error("error: %+v", err)
		return
	}
	exist, err := ExistenceCheck(ctx, db, input.BusinessPartner.BusinessPartner)
	if err != nil {
		l.Info("error: %+v", err)
	}
	output := output_creator.ConvertToOutput(
		input.BusinessPartner.RuntimeSessionID,
		input.BusinessPartner.BusinessPartner,
		exist,
	)
	rmq.Send(c.RMQ.QueueTo()[0], map[string]interface{}{"BusinessPartnerExistence": output})
	l.Info(output)
}
