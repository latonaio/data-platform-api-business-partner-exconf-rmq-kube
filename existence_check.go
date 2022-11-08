package main

import (
	"context"
	"data-platform-api-business-partner-exconf-rmq-kube/database"
	"data-platform-api-business-partner-exconf-rmq-kube/database/models"
	"data-platform-api-business-partner-exconf-rmq-kube/input_reader"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

type ExistencyChecker struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistencyChecker(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistencyChecker {
	return &ExistencyChecker{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistencyChecker) Check(data rabbitmq.RabbitmqMessage) map[string]interface{} {
	existData := map[string]interface{}{
		"ExistenceConf": false,
	}
	input := input_reader.Input{}
	err := json.Unmarshal(data.Raw(), &input)
	if err != nil {
		return existData
	}

	check := "BusinessPartner"
	bpID := input.BusinessPartner.BusinessPartner
	notExistKeys := make([]string, 0, 1)
	ExistKeys := make([]string, 0, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	existData[check] = bpID
	go func() {
		defer wg.Done()
		if !e.checkBusinessPartner(bpID) {
			notExistKeys = append(notExistKeys, check)
			return
		}
		ExistKeys = append(ExistKeys, check)
	}()

	wg.Wait()

	if len(ExistKeys) == 0 {
		return existData
	}
	if len(notExistKeys) > 0 {
		return existData
	}

	existData["ExistenceConf"] = true
	return existData
}

func (e *ExistencyChecker) checkBusinessPartner(val int) bool {
	start := time.Now()
	d, err := models.FindDataPlatformBusinessPartnerGeneralDatum(e.ctx, e.db, val)
	if d == nil || err != nil {
		return false
	}
	fmt.Printf("db check time %v", start)
	return d.BusinessPartner == val
}
