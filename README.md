# Peril
![test coverage badge](https://github.com/chaasfr/peril/actions/workflows/ci.yml/badge.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/chaasfr/peril)

## Overview

Peril is an experimental project in Go demonstrating advanced RabbitMQ patterns with direct, topic, and fanout exchanges alongside durable queues.

It powers a realtime multiplayer strategy game inspired by Risk. Leveraging advanced messaging patterns (direct, topic, and fanout exchanges) alongside durable queues, it delivers dynamic, real-time gameplay while demonstrating high-concurrency and efficiency in Go.

## Features

- Integration with RabbitMQ messaging patterns.
- Built with Go for high concurrency and efficiency.
- Ready for containerization and continuous deployment.
- Easy configuration via environment variables and config files.

## Prerequisites

- Go (version 1.20 or later).
- RabbitMQ installed and running.

## Configuration

Configure connection settings, exchanges, and queues via the provided configuration file or environment variables. Refer to the project Wiki for detailed instructions on customization.

## How to start

- run `rabbit.sh start`
- in the UI of rabbitMQ (by default http://localhost:15672/) create the following:
  - one direct exchange "peril_direct"
  - one topic exchange "peril_topic"
  - one fanout exchange "peril_dlx"
  - one durable queue "peril_dlq" and bind it to the exchange "peril_dlx"
- run the server: `go run ./cmd/server/`
- run clients (one per player): `go run ./cmd/client/`

### ToDo list

- dockerise server and clients
- CD
- proper readme
- unit-test repl
- automated created of exchanges and gamelog queue