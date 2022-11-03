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

## RabbitMQ からの JSON Input
data-platform-api-business-partner-exconf-rmq-kube は、入力ファイルとして、RabbitMQ からのメッセージを JSON 形式で受け取ります。Input の サンプルJSON は、Inputs フォルダ内にあります。

## RabbitMQ からのメッセージ受信による イベントドリヴン の ランタイム実行
data-platform-api-business-partner-exconf-rmq-kube は、RabbitMQ からのメッセージを受け取ると、イベントドリヴンでランタイムを実行します。  
AION の仕様では、Kubernetes 上 の 当該マイクロサービスPod は 立ち上がったまま待機状態で当該メッセージを受け取り、（コンテナ起動などの段取時間をカットして）即座にランタイムを実行します。　 

## RabbitMQ の マスタサーバ環境
data-platform-api-business-partner-exconf-rmq-kube が利用する RabbitMQ のマスタサーバ環境は、rabbitmq-on-kubernetes です。  

## RabbitMQ の Golang Runtime ライブラリ
data-platform-api-business-partner-exconf-rmq-kube は、RabbitMQ の Golang Runtime ライブラリ として、rabbitmq-golang-clientを利用しています。

## デプロイ・稼働
data-platform-api-business-partner-exconf-rmq-kube の デプロイ・稼働 を行うためには、aion-service-definitions の services.yml に、本レポジトリの services.yml を設定する必要があります。

kubectl apply - f 等で Deployment作成後、以下のコマンドで Pod が正しく生成されていることを確認してください。

```
$ kubectl get pods
```


## Output
data-platform-api-business-partner-exconf-rmq-kube では、[golang-logging-library](https://github.com/latonaio/golang-logging-library) により、Output として、RabbitMQ へのメッセージを JSON 形式で出力します。ビジネスパートナの対象値が存在する場合 true、存在しない場合 false、を返します。"cursor" ～ "time"は、golang-logging-library による 定型フォーマットの出力結果です。

```
{
    "cursor": "/go/src/github.com/latonaio/data-platform-business-partner-exconf/main.go#L66",
    "function": "('DPFM_API_ORDERS_SRV', 'creates', 'A_HeaderPartner', 'Customer', 'DPFM_API_BUSINESS_PARTNER_SRV', 'exconf');",
    "level": "INFO",
	"service_label": "ORDERS",
    "Customer": {
        "BusinessPartner": 101,
        "ExistenceConf": true
    },
    "runtime_session_id": "boi9ar543dg91ipdnspi099u231280ab0v8af0ew",
    "time": "2022-10-09T20:36:39+09:00"
}
```