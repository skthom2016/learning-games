import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api/client';
import './RewardResetScreen.css';

/**
 * RewardResetScreen - Admin consolidates virtual rewards into real rewards
 * See ARCHITECTURE.md Section 6 - Screen 11
 *
 * TODO:
 * - Fetch current unconsolidated star balance
 * - Input: number of stars to consolidate (validate <= balance)
 * - Input: description of real-world reward
 * - On submit, call API to create RewardResetEvent
 * - Show previous consolidation history
 * - Confirm before resetting
 */
function RewardResetScreen() {
  const navigate = useNavigate();
  const { playerId } = useParams();
  const [player, setPlayer] = useState(null);
  const [balance, setBalance] = useState({ total: 0, unconsolidated: 0 });
  const [starsToReset, setStarsToReset] = useState('');
  const [rewardDescription, setRewardDescription] = useState('');
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    loadData();
  }, [playerId]);

  const loadData = async () => {
    try {
      // TODO: Fetch player, reward balance, and reset history
      // const playerData = await api.getPlayer(playerId);
      // const balanceData = await api.getRewardBalance(playerId);

      // Mock data
      setPlayer({ player_name: 'Emma' });
      setBalance({ total: 1234, unconsolidated: 856 });
      setHistory([
        { date: '2025-12-20', stars: 500, reward: 'New puzzle book' },
        { date: '2025-11-15', stars: 350, reward: 'Movie night' },
      ]);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load data:', error);
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const stars = parseInt(starsToReset, 10);
    if (isNaN(stars) || stars <= 0) {
      alert('Please enter a valid number of stars');
      return;
    }
    if (stars > balance.unconsolidated) {
      alert(`Cannot consolidate more than ${balance.unconsolidated} stars`);
      return;
    }
    if (!rewardDescription.trim()) {
      alert('Please describe the real-world reward');
      return;
    }

    if (!window.confirm(`Confirm: Reset ${stars} stars for "${rewardDescription}"?`)) {
      return;
    }

    try {
      setSubmitting(true);
      // TODO: Call API to reset rewards
      // await api.resetRewards(playerId, stars, rewardDescription);

      alert(`Successfully consolidated ${stars} stars!`);
      setStarsToReset('');
      setRewardDescription('');
      loadData(); // Reload balance and history
    } catch (error) {
      console.error('Failed to reset rewards:', error);
      alert('Failed to reset rewards');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading || !player) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="reward-reset-screen">
      <header className="admin-header">
        <h1>Reset Rewards for {player.player_name}</h1>
        <button onClick={() => navigate('/admin')}>← Back to Dashboard</button>
      </header>

      <div className="reset-content">
        <section className="balance-section">
          <h2>Current Balance:</h2>
          <div className="balance-display">
            <p>⭐ Unconsolidated Stars: <strong>{balance.unconsolidated}</strong></p>
            <p className="balance-note">(Total earned: {balance.total})</p>
          </div>
        </section>

        <section className="reset-form-section">
          <h2>Consolidation:</h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>How many stars to consolidate?</label>
              <input
                type="number"
                min="1"
                max={balance.unconsolidated}
                value={starsToReset}
                onChange={(e) => setStarsToReset(e.target.value)}
                placeholder={`Max: ${balance.unconsolidated}`}
                required
              />
            </div>

            <div className="form-group">
              <label>What real-world reward did you give?</label>
              <input
                type="text"
                value={rewardDescription}
                onChange={(e) => setRewardDescription(e.target.value)}
                placeholder="e.g., Ice cream trip to Dairy Queen"
                maxLength={500}
                required
              />
            </div>

            <div className="form-actions">
              <button type="button" onClick={() => navigate('/admin')}>
                Cancel
              </button>
              <button type="submit" disabled={submitting}>
                {submitting ? 'Confirming...' : 'Confirm Reset'}
              </button>
            </div>
          </form>
        </section>

        <section className="history-section">
          <h2>Previous Consolidations:</h2>
          {history.length === 0 ? (
            <p className="no-history">No previous consolidations</p>
          ) : (
            <ul className="history-list">
              {history.map((item, i) => (
                <li key={i}>
                  • {item.date}: {item.stars} stars → "{item.reward}"
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </div>
  );
}

export default RewardResetScreen;
