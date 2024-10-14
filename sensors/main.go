package main

import (
	"fmt"
	"math/rand"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func publish(client MQTT.Client, topic string, value float64) {
	token := client.Publish(topic, 0, false, fmt.Sprintf("%f", value))
	token.Wait()
}

func simulateSensors() {
	broker := "tcp://mqtt-broker:1883" // Link to the broker via Docker Compose service name
	opts := MQTT.NewClientOptions().AddBroker(broker).SetClientID("sensor-simulator")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		// Simulating random sensor data
		temperature := 15.0 + rand.Float64()*(35.0-15.0)
		smoke := rand.Float64() * 100.0
		co2 := rand.Float64() * 1000.0

		// Publish sensor data to the MQTT broker
		publish(client, "sensors/temperature", temperature)
		publish(client, "sensors/smoke", smoke)
		publish(client, "sensors/co2", co2)

		fmt.Printf("Published: Temperature=%.2f, Smoke=%.2f, CO2=%.2f\n", temperature, smoke, co2)
		time.Sleep(5 * time.Second) // Send data every 5 seconds
	}
}

func main() {
	simulateSensors()
}
