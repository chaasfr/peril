package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

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
			replSpawn(gs, words)
		case "move":
			replMove(gs, amqpChan, words)
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			replSpam(gs, amqpChan, words)
		case "quit":
			gamelogic.PrintQuit()
			stop = true
		default:
			log.Println("unknown command.")
		}
	}
}

func replSpawn(gs *gamelogic.GameState, words []string) {
	err := gs.CommandSpawn(words)
	if err != nil {
		log.Println(err)
	}
}

func  replMove(gs *gamelogic.GameState, amqpChan *amqp.Channel, words []string) {
	move, err := gs.CommandMove(words)
			if err == nil {
				log.Println("move successful")
				err = pubsub.PublishJSON(amqpChan, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, gs.Player.Username), move)
				log.Println("move published successfully")
			}
			if err != nil {
				log.Println(err)
			}
}

func replSpam(gs *gamelogic.GameState, amqpChan *amqp.Channel, words []string) {
	if len(words) < 2 {
		log.Println("usage: spam <int>")
		return
	}
	n, err := strconv.Atoi(words[1])
	if err != nil {
		log.Println("please use an integer instead of "+ words[1])
		return
	}

	for range n {
		maliciousLog := gamelogic.GetMaliciousLog()
		pubsub.PublishGob(
			amqpChan,
			routing.ExchangePerilTopic,
			fmt.Sprintf("%s.%s", routing.GameLogSlug,gs.Player.Username),
			routing.GameLog{
				CurrentTime: time.Now(),
				Message: maliciousLog,
				Username: gs.Player.Username,
			},
		)
	}
}