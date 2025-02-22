package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)


func repl(gs *gamelogic.GameState, amqpChan *amqp.Channel) {
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
			err := gs.CommandSpawn(words)
			if err != nil {
				log.Println(err)
			}
		case "move":
			move, err := gs.CommandMove(words)
			if err == nil {
				log.Println("move successful")
				err = pubsub.PublishJSON(amqpChan, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, gs.Player.Username), move)
				log.Println("move published successfully")
			}
			if err != nil {
				log.Println(err)
			}
		case "status":
			gs.CommandStatus()
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