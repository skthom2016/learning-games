import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { api } from '../../api/client';
import './ProgressOverviewScreen.css';

/**
 * ProgressOverviewScreen - Show child's progress (simplified)
 * See ARCHITECTURE.md Section 6 - Screen 7
 *
 * TODO:
 * - Fetch mastery heatmap from API
 * - Display total stars, badges (Phase 2), characters (Phase 2)
 * - Show topic list with simple icons (✓ mastered, ⭐ learning, ○ not started)
 * - Keep it simple and encouraging (no percentages)
 * - Add "Back to Games" button
 */
function ProgressOverviewScreen() {
  const navigate = useNavigate();
  const location = useLocation();
  const { playerId } = location.state || {};

  const [player, setPlayer] = useState(null);
  const [mastery, setMastery] = useState([]);
  const [starCount, setStarCount] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!playerId) {
      navigate('/');
      return;
    }
    loadProgress();
  }, [playerId, navigate]);

  const loadProgress = async () => {
    try {
      console.log('Loading progress for player:', playerId);
      const [playerData, masteryData, rewardData] = await Promise.all([
        api.getPlayer(playerId),
        api.getMasteryHeatmap(playerId, 'multiplication-tables'),
        api.getRewardBalance(playerId),
      ]);

      console.log('Player data:', playerData.data);
      console.log('Mastery data:', masteryData.data);
      console.log('Reward data:', rewardData.data);

      setPlayer(playerData.data);
      setStarCount(rewardData.data.total_stars || 0);
      setMastery(masteryData.data.topics || []);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load progress:', error);
      console.error('Error details:', error.response?.data || error.message);
      setLoading(false);
    }
  };

  const getMasteryIcon = (state) => {
    switch (state) {
      case 'MASTERED':
        return '✓';
      case 'STRONG':
      case 'LEARNING':
        return '⭐';
      case 'WEAK':
        return '⚠️';
      case 'UNKNOWN':
      default:
        return '○';
    }
  };

  const getMasteryLabel = (state) => {
    switch (state) {
      case 'MASTERED':
        return 'Mastered';
      case 'STRONG':
        return 'Doing great!';
      case 'LEARNING':
        return 'Learning';
      case 'WEAK':
        return 'Needs practice';
      case 'UNKNOWN':
      default:
        return 'Not started yet';
    }
  };

  if (loading || !player) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="progress-overview-screen">
      <h1>{player.player_name}'s Progress</h1>

      <div className="progress-summary">
        <div className="summary-item">
          <span className="summary-icon">⭐</span>
          <span className="summary-label">Total Stars:</span>
          <span className="summary-value">{starCount}</span>
        </div>
        {/* TODO: Add badges and characters in Phase 2 */}
      </div>

      <h2>Times Tables:</h2>
      <div className="topic-list">
        {mastery.map((topic) => (
          <div key={topic.topic_id} className="topic-item">
            <span className="topic-icon">{getMasteryIcon(topic.mastery_state)}</span>
            <span className="topic-name">{topic.display_name}</span>
            <span className="topic-status">{getMasteryLabel(topic.mastery_state)}</span>
          </div>
        ))}
      </div>

      <button className="back-button" onClick={() => navigate('/games', { state: { playerId } })}>
        Back to Games
      </button>
    </div>
  );
}

export default ProgressOverviewScreen;
