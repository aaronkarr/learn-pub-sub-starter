package main

import (
	"fmt"
	"log"

	pubsub "github.com/aaronkarr/learn-pub-sub-starter/internal"
	"github.com/aaronkarr/learn-pub-sub-starter/internal/gamelogic"
	"github.com/aaronkarr/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	rmqChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not create MQ channel: %v", err)
	}

	gamelogic.PrintServerHelp()

repl:
	for {
		command := gamelogic.GetInput()
		if len(command) == 0 {
			continue
		}
		switch command[0] {
		case "pause":
			fmt.Println("Pausing game...")
			paused := routing.PlayingState{
				IsPaused: true,
			}
			err = pubsub.PublishJSON(
				rmqChannel,
				string(routing.ExchangePerilDirect),
				string(routing.PauseKey),
				paused,
			)
			if err != nil {
				log.Printf("Could not pause: %v", err)
			}
		case "resume":
			fmt.Println("Resuming game...")
			paused := routing.PlayingState{
				IsPaused: false,
			}
			err = pubsub.PublishJSON(
				rmqChannel,
				string(routing.ExchangePerilDirect),
				string(routing.PauseKey),
				paused,
			)
			if err != nil {
				log.Printf("Could not resume: %v", err)
			}
		case "quit":
			fmt.Println("Goodbye.")
			break repl
		default:
			fmt.Println("Unknown command.")
			gamelogic.PrintServerHelp()
			continue
		}
	}
}
