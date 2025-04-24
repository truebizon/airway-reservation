## ドローン航路予約システム

## 概要

ドローン航路予約システムは複数のAPIで構成されます。
- 航路予約API
- 航路予約一覧取得API（運航事業者向け）
- 航路予約一覧取得API（航路運営者向け）
- 航路予約詳細取得API
- 航路予約取消API（運航事業者向け）
- 航路予約撤回API（航路運営者向け）


### ディレクトリ構成

```
$ tree -L 2 -I ".aws-sam|.DS_Store|.git|tmp|.env" -a
.
├── .dockerignore                  # Dockerfileでコピーしないファイルを定義
├── .env.local                     # 環境変数ファイル（Local確認用。コンテナとしてデプロイする際は、.envにリネーム）
├── .gitignore
├── Makefile
├── README.md
├── bin
│   └── gen-proto.sh
├── cmd
│   ├── app                        # アプリ起動時に実行される
│   ├── cleanSchema
│   └── migration                  # DBマイグレーションの設定と実行
├── containers
│   └── airway_reservation
├── database
│   ├── dbdoc
│   ├── migration
│   └── tbls.yml
├── docker-compose.local.yml       # ローカルで起動する際のコンテナ用
├── docker-compose.yml             # 外部に連携する際のコンテナ用
├── go.mod
├── go.sum
├── internal
│   ├── app
│   └── pkg
├── proto
│   ├── pkg
│   └── third_party
├── samconfig.toml
└── template.yaml
```

### テーブル構成

[README.md](database/dbdoc/README.md)

## 開発環境構築手順

### ソフトウェアバージョン

- git
- docker
- go(v1.21.6)

### リポジトリクローン

```
git clone git@github.com:ODS-IS-UASL/airway-reservation.git
```

### 環境変数ファイル設定（コンテナデプロイ用）

```
cp -pr .env.local .env
```

.env
USE_MQTT=true　// ./.broker.crt配置が必須

※ローカル環境での動作確認ではUSE_MQTT=falseとしてMQTTブローカーへのPublish処理をスキップ可能。

### go モジュールインストール

go.mod ファイルに書いてあるモジュールがローカルにインストールされる

```
go mod tidy -v
```

### RDS のマイグレーション

[README.md](cmd/migration/README.md)

ローカル環境のテーブルやカラムを作成更新する

```
make migrate
```

### broker.crt ファイルの配置

API 実行時、主に create/update の際に mqtt 通信で publish を実施するため、ブローカーに対する証明書を配置しておく。<br>
なお\*.cert の拡張子は git 管理しない。<br>
MQTTブローカーへのイベントパブリッシュを行わない場合は、.envのUSE_MQTT=falseとすること。

### ローカルで動作確認する方法

#### docker を起動

docker の起動と api 動作確認ができる

```
# ローカル用
$ make docker-local-up

# 外部連携用
$ make docker-up
```

以下のようなログが出力されるまで待つ

```
データベースとの接続に成功しました

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.11.4
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8081
```

#### api を実行する

postman や curl を使って api を実行できる  
動作確認したいサービスの docker を起動してから実行する  
以下は curl を使った 航路予約 api の一例

```
# 航路予約API
$ curl -i -X POST 'http://localhost:8088/v1/airwayReservations' -H "Content-Type: application/json" -d '{"operatorId": "60c895e5-321a-fe8a-af39-f005f3206efb","airwaySections": [{"airwaySectionId": "123e4567-e89b-12d3-a456-426614174000","startAt": "2025-02-01T23:59:59Z","endAt": "2025-02-02T00:00:00Z"},{"airwaySectionId": "123e4567-e89b-12d3-a456-426614174001","startAt": "2025-02-02T00:01:00Z","endAt": "2025-02-02T00:02:00Z"}]}'

# 航路予約一覧取得API（運航事業者向け）
$ curl -i -X GET 'http://localhost:8088/v1/operator/60c895e5-321a-fe8a-af39-f005f3206efb/airwayReservations'
$ curl -i -X GET 'http://localhost:8088/v1/operator/60c895e5-321a-fe8a-af39-f005f3206efb/airwayReservations?page=1'

# 航路予約一覧取得API（航路運営者向け）
$ curl -i -X GET 'http://localhost:8088/v1/admin/airwayReservations'
$ curl -i -X GET 'http://localhost:8088/v1/admin/airwayReservations?page=1'

# 航路予約詳細取得API
$ curl -i -X GET 'http://localhost:8088/v1/airwayReservations/5a50b6b3-f780-40d4-8c9b-7d6d369640ad'

# 航路予約取消API（運航事業者向け）
$ curl -i -X PUT 'http://localhost:8088/v1/airwayReservations/5a50b6b3-f780-40d4-8c9b-7d6d369640ad/cancel'
# 航路予約取消API（運航事業者向け）
$ curl -i -X PUT 'http://localhost:8088/v1/admin/airwayReservations/5a50b6b3-f780-40d4-8c9b-7d6d369640ad/rescind'

```

## 新規テーブルを作る場合

database ディレクトリにサーバー名のディレクトリを追加
必要な sql ファイルを作成したら Makefile の migrate エリアに以下のように名前を追加

```Makefile
migrate:
    make airway_reservation
```

その後作成・更新を再度を行う

```
make migrate
```

問題なければ最後に doc を更新する

```
tbls doc -c database/tbls.yml --rm-dist
```

## 外部にコンテナを連携する手順

外部向けの docker イメージの作成と立ち上げ

```
make docker-up
```

tar ファイルの作成

```
docker save -o ar.tar airway_reservation-app postgis/postgis
```

tar ファイルと env ファイルと compose ファイルを圧縮

```
tar -czvf containers.tar.gz ar.tar docker-compose.yml .env.local
```

## 外部向けコンテナイメージを実行する方法

圧縮ファイルを適当なディレクトリに移して展開

```
tar -xzvf containers.tar.gz
```

tar ファイルのロード

```
docker load < ar.tar
```

コンテナの立ち上げ

```
docker compose -f docker-compose.yml up
```

API の実行

```
curl -i -X GET http://localhost:8088/v1/airwayReservations/5a50b6b3-f780-40d4-8c9b-7d6d369640ad
```

## 問合せ及び要望に関して

- 本リポジトリは現状は主に配布目的の運用となるため、IssueやPull Requestに関しては受け付けておりません。

## ライセンス

- 本リポジトリはMITライセンスで提供されています。
- 本リポジトリ内のソースコードの著作権は、KDDIスマートドローン株式会社に帰属します。

## 免責事項

- 本リポジトリの内容は予告なく変更・削除する可能性があります。
- 本リポジトリの利用により生じた損失及び損害等について、いかなる責任も負わないものとします。
