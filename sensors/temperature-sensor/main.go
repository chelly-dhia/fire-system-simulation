package main

import (
	"fmt"
	"math/rand"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var isFire bool = false

func publish(client MQTT.Client, topic string, value float64) {
	token := client.Publish(topic, 0, false, fmt.Sprintf("%f", value))
	token.Wait()
}

func simulateTemperatureSensor(client MQTT.Client) {
	for {
		var temperature float64
		if isFire {
			temperature = 35.0 + rand.Float64()*(45.0-35.0)
		} else {
			temperature = 15.0 + rand.Float64()*(30.0-15.0)
		}

		publish(client, "sensors/temperature", temperature)
		fmt.Printf("Published: Temperature=%.2f (Fire Status: %t)\n", temperature, isFire)
		time.Sleep(30 * time.Second)
	}
}

func handleFireStatus(client MQTT.Client) {
	client.Subscribe("sensors/fireStatus", 0, func(client MQTT.Client, msg MQTT.Message) {
		status := string(msg.Payload())
		if status == "fire" {
			isFire = true
			fmt.Println("Fire mode activated")
		} else if status == "unfire" {
			isFire = false
			fmt.Println("Unfire mode activated")
		}
	})
}

func main() {
	broker := "tcp://mqtt-broker:1883"
	opts := MQTT.NewClientOptions().AddBroker(broker).SetClientID("temperature-sensor")
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	handleFireStatus(client)
	simulateTemperatureSensor(client)
}
