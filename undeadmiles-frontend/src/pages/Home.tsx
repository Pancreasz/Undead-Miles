import { useEffect, useState, useRef } from 'react';
import axios from 'axios';
import "../App.css";

interface Trip {
  ID: string;
  Origin: string;
  Destination: string;
  PriceThb: number;
}

function App() {
  const [trips, setTrips] = useState<Trip[]>([]);
  const [error, setError] = useState<string>('');
  
  // --- NEW STATE FOR NOTIFICATIONS ---
  const [notifications, setNotifications] = useState<string[]>([]);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const prevTripsRef = useRef<Trip[]>([]);

  const fetchTrips = async (isPolling: boolean = false) => {
    try {
      const response = await axios.get('http://localhost:8080/trips');
      const newData: Trip[] = response.data;

      // WATCHER LOGIC
      if (isPolling && newData.length > prevTripsRef.current.length) {
        const diff = newData.length - prevTripsRef.current.length;
        const newTrip = newData[newData.length - 1];
        
        // Create a new message
        const timestamp = new Date().toLocaleTimeString();
        const msg = `[${timestamp}] New Trip: ${newTrip.Origin} ➔ ${newTrip.Destination} (฿${newTrip.PriceThb})`;
        
        // Add to the TOP of the notification list
        setNotifications(prev => [msg, ...prev]);
      }

      setTrips(newData);
      prevTripsRef.current = newData;
      setError('');
    } catch (err) {
      console.error("Error fetching trips:", err);
      if (!isPolling) setError('Failed to load trips.');
    }
  };

  useEffect(() => {
    fetchTrips(false);
    const interval = setInterval(() => fetchTrips(true), 3000);
    return () => clearInterval(interval);
  }, []);

  // Toggle Dropdown
  const toggleDropdown = () => setIsDropdownOpen(!isDropdownOpen);

  // Clear Notifications
  const clearNotifications = () => {
    setNotifications([]);
    setIsDropdownOpen(false);
  };

  return (
    <div className="container">
      {/* --- HEADER WITH NOTIFICATION BELL --- */}
      <header className="app-header">
        <h1 id="page-title">Undead Miles Marketplace</h1>
        
        <div className="notification-wrapper">
          <button className="bell-btn" onClick={toggleDropdown}>
            🔔
            {notifications.length > 0 && (
              <span className="badge">{notifications.length}</span>
            )}
          </button>

          {/* DROPDOWN MENU */}
          {isDropdownOpen && (
            <div className="dropdown-menu">
              <div className="dropdown-header">
                <span>Notifications</span>
                <button className="clear-btn" onClick={clearNotifications}>Clear</button>
              </div>
              <div className="dropdown-content">
                {notifications.length === 0 ? (
                  <p className="empty-msg">No new notifications</p>
                ) : (
                  notifications.map((note, index) => (
                    <div key={index} className="notification-item">
                      {note}
                    </div>
                  ))
                )}
              </div>
            </div>
          )}
        </div>
      </header>

      {error && <p className="error">{error}</p>}

      <div className="trip-list">
        {trips.length === 0 ? (
          <p>No trips available...</p>
        ) : (
          trips.map(trip => (
            <div key={trip.ID} className="trip-card">
              <h3>{trip.Origin} ➔ {trip.Destination}</h3>
              <p>Price: ฿{trip.PriceThb}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

export default App;