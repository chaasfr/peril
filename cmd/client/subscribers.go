package main

import (
	"fmt"
	"log"

	"github.com/chaasfr/peril/internal/gamelogic"
	"github.com/chaasfr/peril/internal/pubsub"
	"github.com/chaasfr/peril/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribes(connection *amqp.Connection, gs *gamelogic.GameState, amqpChan *amqp.Channel) {
	err := pubsub.SubscribeJson(
		connection,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", gs.Player.Username),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = pubsub.SubscribeJson(
		connection,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, gs.Player.Username),
		fmt.Sprintf("%s.*", routing.ArmyMovesPrefix),
		pubsub.Transient,
		handleMove(gs, amqpChan),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = pubsub.SubscribeJson(
		connection,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		fmt.Sprintf("%s.*", routing.WarRecognitionsPrefix),
		pubsub.Durable,
		handleWar(gs, amqpChan),
	)
	if err != nil {
		log.Fatal(err)
	}
}
