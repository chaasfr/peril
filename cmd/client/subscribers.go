package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)


func subscribes(connection *amqp.Connection, gs *gamelogic.GameState) {
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
		fmt.Sprintf("%s.*",routing.ArmyMovesPrefix),
		pubsub.Transient,
		handleMove(gs),
	 )
	 if err != nil {
		log.Fatal(err)
	}
}