package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

	 _,_, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s",username),
		routing.PauseKey,
		0,
	 )
	 if err != nil {
		log.Fatal(err)
	 }

	gamelogic.PrintClientHelp()
	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJson(
		connection,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", username),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatal(err)
	}

	repl(gameState)
	
	log.Println("shutting down")
	connection.Close()
}


func repl(gameState *gamelogic.GameState) {
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
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				log.Println(err)
			}
		case "move":
			_, err := gameState.CommandMove(words)
			if err == nil {
				log.Println("move successful")
			} else {
				log.Println(err)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			log.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			stop = true
		default:
			log.Println("unknown command.")
		}
	}
}