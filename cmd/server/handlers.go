package main

import (
	"fmt"
	"log"

	"github.com/chaasfr/peril/internal/gamelogic"
	"github.com/chaasfr/peril/internal/pubsub"
	"github.com/chaasfr/peril/internal/routing"
)

func handleLog() func(routing.GameLog) pubsub.AckType {
	return func(gl routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(gl)
		if err != nil {
			log.Fatal(err)
		}
		return pubsub.Ack
	}
}
