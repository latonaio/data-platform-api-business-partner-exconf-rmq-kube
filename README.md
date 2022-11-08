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
func ExistenceCheck(ctx context.Context, db *database.Mysql, partnerId string) (bool, error) {
	d, err := models.FindDataPlatformBusinessPartnerGeneralDatum(
		ctx, db, partnerId, models.DataPlatformBusinessPartnerGeneralDatumColumns.BusinessPartner,
	)
	if err != nil {
		return false, xerrors.Errorf("cannot get data: %w", err)
	}
	if d == nil {
		return false, nil
	}
	if d.BusinessPartner != partnerId {
		return false, nil
	}
	return true, nil
}
```

## Input
data-platform-api-business-partner-exconf-rmq-kube では、以下のInputファイルをRabbitMQからJSON形式で受け取ります。  

```
{
	"connection_key": "response",
	"result": true,
	"redis_key": "abcdefg",
	"runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
	"business_partner": null,
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