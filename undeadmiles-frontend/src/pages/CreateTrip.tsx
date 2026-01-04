// src/pages/CreateTrip.tsx
import React, { useState } from 'react';
import axios from 'axios';

const CreateTrip = () => {
  const [formData, setFormData] = useState({
    origin: '',
    destination: '',
    departureTime: '',
    price: ''
  });

  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setStatus('loading');

    try {
      // Convert local datetime input to ISO format (e.g., 2026-01-05T09:00:00Z)
      const isoDate = new Date(formData.departureTime).toISOString();

      const payload = {
        driver_id: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13", // Hardcoded for demo
        origin: formData.origin,
        destination: formData.destination,
        price_thb: Number(formData.price),
        departure_time: isoDate
      };

      await axios.post('http://localhost:8080/trips', payload);
      setStatus('success');
      // Optional: Clear form
      setFormData({ origin: '', destination: '', departureTime: '', price: '' });
      
    } catch (error) {
      console.error("Error creating trip:", error);
      setStatus('error');
    }
  };

  return (
    <div className="max-w-2xl mx-auto mt-10 p-6 bg-white rounded-xl shadow-sm border border-gray-100">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 flex items-center gap-2">
        🚗 Create a Trip (Driver)
      </h2>
      
      {status === 'success' && (
        <div className="mb-4 p-3 bg-green-100 text-green-700 rounded-lg">
          Trip created successfully!
        </div>
      )}

      {status === 'error' && (
        <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
          Failed to create trip. Check console/backend.
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Origin</label>
            <input 
              type="text" 
              required
              placeholder="Start location" 
              className="w-full p-2 border rounded-lg"
              value={formData.origin}
              onChange={(e) => setFormData({...formData, origin: e.target.value})}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Destination</label>
            <input 
              type="text" 
              required
              placeholder="End location" 
              className="w-full p-2 border rounded-lg"
              value={formData.destination}
              onChange={(e) => setFormData({...formData, destination: e.target.value})}
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Departure Time</label>
          <input 
            type="datetime-local" 
            required
            className="w-full p-2 border rounded-lg"
            value={formData.departureTime}
            onChange={(e) => setFormData({...formData, departureTime: e.target.value})}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Price per Seat (THB)</label>
          <input 
            type="number" 
            required
            placeholder="0.00" 
            className="w-full p-2 border rounded-lg"
            value={formData.price}
            onChange={(e) => setFormData({...formData, price: e.target.value})}
          />
        </div>

        <button 
          type="submit" 
          disabled={status === 'loading'}
          className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition disabled:bg-gray-400"
        >
          {status === 'loading' ? 'Publishing...' : 'Publish Trip'}
        </button>
      </form>
    </div>
  );
};

export default CreateTrip;