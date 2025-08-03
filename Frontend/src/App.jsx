import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { Toaster } from 'react-hot-toast';
import Login from './pages/Login';
import Items from './pages/Items';
import ProtectedRoute from './components/ProtectedRoute';
import { auth } from './utils/api';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(null); // null = loading, boolean = auth state

  useEffect(() => {
    // Check authentication status on app load
    const isAuth = auth.isAuthenticated();
    setIsAuthenticated(isAuth);
    
    // If not authenticated, clear any stale user data
    if (!isAuth) {
      auth.clearUserData();
    }
  }, []);

  // Show loading while checking auth
  if (isAuthenticated === null) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    );
  }

  return (
    <>
      <Router>
        <Routes>
          <Route 
            path="/login" 
            element={
              isAuthenticated ? <Navigate to="/items" replace /> : <Login setIsAuthenticated={setIsAuthenticated} />
            } 
          />
          <Route 
            path="/items" 
            element={
              <ProtectedRoute isAuthenticated={isAuthenticated} setIsAuthenticated={setIsAuthenticated}>
                <Items setIsAuthenticated={setIsAuthenticated} />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/" 
            element={<Navigate to={isAuthenticated ? "/items" : "/login"} replace />} 
          />
        </Routes>
      </Router>
      <Toaster 
        position="top-right"
        toastOptions={{
          duration: 4000,
          style: {
            background: '#363636',
            color: '#fff',
          },
          success: {
            style: {
              background: 'green',
            },
          },
          error: {
            style: {
              background: 'red',
            },
          },
        }}
      />
    </>
  );
}

export default App;
