import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api/client';
import './PlayerProgressView.css';

/**
 * PlayerProgressView - Detailed analytics for a single player
 * Shows progress across all games with mastery heatmaps
 */
function PlayerProgressView() {
  const navigate = useNavigate();
  const { playerId } = useParams();
  const [player, setPlayer] = useState(null);
  const [analytics, setAnalytics] = useState(null);
  const [selectedGame, setSelectedGame] = useState('multiplication-tables');
  const [mastery, setMastery] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const games = [
    { id: 'multiplication-tables', name: 'Multiplication' },
    { id: 'division-facts', name: 'Division' },
    { id: 'addition-facts', name: 'Addition' },
    { id: 'subtraction-facts', name: 'Subtraction' },
  ];

  useEffect(() => {
    loadData();
  }, [playerId]);

  useEffect(() => {
    if (playerId && selectedGame) {
      loadMastery(selectedGame);
    }
  }, [playerId, selectedGame]);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch player and analytics data
      const [playerResponse, analyticsResponse] = await Promise.all([
        api.getPlayer(playerId),
        api.getPlayerAnalytics(playerId),
      ]);

      setPlayer(playerResponse.data);
      setAnalytics(analyticsResponse.data);
      setLoading(false);
    } catch (err) {
      console.error('Failed to load player data:', err);
      setError('Failed to load player data. Please try again.');
      setLoading(false);
    }
  };

  const loadMastery = async (gameId) => {
    try {
      const masteryResponse = await api.getMasteryHeatmap(playerId, gameId);
      setMastery(masteryResponse.data.mastery || []);
    } catch (err) {
      console.error('Failed to load mastery data:', err);
      // If no mastery data, show empty state
      setMastery([]);
    }
  };

  const getMasteryColor = (state) => {
    switch (state) {
      case 'MASTERED':
      case 'STRONG':
        return '#4caf50'; // Green
      case 'LEARNING':
        return '#ff9800'; // Yellow
      case 'WEAK':
        return '#f44336'; // Red
      case 'UNKNOWN':
      default:
        return '#9e9e9e'; // Gray
    }
  };

  const getMasteryIcon = (state) => {
    switch (state) {
      case 'MASTERED':
      case 'STRONG':
        return '🟢';
      case 'LEARNING':
        return '🟡';
      case 'WEAK':
        return '🔴';
      case 'UNKNOWN':
      default:
        return '⚪';
    }
  };

  const getRecommendations = () => {
    const recs = [];
    if (!mastery || mastery.length === 0) {
      recs.push('💡 Start playing to see recommendations!');
      return recs;
    }

    const weak = mastery.filter(m => m.mastery_state === 'WEAK');
    const unknown = mastery.filter(m => m.mastery_state === 'UNKNOWN');
    const strong = mastery.filter(m => m.mastery_state === 'MASTERED' || m.mastery_state === 'STRONG');

    if (weak.length > 0) {
      recs.push(`⚠️ ${weak[0].display_name} needs attention (${Math.round((weak[0].recent_accuracy || 0) * 100)}% accuracy)`);
    }
    if (unknown.length > 0 && strong.length >= 2) {
      recs.push(`💡 Ready to introduce ${unknown[0].display_name}`);
    }
    if (strong.length > 0) {
      const strongNames = strong.slice(0, 3).map(s => s.display_name).join(', ');
      recs.push(`✓ Great progress on ${strongNames}!`);
    }
    if (recs.length === 0) {
      recs.push('📚 Keep practicing to improve mastery!');
    }

    return recs;
  };

  const getGameStats = (gameId) => {
    if (!analytics || !analytics.games) return null;
    return analytics.games.find(g => g.game_id === gameId);
  };

  if (loading) {
    return <div className="loading">Loading player progress...</div>;
  }

  if (error) {
    return (
      <div className="error-screen">
        <p>{error}</p>
        <button onClick={() => navigate('/admin')}>← Back to Dashboard</button>
      </div>
    );
  }

  if (!player) {
    return (
      <div className="error-screen">
        <p>Player not found</p>
        <button onClick={() => navigate('/admin')}>← Back to Dashboard</button>
      </div>
    );
  }

  const currentGameStats = getGameStats(selectedGame);

  return (
    <div className="player-progress-view">
      <header className="admin-header">
        <h1>{player.player_name}'s Progress</h1>
        <div className="header-buttons">
          <button className="secondary-button" onClick={() => navigate(`/admin/number-ranges/${playerId}`)}>
            ⚙️ Number Ranges
          </button>
          <button onClick={() => navigate('/admin')}>← Back to Dashboard</button>
        </div>
      </header>

      <div className="progress-content">
        {/* Overall Summary Section */}
        <section className="summary-section">
          <h2>Overall Summary</h2>
          <div className="summary-grid">
            <div className="summary-stat">
              <span className="stat-label">Total Questions:</span>
              <span className="stat-value">{analytics?.total_questions || 0}</span>
            </div>
            <div className="summary-stat">
              <span className="stat-label">Overall Accuracy:</span>
              <span className="stat-value">{Math.round((analytics?.overall_accuracy || 0) * 100)}%</span>
            </div>
            <div className="summary-stat">
              <span className="stat-label">Total Correct:</span>
              <span className="stat-value">{analytics?.total_correct || 0}</span>
            </div>
          </div>
        </section>

        {/* Game Selector */}
        <section className="game-selector-section">
          <h2>Select Game</h2>
          <div className="game-tabs">
            {games.map((game) => {
              const stats = getGameStats(game.id);
              return (
                <button
                  key={game.id}
                  className={`game-tab ${selectedGame === game.id ? 'active' : ''}`}
                  onClick={() => setSelectedGame(game.id)}
                >
                  {game.name}
                  {stats && (
                    <span className="game-tab-stats">
                      ({Math.round((stats.accuracy || 0) * 100)}%)
                    </span>
                  )}
                </button>
              );
            })}
          </div>
        </section>

        {/* Game-specific Stats */}
        {currentGameStats && (
          <section className="game-stats-section">
            <h2>{games.find(g => g.id === selectedGame)?.name} Stats</h2>
            <div className="summary-grid">
              <div className="summary-stat">
                <span className="stat-label">Questions:</span>
                <span className="stat-value">{currentGameStats.total_questions}</span>
              </div>
              <div className="summary-stat">
                <span className="stat-label">Correct:</span>
                <span className="stat-value">{currentGameStats.total_correct}</span>
              </div>
              <div className="summary-stat">
                <span className="stat-label">Accuracy:</span>
                <span className="stat-value">{Math.round((currentGameStats.accuracy || 0) * 100)}%</span>
              </div>
              <div className="summary-stat">
                <span className="stat-label">Sessions:</span>
                <span className="stat-value">{currentGameStats.sessions}</span>
              </div>
            </div>
          </section>
        )}

        {/* Mastery Heatmap */}
        <section className="heatmap-section">
          <h2>Mastery Heatmap — {games.find(g => g.id === selectedGame)?.name}</h2>
          {mastery.length === 0 ? (
            <p className="no-data">No mastery data yet. Start playing to see progress!</p>
          ) : (
            <div className="heatmap">
              {mastery.map((topic) => (
                <div key={topic.topic_id} className="heatmap-row">
                  <span className="heatmap-icon">{getMasteryIcon(topic.mastery_state)}</span>
                  <span className="heatmap-topic">{topic.display_name}:</span>
                  <span className="heatmap-state" style={{ color: getMasteryColor(topic.mastery_state) }}>
                    {topic.mastery_state}
                  </span>
                  <span className="heatmap-stats">
                    ({topic.times_practiced > 0 ? `${Math.round((topic.recent_accuracy || 0) * 100)}% accuracy, ${topic.times_practiced} attempts` : `${topic.times_practiced || 0} attempts`})
                  </span>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Recommendations */}
        <section className="recommendations-section">
          <h2>Recommendations</h2>
          <ul className="recommendations-list">
            {getRecommendations().map((rec, i) => (
              <li key={i}>{rec}</li>
            ))}
          </ul>
        </section>

        {/* Reward History */}
        {analytics?.consolidations && analytics.consolidations.length > 0 && (
          <section className="history-section">
            <h2>Reward Redemption History</h2>
            <ul className="history-list">
              {analytics.consolidations.map((item, i) => (
                <li key={i}>
                  • {new Date(item.reset_at).toLocaleDateString()}: {item.stars_consolidated} stars → "{item.real_world_reward}"
                </li>
              ))}
            </ul>
          </section>
        )}
      </div>
    </div>
  );
}

export default PlayerProgressView;
