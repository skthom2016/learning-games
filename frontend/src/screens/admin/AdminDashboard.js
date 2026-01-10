import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api/client';
import './AdminDashboard.css';

/**
 * AdminDashboard - Admin overview of all players
 * See ARCHITECTURE.md Section 6 - Screen 8
 *
 * TODO:
 * - Fetch list of all players
 * - Display player cards with:
 *   - Name
 *   - Star count
 *   - Last played timestamp
 *   - Mastery summary (e.g., "3/11 topics mastered")
 * - Add "View Progress" button per player
 * - Add "Reset Rewards" button per player
 * - Add "Add New Player" button
 */
function AdminDashboard() {
  const navigate = useNavigate();
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPlayers();
  }, []);

  const loadPlayers = async () => {
    try {
      // TODO: Fetch players and their stats
      // const playersData = await api.listPlayers();
      // For each player, fetch reward balance and mastery summary

      // Mock data
      setPlayers([
        {
          player_id: '1',
          player_name: 'Emma',
          total_stars: 1234,
          last_played: 'Today',
          mastery_summary: '3/11 mastered',
        },
        {
          player_id: '2',
          player_name: 'Oliver',
          total_stars: 856,
          last_played: '2 days ago',
          mastery_summary: '1/11 mastered',
        },
      ]);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load players:', error);
      setLoading(false);
    }
  };

  const handleViewProgress = (playerId) => {
    navigate(`/admin/progress/${playerId}`);
  };

  const handleResetRewards = (playerId) => {
    navigate(`/admin/reset-rewards/${playerId}`);
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="admin-dashboard">
      <header className="admin-header">
        <h1>Learning Game Platform — Admin Console</h1>
        <button onClick={() => navigate('/')}>← Back to Game</button>
      </header>

      <div className="players-section">
        <h2>Players:</h2>

        <div className="player-cards-admin">
          {players.map((player) => (
            <div key={player.player_id} className="player-card-admin">
              <div className="player-info">
                <h3>{player.player_name}</h3>
                <p>⭐ {player.total_stars} stars</p>
                <p>Last played: {player.last_played}</p>
                <p>Times Tables: {player.mastery_summary}</p>
              </div>
              <div className="player-actions">
                <button onClick={() => handleViewProgress(player.player_id)}>
                  View Progress
                </button>
                <button onClick={() => handleResetRewards(player.player_id)}>
                  Reset Rewards
                </button>
              </div>
            </div>
          ))}
        </div>

        {/* TODO: Add "Add New Player" button */}
      </div>
    </div>
  );
}

export default AdminDashboard;
