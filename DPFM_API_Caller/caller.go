package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-business-partner-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-business-partner-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-business-partner-exconf-rmq-kube/database"
	"sync"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
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

func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.BusinessPartnerGeneral {
	businessPartner := *input.BusinessPartnerGeneral.BusinessPartner
	notKeyExistence := make([]int, 0, 1)
	KeyExistence := make([]int, 0, 1)

	existData := &dpfm_api_output_formatter.BusinessPartnerGeneral{
		BusinessPartner: businessPartner,
		ExistenceConf:   false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !e.confBusinessPartnerGeneral(businessPartner) {
			notKeyExistence = append(notKeyExistence, businessPartner)
			return
		}
		KeyExistence = append(KeyExistence, businessPartner)
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData.ExistenceConf = true
	return existData
}

func (e *ExistenceConf) confBusinessPartnerGeneral(val int) bool {
	rows, err := e.db.Query(
		`SELECT BusinessPartner 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_business_partner_general_data 
		WHERE BusinessPartner = ?;`, val,
	)
	if err != nil {
		e.l.Error(err)
		return false
	}

	for rows.Next() {
		var businessPartner int
		err := rows.Scan(&businessPartner)
		if err != nil {
			e.l.Error(err)
			continue
		}
		if businessPartner == val {
			return true
		}
	}
	return false
}
