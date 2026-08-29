import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../api/client';
import './PlayerSelectionScreen.css';

/**
 * PlayerSelectionScreen - Child selects their name to play
 * See ARCHITECTURE.md Section 6 - Screen 1
 *
 * TODO:
 * - Fetch list of players from API (api.listPlayers())
 * - Display player cards with name and character avatar
 * - Add "Add Player" button that opens modal/form
 * - On player click, navigate to /games with player ID
 * - Style with large, colorful, kid-friendly UI
 */
function PlayerSelectionScreen() {
  const navigate = useNavigate();
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showAddPlayer, setShowAddPlayer] = useState(false);
  const [newPlayerName, setNewPlayerName] = useState('');

  useEffect(() => {
    loadPlayers();
  }, []);

  const loadPlayers = async () => {
    try {
      const response = await api.listPlayers();
      // Ensure response.data is an array, fallback to empty array
      setPlayers(Array.isArray(response.data) ? response.data : []);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load players:', error);
      setPlayers([]); // Ensure players is always an array
      setLoading(false);
    }
  };

  const handlePlayerClick = (playerId) => {
    // TODO: Navigate to game selection with player context
    navigate('/games', { state: { playerId } });
  };

  const handleAddPlayer = async () => {
    if (!newPlayerName.trim()) {
      alert('Please enter a player name');
      return;
    }

    try {
      console.log('Creating player:', newPlayerName);
      const response = await api.createPlayer(newPlayerName);
      console.log('Player created response:', response);

      setShowAddPlayer(false);
      setNewPlayerName('');
      await loadPlayers();

      // Show success message
      alert(`${newPlayerName} has been added!`);
    } catch (error) {
      console.error('Failed to create player:', error);

      let errorMsg = 'Failed to add player';
      if (error.response) {
        errorMsg += `\n\nStatus: ${error.response.status}\nError: ${JSON.stringify(error.response.data)}`;
      } else if (error.message) {
        errorMsg += `\n\nError: ${error.message}`;
      }
      alert(errorMsg);
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="player-selection-screen">
      <h1 className="title">🎮 Let's Play & Learn! 🎮</h1>
      <h2 className="subtitle">Who's playing today?</h2>

      <div className="player-cards">
        {players && players.map((player) => (
          <div
            key={player.player_id}
            className="player-card"
            onClick={() => handlePlayerClick(player.player_id)}
          >
            <div className="player-avatar">
              {player.preferred_character_id || '👤'}
            </div>
            <div className="player-name">{player.player_name}</div>
          </div>
        ))}

        <div className="player-card add-player" onClick={() => setShowAddPlayer(true)}>
          <div className="player-avatar">➕</div>
          <div className="player-name">Add Player</div>
        </div>
      </div>

      {/* TODO: Implement proper modal component */}
      {showAddPlayer && (
        <div className="modal">
          <div className="modal-content">
            <h3>Add New Player</h3>
            <input
              type="text"
              value={newPlayerName}
              onChange={(e) => setNewPlayerName(e.target.value)}
              placeholder="Enter name"
              maxLength={50}
            />
            <div className="modal-buttons">
              <button onClick={handleAddPlayer}>Add</button>
              <button onClick={() => setShowAddPlayer(false)}>Cancel</button>
            </div>
          </div>
        </div>
      )}

      {/* Subtle admin link for parents */}
      <div className="admin-link">
        <button onClick={() => navigate('/admin/login')}>Admin</button>
      </div>
    </div>
  );
}

export default PlayerSelectionScreen;
