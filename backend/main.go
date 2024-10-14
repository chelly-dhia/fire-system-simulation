package main

import (
	"fmt"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var temperatureThreshold = 30.0
var smokeThreshold = 50.0
var co2Threshold = 500.0

func messageHandler(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())
	value, _ := strconv.ParseFloat(payload, 64)

	if topic == "sensors/temperature" && value > temperatureThreshold {
		fmt.Printf("ALERT! High temperature detected: %.2f\n", value)
	} else if topic == "sensors/smoke" && value > smokeThreshold {
		fmt.Printf("ALERT! High smoke level detected: %.2f\n", value)
	} else if topic == "sensors/co2" && value > co2Threshold {
		fmt.Printf("ALERT! High CO2 level detected: %.2f\n", value)
	}
}

func main() {
	broker := "tcp://mqtt-broker:1883"
	opts := MQTT.NewClientOptions().AddBroker(broker).SetClientID("fire-alarm-backend")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client.Subscribe("sensors/temperature", 0, messageHandler)
	client.Subscribe("sensors/smoke", 0, messageHandler)
	client.Subscribe("sensors/co2", 0, messageHandler)

	fmt.Println("Subscribed to sensor topics, waiting for data...")
	select {} // Keep the service running
}
