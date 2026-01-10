import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';

// Child UI Screens
import PlayerSelectionScreen from './screens/child/PlayerSelectionScreen';
import GameSelectionScreen from './screens/child/GameSelectionScreen';
import GamePlayScreen from './screens/child/GamePlayScreen';
import FeedbackScreen from './screens/child/FeedbackScreen';
import ProgressOverviewScreen from './screens/child/ProgressOverviewScreen';

// Admin UI Screens
import AdminDashboard from './screens/admin/AdminDashboard';
import PlayerProgressView from './screens/admin/PlayerProgressView';
import RewardResetScreen from './screens/admin/RewardResetScreen';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          {/* Child UI Routes */}
          <Route path="/" element={<PlayerSelectionScreen />} />
          <Route path="/games" element={<GameSelectionScreen />} />
          <Route path="/play" element={<GamePlayScreen />} />
          <Route path="/feedback" element={<FeedbackScreen />} />
          <Route path="/progress" element={<ProgressOverviewScreen />} />

          {/* Admin UI Routes */}
          <Route path="/admin" element={<AdminDashboard />} />
          <Route path="/admin/progress/:playerId" element={<PlayerProgressView />} />
          <Route path="/admin/reset-rewards/:playerId" element={<RewardResetScreen />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
