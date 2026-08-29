import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api/client';

// Default star rewards matching backend defaults
const DEFAULT_REWARDS = {
  easy: 5,
  medium: 10,
  hard: 15,
};

function StarRewardSettingsScreen() {
  const navigate = useNavigate();
  const { playerId } = useParams();
  const [player, setPlayer] = useState(null);
  const [rewards, setRewards] = useState({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  // Game configurations
  const games = [
    { id: 'multiplication-tables', name: 'Multiplication', icon: '×' },
    { id: 'division-facts', name: 'Division', icon: '÷' },
    { id: 'addition-facts', name: 'Addition', icon: '+' },
    { id: 'subtraction-facts', name: 'Subtraction', icon: '−' },
  ];

  useEffect(() => {
    loadData();
  }, [playerId]);

  const loadData = async () => {
    try {
      const [playerResponse, rewardsResponse] = await Promise.all([
        api.getPlayer(playerId),
        api.getPlayerStarRewards(playerId),
      ]);

      setPlayer(playerResponse.data);

      // Convert rewards response to map
      const rewardsMap = {};
      if (rewardsResponse.data.star_rewards) {
        Object.entries(rewardsResponse.data.star_rewards).forEach(([gameId, config]) => {
          // Check if values differ from defaults to determine if custom
          const isCustom =
            config.easy_stars !== DEFAULT_REWARDS.easy ||
            config.medium_stars !== DEFAULT_REWARDS.medium ||
            config.hard_stars !== DEFAULT_REWARDS.hard;

          rewardsMap[gameId] = {
            easy: config.easy_stars,
            medium: config.medium_stars,
            hard: config.hard_stars,
            enabled: isCustom,
          };
        });
      }
      setRewards(rewardsMap);
      setLoading(false);
    } catch (error) {
      console.error('Failed to load data:', error);
      setLoading(false);
    }
  };

  const handleStarChange = (gameId, difficulty, value) => {
    const intValue = value === '' ? 1 : Math.max(1, parseInt(value, 10) || 1);
    setRewards((prev) => ({
      ...prev,
      [gameId]: {
        ...prev[gameId],
        [difficulty]: intValue,
        enabled: true,
      },
    }));
  };

  const handleToggleEnabled = (gameId) => {
    setRewards((prev) => {
      const current = prev[gameId];
      if (current && current.enabled) {
        // Disable - remove the custom config (revert to defaults)
        const { [gameId]: _, ...rest } = prev;
        return rest;
      } else {
        // Enable - add default reward config
        return {
          ...prev,
          [gameId]: { easy: DEFAULT_REWARDS.easy, medium: DEFAULT_REWARDS.medium, hard: DEFAULT_REWARDS.hard, enabled: true },
        };
      }
    });
  };

  const handleSave = async () => {
    // Validate rewards
    for (const [gameId, config] of Object.entries(rewards)) {
      const gameName = games.find((g) => g.id === gameId)?.name;
      if (config.easy < 1 || config.medium < 1 || config.hard < 1) {
        alert(`${gameName}: Star values must be at least 1`);
        return;
      }
    }

    if (!window.confirm('Save these star reward settings? This will affect the stars earned for correct answers.')) {
      return;
    }

    try {
      setSaving(true);
      const starRewards = {};
      for (const [gameId, config] of Object.entries(rewards)) {
        starRewards[gameId] = {
          easy_stars: config.easy,
          medium_stars: config.medium,
          hard_stars: config.hard,
        };
      }

      await api.setPlayerStarRewards(playerId, starRewards);
      alert('Star reward settings saved successfully!');
    } catch (error) {
      console.error('Failed to save rewards:', error);
      alert('Failed to save star reward settings');
    } finally {
      setSaving(false);
    }
  };

  const handleResetToDefaults = async () => {
    if (!window.confirm('Reset all star rewards to system defaults?')) {
      return;
    }

    try {
      setSaving(true);
      for (const gameId of Object.keys(rewards)) {
        await api.deletePlayerStarReward(playerId, gameId);
      }
      setRewards({});
      alert('Reset to defaults successfully!');
    } catch (error) {
      console.error('Failed to reset:', error);
      alert('Failed to reset to defaults');
    } finally {
      setSaving(false);
    }
  };

  if (loading || !player) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="star-reward-settings-screen">
      <style>{`
        .star-reward-settings-screen {
          max-width: 900px;
          margin: 0 auto;
          padding: 20px;
        }
        .admin-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 30px;
          padding-bottom: 20px;
          border-bottom: 1px solid #e0e0e0;
        }
        .admin-header h1 {
          margin: 0;
          font-size: 1.5rem;
          color: #333;
        }
        .admin-header button {
          padding: 8px 16px;
          background: #f0f0f0;
          border: none;
          border-radius: 4px;
          cursor: pointer;
          font-size: 0.9rem;
        }
        .admin-header button:hover {
          background: #e0e0e0;
        }
        .info-section {
          background: #f8f9fa;
          padding: 20px;
          border-radius: 8px;
          margin-bottom: 30px;
        }
        .info-section h3 {
          margin-top: 0;
          color: #495057;
        }
        .info-section p {
          margin-bottom: 10px;
          color: #6c757d;
        }
        .games-section {
          margin-bottom: 30px;
        }
        .games-section h2 {
          margin-bottom: 20px;
          color: #333;
        }
        .games-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
          gap: 20px;
        }
        .game-card {
          background: white;
          border: 2px solid #e0e0e0;
          border-radius: 12px;
          padding: 20px;
          transition: all 0.2s;
        }
        .game-card.enabled {
          border-color: #4caf50;
          box-shadow: 0 4px 12px rgba(76, 175, 80, 0.15);
        }
        .game-card-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 15px;
        }
        .game-icon {
          font-size: 1.5rem;
          font-weight: bold;
          color: #666;
          width: 40px;
          height: 40px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: #f5f5f5;
          border-radius: 50%;
        }
        .game-card-header h3 {
          margin: 0;
          flex: 1;
          margin-left: 10px;
          font-size: 1.1rem;
        }
        .toggle-switch {
          position: relative;
          display: inline-block;
          width: 50px;
          height: 26px;
        }
        .toggle-switch input {
          opacity: 0;
          width: 0;
          height: 0;
        }
        .toggle-switch .slider {
          position: absolute;
          cursor: pointer;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background-color: #ccc;
          transition: 0.3s;
          border-radius: 26px;
        }
        .toggle-switch .slider:before {
          position: absolute;
          content: "";
          height: 20px;
          width: 20px;
          left: 3px;
          bottom: 3px;
          background-color: white;
          transition: 0.3s;
          border-radius: 50%;
        }
        .toggle-switch input:checked + .slider {
          background-color: #4caf50;
        }
        .toggle-switch input:checked + .slider:before {
          transform: translateX(24px);
        }
        .star-inputs {
          display: flex;
          flex-direction: column;
          gap: 15px;
        }
        .star-row {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 10px;
          background: #f9f9f9;
          border-radius: 8px;
        }
        .star-row label {
          display: flex;
          align-items: center;
          gap: 8px;
          font-weight: 500;
        }
        .difficulty-badge {
          padding: 4px 12px;
          border-radius: 12px;
          font-size: 0.8rem;
          font-weight: bold;
        }
        .difficulty-badge.easy {
          background: #e8f5e9;
          color: #2e7d32;
        }
        .difficulty-badge.medium {
          background: #fff3e0;
          color: #ef6c00;
        }
        .difficulty-badge.hard {
          background: #ffebee;
          color: #c62828;
        }
        .star-input-group {
          display: flex;
          align-items: center;
          gap: 8px;
        }
        .star-input-group input {
          width: 70px;
          padding: 6px 10px;
          border: 1px solid #ddd;
          border-radius: 4px;
          text-align: center;
          font-size: 1rem;
        }
        .star-icon {
          color: #ffb300;
          font-size: 1.2rem;
        }
        .range-disabled {
          text-align: center;
          padding: 20px;
          color: #999;
          font-style: italic;
        }
        .default-badge {
          display: inline-block;
          padding: 2px 8px;
          background: #e0e0e0;
          border-radius: 4px;
          font-size: 0.75rem;
          color: #666;
        }
        .actions-section {
          display: flex;
          justify-content: flex-end;
          gap: 15px;
          padding-top: 20px;
          border-top: 1px solid #e0e0e0;
        }
        .primary-button, .secondary-button {
          padding: 12px 24px;
          border: none;
          border-radius: 6px;
          font-size: 1rem;
          cursor: pointer;
          transition: all 0.2s;
        }
        .primary-button {
          background: #4caf50;
          color: white;
        }
        .primary-button:hover:not(:disabled) {
          background: #45a049;
        }
        .primary-button:disabled {
          background: #ccc;
          cursor: not-allowed;
        }
        .secondary-button {
          background: #f0f0f0;
          color: #333;
        }
        .secondary-button:hover:not(:disabled) {
          background: #e0e0e0;
        }
        .secondary-button:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }
      `}</style>

      <header className="admin-header">
        <h1>Star Reward Settings for {player.player_name}</h1>
        <button onClick={() => navigate(`/admin/player/${playerId}`)}>← Back to Player</button>
      </header>

      <div className="settings-content">
        <section className="info-section">
          <h3>About Star Rewards</h3>
          <p>
            Configure custom star rewards for each difficulty level based on the player's age and ability.
            Stars are awarded for each correct answer.
          </p>
          <p>
            <strong>Default rewards:</strong> Easy = 5 stars, Medium = 10 stars, Hard = 15 stars
          </p>
          <p>
            <strong>Difficulty Progression:</strong> Players start at Easy level. After 10 consecutive correct answers,
            they advance to Medium. After 10 more consecutive correct, they advance to Hard. Two consecutive wrong
            answers will drop them down one level.
          </p>
        </section>

        <section className="games-section">
          <h2>Star Rewards per Game</h2>
          <div className="games-grid">
            {games.map((game) => {
              const config = rewards[game.id];
              const isEnabled = config && config.enabled;

              return (
                <div key={game.id} className={`game-card ${isEnabled ? 'enabled' : ''}`}>
                  <div className="game-card-header">
                    <span className="game-icon">{game.icon}</span>
                    <h3>{game.name}</h3>
                    <label className="toggle-switch">
                      <input
                        type="checkbox"
                        checked={isEnabled}
                        onChange={() => handleToggleEnabled(game.id)}
                      />
                      <span className="slider"></span>
                    </label>
                  </div>

                  {isEnabled ? (
                    <div className="star-inputs">
                      <div className="star-row">
                        <label>
                          <span className="difficulty-badge easy">Easy</span>
                        </label>
                        <div className="star-input-group">
                          <input
                            type="number"
                            min="1"
                            value={config.easy}
                            onChange={(e) => handleStarChange(game.id, 'easy', e.target.value)}
                          />
                          <span className="star-icon">⭐</span>
                        </div>
                      </div>

                      <div className="star-row">
                        <label>
                          <span className="difficulty-badge medium">Medium</span>
                        </label>
                        <div className="star-input-group">
                          <input
                            type="number"
                            min="1"
                            value={config.medium}
                            onChange={(e) => handleStarChange(game.id, 'medium', e.target.value)}
                          />
                          <span className="star-icon">⭐</span>
                        </div>
                      </div>

                      <div className="star-row">
                        <label>
                          <span className="difficulty-badge hard">Hard</span>
                        </label>
                        <div className="star-input-group">
                          <input
                            type="number"
                            min="1"
                            value={config.hard}
                            onChange={(e) => handleStarChange(game.id, 'hard', e.target.value)}
                          />
                          <span className="star-icon">⭐</span>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="range-disabled">
                      Using defaults: Easy 5⭐ | Medium 10⭐ | Hard 15⭐
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </section>

        <section className="actions-section">
          <button className="secondary-button" onClick={handleResetToDefaults} disabled={saving || Object.keys(rewards).length === 0}>
            Reset to Defaults
          </button>
          <button className="primary-button" onClick={handleSave} disabled={saving || Object.keys(rewards).length === 0}>
            {saving ? 'Saving...' : 'Save Settings'}
          </button>
        </section>
      </div>
    </div>
  );
}

export default StarRewardSettingsScreen;
