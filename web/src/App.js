import mqtt from 'mqtt';
import React, { useState, useEffect } from 'react';

const client = mqtt.connect('ws://localhost:9001');  // WebSocket to the MQTT broker

function App() {
  const [temperature, setTemperature] = useState(0);
  const [smoke, setSmoke] = useState(0);
  const [co2, setCO2] = useState(0);
  const [alarm, setAlarm] = useState(false);

  useEffect(() => {
    client.on('connect', () => {
      client.subscribe('sensors/temperature');
      client.subscribe('sensors/smoke');
      client.subscribe('sensors/co2');
    });

    client.on('message', (topic, message) => {
      const value = parseFloat(message.toString());
      if (topic === 'sensors/temperature') setTemperature(value);
      if (topic === 'sensors/smoke') setSmoke(value);
      if (topic === 'sensors/co2') setCO2(value);

      // Trigger alarm if thresholds are exceeded
      if (temperature > 30 || smoke > 50 || co2 > 500) setAlarm(true);
    });
  }, [temperature, smoke, co2]);

  return (
    <div>
      <h1>Fire Alarm System</h1>
      <p>Temperature: {temperature}</p>
      <p>Smoke Level: {smoke}</p>
      <p>CO2 Level: {co2}</p>
      <p>Alarm: {alarm ? 'Active' : 'Inactive'}</p>
    </div>
  );
}

export default App;
