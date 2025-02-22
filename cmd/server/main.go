package main

import (
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

	_,_, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug + ".*",
		1,
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
		switch(words[0]) {
		case "pause":
			log.Println("sending a pause message")
			pubsub.PublishJSON(
				amqpChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
		case "resume":
			log.Println("sending a resume message")
			pubsub.PublishJSON(
				amqpChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
		case "quit":
			log.Println("exiting")
			stop = true
		default:
			log.Println("I donùt understand this command")
		}
	}


	// // wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan 
	// log.Println("shutting down")
	// connection.Close()
	
}
