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

func simulateCO2Sensor(client MQTT.Client) {
	for {
		var co2 float64
		if isFire {
			co2 = 600.0 + rand.Float64()*(1000.0-600.0)
		} else {
			co2 = rand.Float64() * 500.0
		}

		publish(client, "sensors/co2", co2)
		fmt.Printf("Published: CO2=%.2f (Fire Status: %t)\n", co2, isFire)
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
	opts := MQTT.NewClientOptions().AddBroker(broker).SetClientID("co2-sensor")
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	handleFireStatus(client)
	simulateCO2Sensor(client)
}
