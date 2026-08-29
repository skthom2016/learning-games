import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import './AdminLoginScreen.css';

/**
 * AdminLoginScreen - Simple password protection for admin panel
 * Prevents kids from accidentally accessing admin features
 */
function AdminLoginScreen() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    // Simulate a small delay for better UX
    await new Promise(resolve => setTimeout(resolve, 300));

    const result = login(password);

    if (result.success) {
      // Redirect to admin dashboard
      navigate('/admin');
    } else {
      setError('Incorrect password. Please try again.');
      setPassword('');
      setIsLoading(false);
    }
  };

  return (
    <div className="admin-login-screen">
      <div className="login-container">
        <div className="login-card">
          <div className="login-header">
            <h1>🔒 Admin Login</h1>
            <p>Enter password to access admin panel</p>
          </div>

          <form onSubmit={handleSubmit} className="login-form">
            <div className="form-group">
              <label htmlFor="password">Admin Password:</label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter admin password"
                autoFocus
                disabled={isLoading}
              />
            </div>

            {error && <div className="error-message">{error}</div>}

            <div className="login-actions">
              <button
                type="button"
                className="cancel-button"
                onClick={() => navigate('/')}
                disabled={isLoading}
              >
                ← Back to Games
              </button>
              <button
                type="submit"
                className="submit-button"
                disabled={!password.trim() || isLoading}
              >
                {isLoading ? 'Logging in...' : 'Login'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

export default AdminLoginScreen;
