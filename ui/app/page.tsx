"use client"
import mqtt from 'mqtt';
import React, { useState, useEffect } from 'react';
import { FaTemperatureHigh, FaSmoking, FaCloud } from 'react-icons/fa';
import { MdOutlineAlarmOn, MdOutlineAlarmOff } from 'react-icons/md';

export default function Home() {
  const [temperature, setTemperature] = useState(0);
  const [smoke, setSmoke] = useState(0);
  const [co2, setCO2] = useState(0);
  const [alarm, setAlarm] = useState(false);

  // Connect to MQTT broker using WebSocket
  const client = mqtt.connect('ws://localhost:9001'); 

  useEffect(() => {
    const topics = ['sensors/temperature', 'sensors/smoke', 'sensors/co2', 'alarm/status'];
  
    // Subscribe to sensor topics on MQTT connection
    client.on('connect', () => {
      console.log('Connected to MQTT broker');
      topics.forEach(topic => {
        client.subscribe(topic, (err) => {
          if (err) {
            console.error(`Failed to subscribe to ${topic}:`, err);
          } else {
            console.log(`Subscribed to ${topic}`);
          }
        });
      });
    });
  
    // Handle incoming messages and update state
    client.on('message', (topic, message) => {
      const value = parseFloat(message.toString());
  
      switch (topic) {
        case 'sensors/temperature':
          setTemperature(value);
          break;
        case 'sensors/smoke':
          setSmoke(value);
          break;
        case 'sensors/co2':
          setCO2(value);
          break;
        case 'alarm/status':
          setAlarm(message.toString() === '1');
          break;
        default:
          console.warn(`Unhandled topic: ${topic}`);
      }
    });
  
    return () => {
      if (client) {
        client.end(); // Clean up MQTT connection on unmount
      }
    };
  }, []);
  

  // Send fire or unfire status via MQTT
  const handleFireStatus = (status: string ) => {
    client.publish('sensors/fireStatus', status);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 p-4">
      <div className="max-w-3xl w-full bg-white shadow-lg rounded-lg p-6 space-y-6">
        <h1 className="text-3xl font-bold text-center text-gray-800">Fire Alarm System</h1>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {/* Temperature Sensor */}
          <div className="p-4 bg-blue-50 rounded-lg shadow flex flex-col items-center">
            <FaTemperatureHigh className="text-4xl text-blue-500 mb-2" />
            <h2 className="text-lg font-semibold text-gray-600">Temperature</h2>
            <p className={`text-2xl font-bold ${temperature > 30 ? 'text-red-500' : 'text-blue-500'}`}>
              {temperature}°C
            </p>
          </div>

          {/* Smoke Sensor */}
          <div className="p-4 bg-yellow-50 rounded-lg shadow flex flex-col items-center">
            <FaSmoking className="text-4xl text-yellow-500 mb-2" />
            <h2 className="text-lg font-semibold text-gray-600">Smoke Level</h2>
            <p className={`text-2xl font-bold ${smoke > 50 ? 'text-red-500' : 'text-yellow-500'}`}>
              {smoke} ppm
            </p>
          </div>

          {/* CO2 Sensor */}
          <div className="p-4 bg-green-50 rounded-lg shadow flex flex-col items-center">
            <FaCloud className="text-4xl text-green-500 mb-2" />
            <h2 className="text-lg font-semibold text-gray-600">CO2 Level</h2>
            <p className={`text-2xl font-bold ${co2 > 500 ? 'text-red-500' : 'text-green-500'}`}>
              {co2} ppm
            </p>
          </div>
        </div>

        {/* Alarm Status */}
        <div className={`p-4 rounded-lg ${alarm ? 'bg-red-100' : 'bg-green-100'} shadow flex items-center justify-center`}>
          {alarm ? (
            <MdOutlineAlarmOn className="text-4xl text-red-600 mr-2" />
          ) : (
            <MdOutlineAlarmOff className="text-4xl text-green-600 mr-2" />
          )}
          <div>
            <h2 className="text-lg font-semibold text-gray-600">Alarm Status</h2>
            <p className={`text-2xl font-bold ${alarm ? 'text-red-600' : 'text-green-600'}`}>
              {alarm ? 'Active' : 'Inactive'}
            </p>
          </div>
        </div>

        {/* Fire/Unfire Control Buttons */}
        <div className="flex space-x-4">
          <button
            onClick={() => handleFireStatus('fire')}
            className="px-4 py-2 bg-red-500 text-white rounded"
          >
            Fire
          </button>
          <button
            onClick={() => handleFireStatus('unfire')}
            className="px-4 py-2 bg-green-500 text-white rounded"
          >
            Unfire
          </button>
        </div>
      </div>
    </div>
  );
}
