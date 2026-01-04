// App.tsx
import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

// Import your components
import Navbar from './components/Navbar';
import CreateTrip from './pages/CreateTrip';
import CreateWatcher from './pages/CreateWatcher';

// Placeholder for your existing Home page code
import Home from './pages/Home'; // Assuming you moved your previous code here

const App = () => {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        {/* Navigation is always visible */}
        <Navbar />

        {/* The content area changes based on the URL */}
        <div className="p-4">
          <Routes>
            {/* 1. Home Page */}
            <Route path="/" element={<Home />} />
            
            {/* 2. Create Trip Page */}
            <Route path="/create-trip" element={<CreateTrip />} />
            
            {/* 3. Create Watcher Page */}
            <Route path="/create-watcher" element={<CreateWatcher />} />
          </Routes>
        </div>
      </div>
    </BrowserRouter>
  );
};

export default App;