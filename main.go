package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	var displayName string
	if len(os.Args) > 1 {
		displayName = os.Args[1]
	} else {
		displayName = "go-simple"
	}

	opts := mqtt.NewClientOptions().AddBroker("0.0.0.0:1883")
	opts.SetClientID(displayName)
	opts.SetCleanSession(true)
	
	// define a function for the default message handler that prints to console
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		// don't print own messages
		if msg.Topic() == "topicul_de_miercuri_seara/" + displayName {
			return
		}
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	// define a function for the default message handler that prints to console
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println(displayName + " Connected")
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

	// listen for keyboard input and send message to topic
	for {
		var text string
		fmt.Scanln(&text)
		token := c.Publish("topicul_de_miercuri_seara/" + displayName, 0, false, text) 
		token.Wait()
	}
}
