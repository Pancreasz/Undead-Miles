// src/pages/CreateWatcher.tsx
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const CreateWatcher = () => {
  const [formData, setFormData] = useState({
    origin: '',
    destination: ''
  });
  
  // This simulates the logged-in user ID you provided
  const USER_ID = "passenger_007"; 

  const [isWatching, setIsWatching] = useState(false);
  const [notifications, setNotifications] = useState<any[]>([]);

  // 1. Handle Create Watcher
  const handleStartWatching = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = {
        user_id: USER_ID,
        origin: formData.origin,
        destination: formData.destination
      };
      
      await axios.post('http://localhost:8081/watchers', payload);
      setIsWatching(true); // Enable polling
      alert("Watcher created! Scanning for rides...");
    } catch (error) {
      console.error("Error creating watcher:", error);
      alert("Failed to start watcher");
    }
  };

  // 2. Poll for Notifications (Runs only when isWatching is true)
  useEffect(() => {
    let interval: number;

    if (isWatching) {
      const fetchNotifications = async () => {
        try {
          const res = await axios.get(`http://localhost/notifications/${USER_ID}`);
          // Assuming the API returns a list of notifications
          if (res.data) {
             console.log("New notification found:", res.data);
             // Append or set notifications (depends on your API structure)
             // If API returns a single object, wrap it in array [res.data]
             setNotifications(prev => Array.isArray(res.data) ? res.data : [res.data]); 
          }
        } catch (error) {
          console.error("No notifications yet or error fetching:", error);
        }
      };

      // Poll every 5 seconds
      interval = setInterval(fetchNotifications, 5000);
      
      // Run immediately once
      fetchNotifications();
    }

    return () => clearInterval(interval); // Cleanup on unmount
  }, [isWatching]);

  return (
    <div className="max-w-2xl mx-auto mt-10 p-6 bg-white rounded-xl shadow-sm border border-gray-100">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 flex items-center gap-2">
        👀 Create Watcher (Passenger)
      </h2>
      <p className="text-gray-500 mb-6">
        Tell us your route, and we'll notify you when a driver is available.
      </p>

      {!isWatching ? (
        <form onSubmit={handleStartWatching} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Where are you starting from?</label>
            <input 
              type="text" 
              required
              placeholder="Start location..." 
              className="w-full p-2 border rounded-lg"
              value={formData.origin}
              onChange={(e) => setFormData({...formData, origin: e.target.value})}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Where are you going?</label>
            <input 
              type="text" 
              required
              placeholder="Destination..." 
              className="w-full p-2 border rounded-lg"
              value={formData.destination}
              onChange={(e) => setFormData({...formData, destination: e.target.value})}
            />
          </div>
          
          <button type="submit" className="w-full bg-emerald-600 text-white py-3 rounded-lg font-semibold hover:bg-emerald-700 transition">
            Start Watching
          </button>
        </form>
      ) : (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 text-blue-800 rounded-lg text-center animate-pulse">
            🔍 Scanning for drivers... (ID: {USER_ID})
          </div>

          {/* Display Notifications */}
          <div className="mt-6">
            <h3 className="font-bold text-lg mb-2">Notifications</h3>
            {notifications.length === 0 && <p className="text-gray-400">No matches found yet...</p>}
            
            {notifications.map((notif, index) => (
              <div key={index} className="p-3 mb-2 border border-green-200 bg-green-50 rounded-lg">
                {/* Render notification content based on your API response structure */}
                <p><strong>Match Found!</strong></p>
                <pre className="text-xs mt-2">{JSON.stringify(notif, null, 2)}</pre>
              </div>
            ))}
          </div>
          
          <button 
            onClick={() => setIsWatching(false)}
            className="text-sm text-red-500 underline mt-4"
          >
            Stop Watching
          </button>
        </div>
      )}
    </div>
  );
};

export default CreateWatcher;