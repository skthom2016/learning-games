import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { api } from '../../api/client';
import './GameSelectionScreen.css';

/**
 * GameSelectionScreen - Choose which game to play
 * See ARCHITECTURE.md Section 6 - Screen 2
 *
 * TODO:
 * - Get player ID from navigation state
 * - Fetch player data to show name and star count
 * - Display available games (MVP: only multiplication)
 * - Show "My Progress" button
 * - On "Let's Play" click, start session and navigate to /play
 */
function GameSelectionScreen() {
  const navigate = useNavigate();
  const location = useLocation();
  const { playerId } = location.state || {};

  const [player, setPlayer] = useState(null);
  const [starCount, setStarCount] = useState(0);

  useEffect(() => {
    if (!playerId) {
      navigate('/');
      return;
    }
    loadPlayerData();
  }, [playerId, navigate]);

  const loadPlayerData = async () => {
    try {
      const playerData = await api.getPlayer(playerId);
      const rewardData = await api.getRewardBalance(playerId);

      setPlayer(playerData.data);
      setStarCount(rewardData.data.total_stars || 0);
    } catch (error) {
      console.error('Failed to load player data:', error);
    }
  };

  const handlePlayGame = async (gameId) => {
    try {
      await api.startSession(playerId, gameId);
      navigate('/play', { state: { playerId, gameId } });
    } catch (error) {
      console.error('Failed to start session:', error);
    }
  };

  const handleViewProgress = () => {
    navigate('/progress', { state: { playerId } });
  };

  if (!player) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="game-selection-screen">
      <div className="header">
        <h1>Hi {player.player_name}! 👋</h1>
        <div className="star-count">⭐ {starCount}</div>
      </div>

      <h2 className="subtitle">What do you want to play?</h2>

      <div className="game-cards">
        {/* Addition game */}
        <div className="game-card">
          <div className="game-icon">➕</div>
          <h3 className="game-title">Addition Adventure</h3>
          <p className="game-description">Master addition and become a math explorer!</p>
          <button
            className="play-button"
            onClick={() => handlePlayGame('addition-facts')}
          >
            Let's Play!
          </button>
        </div>

        {/* Subtraction game */}
        <div className="game-card">
          <div className="game-icon">➖</div>
          <h3 className="game-title">Subtraction Safari</h3>
          <p className="game-description">Explore subtraction and become a math ranger!</p>
          <button
            className="play-button"
            onClick={() => handlePlayGame('subtraction-facts')}
          >
            Let's Play!
          </button>
        </div>

        {/* Multiplication game */}
        <div className="game-card">
          <div className="game-icon">✖️</div>
          <h3 className="game-title">Multiplication Master</h3>
          <p className="game-description">Practice times tables with fun characters!</p>
          <button
            className="play-button"
            onClick={() => handlePlayGame('multiplication-tables')}
          >
            Let's Play!
          </button>
        </div>

        {/* Division game */}
        <div className="game-card">
          <div className="game-icon">➗</div>
          <h3 className="game-title">Division Detective</h3>
          <p className="game-description">Solve division puzzles and become a math detective!</p>
          <button
            className="play-button"
            onClick={() => handlePlayGame('division-facts')}
          >
            Let's Play!
          </button>
        </div>
      </div>

      <div className="footer-buttons">
        <button onClick={() => navigate('/')}>← Back</button>
        <button onClick={handleViewProgress}>My Progress</button>
      </div>
    </div>
  );
}

export default GameSelectionScreen;
