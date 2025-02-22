# Peril

playing around with rabbitMQ

# How to setup
- install go & rabbitMQ
- run `rabbit.sh start`
- in the UI of rabbitMQ (by default http://localhost:15672/) create the following:
  - one direct exchange "peril_direct"
  - one topic exchange "peril_topic"
  - one fanout exchange "peril_dlx"
  - one durable queue "peril_dlq" and bind it to the exchange "peril_dlx"