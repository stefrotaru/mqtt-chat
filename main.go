package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// implement nqtt client to subscribe to a topic

func main() {
	opts := mqtt.NewClientOptions().AddBroker("broker.hivemq.com:1883")
	opts.SetClientID("go-simple")
	// opts.SetDefaultPublishHandler()
	opts.SetCleanSession(true)
	
	// define a function for the default message handler that prints to console
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	// define a function for the default message handler that prints to console
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println("Connected")
	}

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("topicul_de_miercuri_seara/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// keep the main thread alive
	select {} 
}