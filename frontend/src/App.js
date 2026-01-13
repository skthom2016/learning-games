import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';

// Child UI Screens
import PlayerSelectionScreen from './screens/child/PlayerSelectionScreen';
import GameSelectionScreen from './screens/child/GameSelectionScreen';
import GamePlayScreen from './screens/child/GamePlayScreen';
import FeedbackScreen from './screens/child/FeedbackScreen';
import ProgressOverviewScreen from './screens/child/ProgressOverviewScreen';

// Admin UI Screens
import AdminDashboard from './screens/admin/AdminDashboard';
import AdminLoginScreen from './screens/admin/AdminLoginScreen';
import PlayerProgressView from './screens/admin/PlayerProgressView';
import RewardResetScreen from './screens/admin/RewardResetScreen';
import NumberRangeSettingsScreen from './screens/admin/NumberRangeSettingsScreen';

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="App">
          <Routes>
            {/* Child UI Routes */}
            <Route path="/" element={<PlayerSelectionScreen />} />
            <Route path="/games" element={<GameSelectionScreen />} />
            <Route path="/play" element={<GamePlayScreen />} />
            <Route path="/feedback" element={<FeedbackScreen />} />
            <Route path="/progress" element={<ProgressOverviewScreen />} />

            {/* Admin Login Route */}
            <Route path="/admin/login" element={<AdminLoginScreen />} />

            {/* Protected Admin UI Routes */}
            <Route
              path="/admin"
              element={
                <ProtectedRoute>
                  <AdminDashboard />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/player/:playerId"
              element={
                <ProtectedRoute>
                  <PlayerProgressView />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/reset-rewards/:playerId"
              element={
                <ProtectedRoute>
                  <RewardResetScreen />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin/number-ranges/:playerId"
              element={
                <ProtectedRoute>
                  <NumberRangeSettingsScreen />
                </ProtectedRoute>
              }
            />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;
