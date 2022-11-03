package main

import (
	"context"
	"data-platform-api-business-partner-exconf-rmq-kube/config"
	"data-platform-api-business-partner-exconf-rmq-kube/database"
	"fmt"

	"github.com/latonaio/golang-logging-library/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
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

	rmq, err := rabbitmq.NewRabbitmqClient(c.RMQ.URL(), c.RMQ.QueueFrom(), "", nil, -1)
	if err != nil {
		l.Fatal(err.Error())
	}
	iter, err := rmq.Iterator()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Stop()
	for msg := range iter {
		go dataCheckProcess(ctx, c, db, msg)
	}
}

func dataCheckProcess(
	ctx context.Context,
	c *config.Conf,
	db *database.Mysql,
	rmqMsg rabbitmq.RabbitmqMessage,
) {
	defer rmqMsg.Success()
	l := logger.NewLogger()
	data := rmqMsg.Data()
	sessionId := getBodyHeader(data)
	l.AddHeaderInfo(map[string]interface{}{"runtime_session_id": sessionId})
	l.Info(rmqMsg.Data())

	exist := (*ExistencyChecker).Check(NewExistencyChecker(ctx, db, l), rmqMsg)
	rmqMsg.Respond(exist)
	l.Info(exist)
}

func getBodyHeader(data map[string]interface{}) string {
	id := fmt.Sprintf("%v", data["runtime_session_id"])
	return id
}
