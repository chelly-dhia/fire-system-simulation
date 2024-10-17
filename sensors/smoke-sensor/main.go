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

func simulateSmokeSensor(client MQTT.Client) {
	for {
		var smoke float64
		if isFire {
			smoke = 60.0 + rand.Float64()*(100.0-60.0)
		} else {
			smoke = rand.Float64() * 50.0
		}

		publish(client, "sensors/smoke", smoke)
		fmt.Printf("Published: Smoke=%.2f (Fire Status: %t)\n", smoke, isFire)
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
	opts := MQTT.NewClientOptions().AddBroker(broker).SetClientID("smoke-sensor")
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	handleFireStatus(client)
	simulateSmokeSensor(client)
}
