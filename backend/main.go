package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	alarmActive   = false                  // Tracks alarm status
	mqttBrokerURL = "tcp://localhost:1883" // MQTT broker address
	clientID      = "alarm-system"         // Unique MQTT client ID
)

// Publish message to MQTT topic
func publish(client MQTT.Client, topic string, message string) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
}

// Monitor sensor data and trigger alarm if conditions are met
func monitorSensors(client MQTT.Client) {
	client.Subscribe("sensors/+", 0, func(client MQTT.Client, msg MQTT.Message) {
		value, _ := strconv.ParseFloat(string(msg.Payload()), 64)
		topic := msg.Topic()

		fmt.Printf("Received %s: %.2f\n", topic, value)

		// Check if any sensor reading exceeds threshold
		if (topic == "sensors/temperature" && value > 30) ||
			(topic == "sensors/smoke" && value > 50) ||
			(topic == "sensors/co2" && value > 500) {
			if !alarmActive {
				activateAlarm(client)
			}
		} else if alarmActive {
			deactivateAlarm(client)
		}
	})
}

// Activate the alarm and notify UI
func activateAlarm(client MQTT.Client) {
	alarmActive = true
	fmt.Println("ALARM ACTIVATED!")
	publish(client, "alarm/status", "1") // Notify UI: alarm active
}

// Deactivate the alarm and notify UI
func deactivateAlarm(client MQTT.Client) {
	alarmActive = false
	fmt.Println("ALARM DEACTIVATED!")
	publish(client, "alarm/status", "0") // Notify UI: alarm inactive
}

func main() {
	// Setup MQTT client
	opts := MQTT.NewClientOptions().AddBroker(mqttBrokerURL).SetClientID(clientID)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("MQTT Connection Error:", token.Error())
		return
	}
	fmt.Println("Connected to MQTT Broker")

	// Monitor sensors for alarm conditions
	monitorSensors(client)

	// Handle termination gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down...")
	client.Disconnect(250)
}
