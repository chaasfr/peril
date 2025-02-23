# Peril
playing around with rabbitMQ


### ToDo list
- dockerise server and clients
- CD
- proper readme
- unit-test repl
- automated created of exchanges and gamelog queue

# How to install
- install go 
- install rabbitMQ CLI `sudo apt install rabbitmq-server`

## How to start
- run `rabbit.sh start`
- in the UI of rabbitMQ (by default http://localhost:15672/) create the following:
  - one direct exchange "peril_direct"
  - one topic exchange "peril_topic"
  - one fanout exchange "peril_dlx"
  - one durable queue "peril_dlq" and bind it to the exchange "peril_dlx"
- run the server: `go run ./cmd/server/`
- run clients (one per player): `go run ./cmd/client/`