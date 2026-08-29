import React, { createContext, useState, useContext, useEffect } from 'react';

const AuthContext = createContext();

// Simple password - in production, this should be stored securely on the backend
// For now, we'll use a default password that can be changed
const DEFAULT_ADMIN_PASSWORD = 'admin123';

export function AuthProvider({ children }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check if user is already logged in (from localStorage)
    const authStatus = localStorage.getItem('admin_authenticated');
    setIsAuthenticated(authStatus === 'true');
    setIsLoading(false);
  }, []);

  const login = (password) => {
    // In production, this should call an API to validate credentials
    if (password === DEFAULT_ADMIN_PASSWORD) {
      localStorage.setItem('admin_authenticated', 'true');
      setIsAuthenticated(true);
      return { success: true };
    }
    return { success: false, error: 'Incorrect password' };
  };

  const logout = () => {
    localStorage.removeItem('admin_authenticated');
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
