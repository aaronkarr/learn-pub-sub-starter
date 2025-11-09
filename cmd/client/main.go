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
	fmt.Println("Starting Peril client...")
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("couldn't connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("login error: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("couldn't create MQ channel: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// wait for ctrl+c
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	fmt.Println("RabbitMQ connection closed.")
}
