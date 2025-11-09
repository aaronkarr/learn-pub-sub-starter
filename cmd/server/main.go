package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	pubsub "github.com/aaronkarr/learn-pub-sub-starter/internal"
	"github.com/aaronkarr/learn-pub-sub-starter/internal/gamelogic"
	"github.com/aaronkarr/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	gamelogic.PrintServerHelp()
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

	paused := routing.PlayingState{
		IsPaused: true,
	}
	err = pubsub.PublishJSON(rmqChannel, string(routing.ExchangePerilDirect), string(routing.PauseKey), paused)
	if err != nil {
		log.Fatalf("Could not pause: %v", err)
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
