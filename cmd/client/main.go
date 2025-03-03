package main

import (
	"fmt"
	"log"

	"github.com/chaasfr/peril/internal/gamelogic"
	"github.com/chaasfr/peril/internal/pubsub"
	"github.com/chaasfr/peril/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)

	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	log.Println("connection successfull to amqp")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", username),
		routing.PauseKey,
		0,
	)
	if err != nil {
		log.Fatal(err)
	}

	gamelogic.PrintClientHelp()
	gameState := gamelogic.NewGameState(username)

	amqpChan, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	subscribes(connection, gameState, amqpChan)
	repl(gameState, amqpChan)

	log.Println("shutting down")
	err = connection.Close()
	if err != nil {
		log.Fatal(err)
	}
}
