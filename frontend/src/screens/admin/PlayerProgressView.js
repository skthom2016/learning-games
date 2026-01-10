import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api/client';
import './PlayerProgressView.css';

/**
 * PlayerProgressView - Detailed analytics for a single player
 * See ARCHITECTURE.md Section 6 - Screen 9
 *
 * TODO:
 * - Fetch player progress and mastery heatmap
 * - Display summary stats (total questions, accuracy, session count)
 * - Show mastery heatmap with color coding:
 *   - Green: MASTERED, STRONG
 *   - Yellow: LEARNING
 *   - Red: WEAK
 *   - Gray: UNKNOWN
 * - Display difficulty distribution (pie chart or bars)
 * - Auto-generate recommendations (e.g., "6× table needs attention")
 * - Add "Export Report" button (Phase 2)
 */
function PlayerProgressView() {
  const navigate = useNavigate();
  const { playerId } = useParams();
  const [player, setPlayer] = useState(null);
  const [progress, setProgress] = useState(null);
  const [mastery, setMastery] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadProgress();
  }, [playerId]);

  const loadProgress = async () => {
    try {
      // TODO: Fetch player, progress, and mastery data
      // const playerData = await api.getPlayer(playerId);
      // const progressData = await api.getPlayerProgress(playerId, 'multiplication-tables');
      // const masteryData = await api.getMasteryHeatmap(playerId, 'multiplication-tables');

      // Mock data
      setPlayer({ player_name: 'Emma' });
      setProgress({
        total_questions: 487,
        total_correct: 370,
        accuracy: 0.76,
        session_count: 12,
        avg_session_duration: '15 min',
        confidence_recovery_percentage: 0.08,
      });
      setMastery([
        { topic_id: 'times-2', display_name: '2× table', mastery_state: 'MASTERED', recent_accuracy: 0.95, times_practiced: 45 },
        { topic_id: 'times-3', display_name: '3× table', mastery_state: 'MASTERED', recent_accuracy: 0.92, times_practiced: 38 },
        { topic_id: 'times-4', display_name: '4× table', mastery_state: 'LEARNING', recent_accuracy: 0.68, times_practiced: 22 },
        { topic_id: 'times-5', display_name: '5× table', mastery_state: 'STRONG', recent_accuracy: 0.82, times_practiced: 31 },
        { topic_id: 'times-6', display_name: '6× table', mastery_state: 'WEAK', recent_accuracy: 0.35, times_practiced: 18 },
        { topic_id: 'times-7', display_name: '7× table', mastery_state: 'LEARNING', recent_accuracy: 0.61, times_practiced: 25 },
        { topic_id: 'times-8', display_name: '8× table', mastery_state: 'LEARNING', recent_accuracy: 0.59, times_practiced: 20 },
        { topic_id: 'times-9', display_name: '9× table', mastery_state: 'UNKNOWN', recent_accuracy: 0, times_practiced: 2 },
        { topic_id: 'times-10', display_name: '10× table', mastery_state: 'UNKNOWN', recent_accuracy: 0, times_practiced: 1 },
        { topic_id: 'times-11', display_name: '11× table', mastery_state: 'UNKNOWN', recent_accuracy: 0, times_practiced: 0 },
        { topic_id: 'times-12', display_name: '12× table', mastery_state: 'UNKNOWN', recent_accuracy: 0, times_practiced: 0 },
      ]);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load progress:', error);
      setLoading(false);
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
    const weak = mastery.filter(m => m.mastery_state === 'WEAK');
    const unknown = mastery.filter(m => m.mastery_state === 'UNKNOWN');
    const strong = mastery.filter(m => m.mastery_state === 'MASTERED');

    if (weak.length > 0) {
      recs.push(`⚠️ ${weak[0].display_name} needs attention (${Math.round(weak[0].recent_accuracy * 100)}% accuracy)`);
    }
    if (unknown.length > 0 && strong.length >= 3) {
      recs.push(`💡 Ready to introduce ${unknown[0].display_name}`);
    }
    if (strong.length > 0) {
      recs.push(`✓ Great progress on ${strong.map(s => s.display_name.split('×')[0]).join(', ')}× tables!`);
    }

    return recs;
  };

  if (loading || !player || !progress) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="player-progress-view">
      <header className="admin-header">
        <h1>{player.player_name}'s Progress — Multiplication Tables</h1>
        <button onClick={() => navigate('/admin')}>← Back to Dashboard</button>
      </header>

      <div className="progress-content">
        <section className="summary-section">
          <h2>Summary:</h2>
          <div className="summary-grid">
            <div className="summary-stat">
              <span className="stat-label">Total Questions:</span>
              <span className="stat-value">{progress.total_questions}</span>
            </div>
            <div className="summary-stat">
              <span className="stat-label">Accuracy:</span>
              <span className="stat-value">{Math.round(progress.accuracy * 100)}%</span>
            </div>
            <div className="summary-stat">
              <span className="stat-label">Session Count:</span>
              <span className="stat-value">{progress.session_count}</span>
            </div>
            <div className="summary-stat">
              <span className="stat-label">Avg Session:</span>
              <span className="stat-value">{progress.avg_session_duration}</span>
            </div>
            <div className="summary-stat">
              <span className="stat-label">Confidence Recovery Time:</span>
              <span className="stat-value">{Math.round(progress.confidence_recovery_percentage * 100)}% (normal)</span>
            </div>
          </div>
        </section>

        <section className="heatmap-section">
          <h2>Mastery Heatmap — Multiplication Tables</h2>
          <div className="heatmap">
            {mastery.map((topic) => (
              <div key={topic.topic_id} className="heatmap-row">
                <span className="heatmap-icon">{getMasteryIcon(topic.mastery_state)}</span>
                <span className="heatmap-topic">{topic.display_name}:</span>
                <span className="heatmap-state" style={{ color: getMasteryColor(topic.mastery_state) }}>
                  {topic.mastery_state}
                </span>
                <span className="heatmap-stats">
                  ({topic.times_practiced > 0 ? `${Math.round(topic.recent_accuracy * 100)}% accuracy, ${topic.times_practiced} attempts` : `${topic.times_practiced} attempts`})
                </span>
              </div>
            ))}
          </div>
        </section>

        <section className="recommendations-section">
          <h2>Recommendations:</h2>
          <ul className="recommendations-list">
            {getRecommendations().map((rec, i) => (
              <li key={i}>{rec}</li>
            ))}
          </ul>
        </section>

        {/* TODO: Add difficulty distribution chart */}
        {/* TODO: Add "Export Report" button */}
      </div>
    </div>
  );
}

export default PlayerProgressView;
