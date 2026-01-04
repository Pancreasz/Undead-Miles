// components/Navbar.tsx
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Car, Eye, Home } from 'lucide-react';

const Navbar = () => {
  const location = useLocation();

  // Helper to check if link is active
  const isActive = (path: string) => location.pathname === path;

  // Common styling for nav items
  const navClass = (path: string) =>
    `flex items-center gap-2 px-4 py-2 rounded-lg transition-colors duration-200 ${
      isActive(path)
        ? 'bg-blue-600 text-white shadow-md'
        : 'text-gray-600 hover:bg-gray-100'
    }`;

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-5xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          {/* Logo / Brand */}
          <div className="text-xl font-bold text-blue-600">
            TripShare
          </div>

          {/* Navigation Links */}
          <div className="flex space-x-2">
            <Link to="/" className={navClass('/')}>
              <Home size={18} />
              <span className="hidden sm:inline">Home</span>
            </Link>
            
            <Link to="/create-trip" className={navClass('/create-trip')}>
              <Car size={18} />
              <span className="hidden sm:inline">Driver Mode</span>
            </Link>

            <Link to="/create-watcher" className={navClass('/create-watcher')}>
              <Eye size={18} />
              <span className="hidden sm:inline">Find Ride</span>
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;