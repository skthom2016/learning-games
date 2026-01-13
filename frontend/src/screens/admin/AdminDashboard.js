import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api/client';
import { useAuth } from '../../contexts/AuthContext';
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
  const { logout } = useAuth();
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPlayers();
  }, []);

  const loadPlayers = async () => {
    try {
      console.log('Fetching players from API...');
      const response = await api.listPlayers();
      console.log('Players API response:', response.data);

      const playersList = Array.isArray(response.data) ? response.data : [];

      // Fetch additional stats for each player
      const playersWithStats = await Promise.all(
        playersList.map(async (player) => {
          try {
            // Get reward balance
            const balanceResponse = await api.getRewardBalance(player.player_id);
            const totalStars = balanceResponse.data?.total_stars || 0;

            // Format last played date
            let lastPlayed = 'Never';
            if (player.last_played_at) {
              const lastPlayedDate = new Date(player.last_played_at);
              const now = new Date();
              const diffMs = now - lastPlayedDate;
              const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

              if (diffDays === 0) {
                lastPlayed = 'Today';
              } else if (diffDays === 1) {
                lastPlayed = 'Yesterday';
              } else if (diffDays < 7) {
                lastPlayed = `${diffDays} days ago`;
              } else if (diffDays < 30) {
                const weeks = Math.floor(diffDays / 7);
                lastPlayed = `${weeks} week${weeks > 1 ? 's' : ''} ago`;
              } else {
                const months = Math.floor(diffDays / 30);
                lastPlayed = `${months} month${months > 1 ? 's' : ''} ago`;
              }
            }

            return {
              player_id: player.player_id,
              player_name: player.player_name,
              total_stars: totalStars,
              last_played: lastPlayed,
              mastery_summary: 'View progress', // Simplified for now
            };
          } catch (error) {
            console.error(`Failed to load stats for player ${player.player_id}:`, error);
            return {
              player_id: player.player_id,
              player_name: player.player_name,
              total_stars: 0,
              last_played: 'Never',
              mastery_summary: 'No data',
            };
          }
        })
      );

      console.log('Players with stats:', playersWithStats);
      setPlayers(playersWithStats);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load players:', error);
      setPlayers([]);
      setLoading(false);
    }
  };

  const handleViewProgress = (playerId) => {
    navigate(`/admin/player/${playerId}`);
  };

  const handleResetRewards = (playerId) => {
    navigate(`/admin/reset-rewards/${playerId}`);
  };

  const handleDeletePlayer = async (playerId, playerName) => {
    if (
      !window.confirm(
        `Are you sure you want to delete ${playerName}?\n\nThis will permanently delete:\n- All progress data\n- All reward history\n- All mastery records\n\nThis action cannot be undone!`
      )
    ) {
      return;
    }

    // Double confirmation
    if (!window.confirm(`This is your last chance! Delete "${playerName}"?`)) {
      return;
    }

    try {
      console.log(`Deleting player ${playerId} (${playerName})...`);
      const response = await api.deletePlayer(playerId);
      console.log('Delete response:', response);

      // Show success message
      alert(`${playerName} has been successfully deleted.`);

      // Reload players list
      await loadPlayers();
    } catch (error) {
      console.error('Failed to delete player:', error);

      // Show detailed error message
      let errorMessage = 'Failed to delete player.\n\n';

      if (error.response) {
        errorMessage += `Status: ${error.response.status}\n`;
        errorMessage += `Error: ${JSON.stringify(error.response.data)}`;
      } else if (error.request) {
        errorMessage += 'No response from server. Please check if the backend is running.';
      } else {
        errorMessage += `Error: ${error.message}`;
      }

      alert(errorMessage);
    }
  };

  const handleLogout = () => {
    if (window.confirm('Are you sure you want to logout?')) {
      logout();
      navigate('/');
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="admin-dashboard">
      <header className="admin-header">
        <h1>Learning Game Platform — Admin Console</h1>
        <div className="header-buttons">
          <button className="logout-button" onClick={handleLogout}>
            🚪 Logout
          </button>
          <button onClick={() => navigate('/')}>← Back to Game</button>
        </div>
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
                <button
                  className="delete-button"
                  onClick={() => handleDeletePlayer(player.player_id, player.player_name)}
                >
                  Delete
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
