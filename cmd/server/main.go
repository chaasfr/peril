package main

import (
	"log"

	"github.com/chaasfr/peril/internal/gamelogic"
	"github.com/chaasfr/peril/internal/pubsub"
	"github.com/chaasfr/peril/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)

	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	log.Println("server connection successfull to amqp")
	amqpChan, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = pubsub.SubscribeGob(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.Durable,
		handleLog(),
	)
	if err != nil {
		log.Fatal(err)
	}

	gamelogic.PrintServerHelp()
	stop := false
	for {
		if stop {
			break
		}

		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			log.Println("sending a pause message")
			err = pubsub.PublishJSON(
				amqpChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				log.Fatal(err)
			}
		case "resume":
			log.Println("sending a resume message")
			err = pubsub.PublishJSON(
				amqpChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				log.Fatal(err)
			}
		case "quit":
			log.Println("exiting")
			stop = true
		default:
			log.Println("I don't understand this command")
		}
	}

}
