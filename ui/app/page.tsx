"use client";
import mqtt from 'mqtt';
import React, { useState, useEffect } from 'react';
import { FaTemperatureHigh, FaSmoking, FaCloud } from 'react-icons/fa';
import { MdOutlineAlarmOn, MdOutlineAlarmOff } from 'react-icons/md';

export default function Home() {
  const client = mqtt.connect('ws://localhost:9001');
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
      else setAlarm(false);
    });
  }, [temperature, smoke, co2]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 p-4">
      <div className="max-w-3xl w-full bg-white shadow-lg rounded-lg p-6 space-y-6">
        <h1 className="text-3xl font-bold text-center text-gray-800">Fire Alarm System</h1>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="p-4 bg-blue-50 rounded-lg shadow flex flex-col items-center">
            <FaTemperatureHigh className="text-4xl text-blue-500 mb-2" />
            <h2 className="text-lg font-semibold text-gray-600">Temperature</h2>
            <p className={`text-2xl font-bold ${temperature > 30 ? 'text-red-500' : 'text-blue-500'}`}>
              {temperature}°C
            </p>
          </div>
          <div className="p-4 bg-yellow-50 rounded-lg shadow flex flex-col items-center">
            <FaSmoking className="text-4xl text-yellow-500 mb-2" />
            <h2 className="text-lg font-semibold text-gray-600">Smoke Level</h2>
            <p className={`text-2xl font-bold ${smoke > 50 ? 'text-red-500' : 'text-yellow-500'}`}>
              {smoke} ppm
            </p>
          </div>
          <div className="p-4 bg-green-50 rounded-lg shadow flex flex-col items-center">
            <FaCloud className="text-4xl text-green-500 mb-2" />
            <h2 className="text-lg font-semibold text-gray-600">CO2 Level</h2>
            <p className={`text-2xl font-bold ${co2 > 500 ? 'text-red-500' : 'text-green-500'}`}>
              {co2} ppm
            </p>
          </div>
        </div>
        <div className={`p-4 rounded-lg ${alarm ? 'bg-red-100' : 'bg-green-100'} shadow flex items-center justify-center`}>
          {alarm ? <MdOutlineAlarmOn className="text-4xl text-red-600 mr-2" /> : <MdOutlineAlarmOff className="text-4xl text-green-600 mr-2" />}
          <div>
            <h2 className="text-lg font-semibold text-gray-600">Alarm Status</h2>
            <p className={`text-2xl font-bold ${alarm ? 'text-red-600' : 'text-green-600'}`}>
              {alarm ? 'Active' : 'Inactive'}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
