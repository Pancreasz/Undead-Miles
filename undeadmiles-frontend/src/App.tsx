import { useEffect, useState } from 'react';
import axios from 'axios';
import './App.css';

// 1. Define what a Trip looks like
interface Trip {
  id: string;
  origin: string;
  destination: string;
  price_thb: number;
}

function App() {
  const [trips, setTrips] = useState<Trip[]>([]);
  const [error, setError] = useState<string>('');

  // 2. Fetch data from K8s when page loads
  useEffect(() => {
    // Note: We access localhost:8080 directly because K8s exposed it via NodePort/LoadBalancer
    axios.get('http://localhost:8080/trips')
      .then(response => {
        console.log("Data loaded:", response.data);
        setTrips(response.data);
      })
      .catch(err => {
        console.error("Error fetching trips:", err);
        setError('Failed to load trips. Is K8s running?');
      });
  }, []);

  return (
    <div className="container">
      <h1 id="page-title">Undead Miles Marketplace</h1>
      
      {error && <p className="error">{error}</p>}

      <div className="trip-list" id="trip-list">
        {trips.length === 0 ? (
          <p>No trips available...</p>
        ) : (
          trips.map(trip => (
            <div key={trip.id} className="trip-card" data-testid="trip-card">
              <h3>{trip.origin} ➔ {trip.destination}</h3>
              <p>Price: ฿{trip.price_thb}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

export default App;