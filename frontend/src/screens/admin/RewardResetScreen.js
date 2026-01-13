import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api/client';
import './RewardResetScreen.css';

/**
 * RewardResetScreen - Admin resets/archives virtual rewards
 * Simplified flow: Show balance, confirm, reset all unconsolidated stars
 */
function RewardResetScreen() {
  const navigate = useNavigate();
  const { playerId } = useParams();
  const [player, setPlayer] = useState(null);
  const [balance, setBalance] = useState({ total: 0, unconsolidated: 0 });
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [resetting, setResetting] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadData();
  }, [playerId]);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch player, reward balance, and analytics (which includes consolidation history)
      const [playerResponse, balanceResponse, analyticsResponse] = await Promise.all([
        api.getPlayer(playerId),
        api.getRewardBalance(playerId),
        api.getPlayerAnalytics(playerId),
      ]);

      setPlayer(playerResponse.data);
      setBalance({
        total: balanceResponse.data.total_stars || 0,
        unconsolidated: balanceResponse.data.unconsolidated_stars || 0,
      });

      // Get consolidation history from analytics
      const consolidations = analyticsResponse.data.consolidations || [];
      setHistory(consolidations.map(c => ({
        date: new Date(c.reset_at).toLocaleDateString(),
        stars: c.stars_consolidated,
        reward: c.real_world_reward,
      })));

      setLoading(false);
    } catch (err) {
      console.error('Failed to load data:', err);
      setError('Failed to load player data. Please try again.');
      setLoading(false);
    }
  };

  const handleReset = async () => {
    if (balance.unconsolidated === 0) {
      alert('No stars to reset!');
      return;
    }

    // Single confirmation
    if (!window.confirm(
      `Reset ${balance.unconsolidated} stars for ${player.player_name}?\n\n` +
      `This will archive the current star balance. The stars will be marked as "redeemed" and won't count towards the visible balance anymore.\n\n` +
      `This action is recorded for reporting purposes.`
    )) {
      return;
    }

    try {
      setResetting(true);

      // Call API to reset rewards
      // Using a default reward description since we're simplifying the flow
      await api.resetRewards(
        playerId,
        balance.unconsolidated,
        `Reward reset on ${new Date().toLocaleDateString()}`
      );

      alert(`Successfully reset ${balance.unconsolidated} stars!`);

      // Reload data to show updated balance and history
      await loadData();
    } catch (err) {
      console.error('Failed to reset rewards:', err);

      let errorMessage = 'Failed to reset rewards. ';
      if (err.response?.data?.error) {
        errorMessage += err.response.data.error;
      } else {
        errorMessage += 'Please try again.';
      }

      alert(errorMessage);
    } finally {
      setResetting(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
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

  return (
    <div className="reward-reset-screen">
      <header className="admin-header">
        <h1>Reset Rewards for {player.player_name}</h1>
        <button onClick={() => navigate('/admin')}>← Back to Dashboard</button>
      </header>

      <div className="reset-content">
        <section className="balance-section">
          <h2>Current Star Balance</h2>
          <div className="balance-display">
            <div className="balance-main">
              <span className="star-icon">⭐</span>
              <span className="balance-value">{balance.unconsolidated}</span>
              <span className="balance-label">Available Stars</span>
            </div>
            <p className="balance-note">
              Total stars earned: {balance.total} |
              Already redeemed: {balance.total - balance.unconsolidated}
            </p>
          </div>
        </section>

        <section className="reset-action-section">
          <h2>Reset Stars</h2>
          <p className="reset-info">
            Resetting stars will archive the current balance and mark all {balance.unconsolidated} stars
            as redeemed. This is useful when you've given the child a real-world reward.
          </p>

          <button
            className="reset-button"
            onClick={handleReset}
            disabled={resetting || balance.unconsolidated === 0}
          >
            {resetting ? 'Resetting...' : `Reset ${balance.unconsolidated} Stars`}
          </button>

          {balance.unconsolidated === 0 && (
            <p className="no-stars-message">No stars available to reset.</p>
          )}
        </section>

        <section className="history-section">
          <h2>Reset History</h2>
          {history.length === 0 ? (
            <p className="no-history">No previous resets</p>
          ) : (
            <div className="history-table">
              <div className="history-header">
                <span>Date</span>
                <span>Stars</span>
                <span>Note</span>
              </div>
              {history.map((item, i) => (
                <div key={i} className="history-row">
                  <span>{item.date}</span>
                  <span>{item.stars} ⭐</span>
                  <span>{item.reward}</span>
                </div>
              ))}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}

export default RewardResetScreen;
