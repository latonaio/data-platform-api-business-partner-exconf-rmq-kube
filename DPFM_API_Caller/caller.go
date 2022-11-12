package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-business-partner-exconf-rmq-kube/DPFM_API_Input_Reader"
	"data-platform-api-business-partner-exconf-rmq-kube/database"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(data rabbitmq.RabbitmqMessage) map[string]interface{} {
	existData := map[string]interface{}{
		"ExistenceConf": false,
	}
	input := dpfm_api_input_reader.SDC{}
	err := json.Unmarshal(data.Raw(), &input)
	if err != nil {
		return existData
	}

	conf := "BusinessPartner"
	businessPartnerID := *input.BusinessPartnerID.BusinessPartnerID
	notKeyExistence := make([]int, 0, 1)
	KeyExistence := make([]int, 0, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	existData[conf] = businessPartnerID
	go func() {
		defer wg.Done()
		if !e.confBusinessPartnerGeneral(businessPartnerID) {
			notKeyExistence = append(notKeyExistence, businessPartnerID)
			return
		}
		KeyExistence = append(KeyExistence, businessPartnerID)
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData["ExistenceConf"] = true
	return existData
}

func (e *ExistenceConf) confBusinessPartnerGeneral(val int) bool {
	rows, err := e.db.Query(
		`SELECT BusinessPartner 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_business_partner_general_data 
		WHERE BusinessPartner = ?;`, val,
	)
	if err != nil {
		fmt.Printf("err = %+v \n", err)
		return false
	}

	for rows.Next() {
		var businessPartner int
		err := rows.Scan(&businessPartner)
		if err != nil {
			fmt.Printf("err = %+v \n", err)
			continue
		}
		fmt.Printf("data = %+v \n", businessPartner)
		if businessPartner == val {
			return true
		}
	}
	return false
}
