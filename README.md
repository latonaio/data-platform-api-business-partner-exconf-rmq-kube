# data-platform-api-business-partner-exconf-rmq-kube
data-platform-api-business-partner-exconf-rmq-kube は、データ連携基盤において、API でビジネスパートナの存在性チェックを行うためのマイクロサービスです。

## 動作環境
・ OS: LinuxOS  
・ CPU: ARM/AMD/Intel  

## 存在確認先テーブル名
以下のsqlファイルに対して、ビジネスパートナの存在確認が行われます。

* data-platform-business-partner-general-data.sql（データ連携基盤 ビジネスパートナ - 一般データ）

## existence_check.go による存在性確認
Input で取得されたファイルに基づいて、existence_check.go で、 API がコールされます。
existence_check.go の 以下の箇所が、指定された API をコールするソースコードです。

```
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
```

## Input
data-platform-api-business-partner-exconf-rmq-kube では、以下のInputファイルをRabbitMQからJSON形式で受け取ります。  

```
{
	"connection_key": "request",
	"result": true,
	"redis_key": "abcdefg",
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"business_partner": 201,
	"filepath": "/var/lib/aion/Data/rededge_sdc/abcdef.json",
	"service_label": "ORDERS",
	"BusinessPartner": {
		"BusinessPartner": 101
	},
	"api_schema": "DPFMOrdersCreates",
	"accepter": ["All"],
	"order_id": null,
	"deleted": false
}
```

## Output
data-platform-api-business-partner-exconf-rmq-kube では、[golang-logging-library-for-data-platform](https://github.com/latonaio/golang-logging-library-for-data-platform) により、Output として、RabbitMQ へのメッセージを JSON 形式で出力します。ビジネスパートナの対象値が存在する場合 true、存在しない場合 false、を返します。"cursor" ～ "time"は、golang-logging-library-for-data-platform による 定型フォーマットの出力結果です。

```
{
	"cursor": "/go/src/github.com/latonaio/existence_check/checker.go#L116",
	"function": "data-platform-api-orders-creates-rmq-kube/existence_check.(*ExistenceChecker).bpExistenceCheck",
	"level": "INFO",
	"message": {
		"BusinessPartner": 201,
		"ExistenceConf": true
	},
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"time": "2022-11-08T07:50:59Z"
}
```