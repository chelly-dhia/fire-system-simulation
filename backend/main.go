package main

import (
	"fmt"
	"math/rand"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var isFire = false

func publish(client MQTT.Client, topic string, value float64) {
	token := client.Publish(topic, 0, false, fmt.Sprintf("%f", value))
	token.Wait()
}

func simulateSensors(client MQTT.Client) {
	for {
		var temperature, smoke, co2 float64

		if isFire {
			// Simulate higher values in fire case
			temperature = 35.0 + rand.Float64()*(45.0-35.0)
			smoke = 60.0 + rand.Float64()*(100.0-60.0)
			co2 = 600.0 + rand.Float64()*(1000.0-600.0)
		} else {
			// Simulate lower values in unfire case
			temperature = 15.0 + rand.Float64()*(30.0-15.0)
			smoke = rand.Float64() * 50.0
			co2 = rand.Float64() * 500.0
		}

		// Publish sensor data to the MQTT broker
		publish(client, "sensors/temperature", temperature)
		publish(client, "sensors/smoke", smoke)
		publish(client, "sensors/co2", co2)

		fmt.Printf("Published: Temperature=%.2f, Smoke=%.2f, CO2=%.2f (Fire Status: %t)\n", temperature, smoke, co2, isFire)

		// Send data every 5 seconds
		time.Sleep(5 * time.Second)
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
	opts := MQTT.NewClientOptions().AddBroker(broker).SetClientID("sensor-simulator")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Handle fire/unfire status
	handleFireStatus(client)

	// Simulate sensor data
	simulateSensors(client)
}
