import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { api } from '../../api/client';
import './NumberRangeSettingsScreen.css';

/**
 * NumberRangeSettingsScreen - Admin configures number ranges per player per game
 * Supports granular operand-specific ranges for each operation type
 */
function NumberRangeSettingsScreen() {
  const navigate = useNavigate();
  const { playerId } = useParams();
  const [player, setPlayer] = useState(null);
  const [ranges, setRanges] = useState({});
  const [starRewards, setStarRewards] = useState({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  // Game configurations with operand labels
  const games = [
    {
      id: 'multiplication-tables',
      name: 'Multiplication',
      icon: '×',
      operand1Label: 'Factor 1 (first number)',
      operand2Label: 'Factor 2 (second number)',
    },
    {
      id: 'division-facts',
      name: 'Division',
      icon: '÷',
      operand1Label: 'Dividend (number being divided)',
      operand2Label: 'Divisor (number dividing by)',
    },
    {
      id: 'addition-facts',
      name: 'Addition',
      icon: '+',
      operand1Label: 'Addend 1 (first number)',
      operand2Label: 'Addend 2 (second number)',
    },
    {
      id: 'subtraction-facts',
      name: 'Subtraction',
      icon: '−',
      operand1Label: 'Minuend (starting number)',
      operand2Label: 'Subtrahend (number being subtracted)',
    },
  ];

  useEffect(() => {
    loadData();
  }, [playerId]);

  const loadData = async () => {
    try {
      const [playerResponse, rangesResponse, starRewardsResponse] = await Promise.all([
        api.getPlayer(playerId),
        api.getPlayerNumberRanges(playerId),
        api.getPlayerStarRewards(playerId),
      ]);

      setPlayer(playerResponse.data);

      // Convert ranges array to map for easier editing
      const rangesMap = {};
      if (rangesResponse.data.number_ranges) {
        rangesResponse.data.number_ranges.forEach((range) => {
          rangesMap[range.game_id] = {
            min: range.min_number,
            max: range.max_number,
            operand1Min: range.operand1_min,
            operand1Max: range.operand1_max,
            operand2Min: range.operand2_min,
            operand2Max: range.operand2_max,
            enabled: true,
            useGranular: !!(range.operand1_min || range.operand1_max || range.operand2_min || range.operand2_max),
          };
        });
      }
      setRanges(rangesMap);

      // Convert star rewards to map
      const starRewardsMap = {};
      if (starRewardsResponse.data.star_rewards) {
        Object.entries(starRewardsResponse.data.star_rewards).forEach(([gameId, config]) => {
          starRewardsMap[gameId] = {
            easy: config.easy_stars,
            medium: config.medium_stars,
            hard: config.hard_stars,
          };
        });
      }
      setStarRewards(starRewardsMap);

      setLoading(false);
    } catch (error) {
      console.error('Failed to load data:', error);
      setLoading(false);
    }
  };

  const handleRangeChange = (gameId, field, value) => {
    setRanges((prev) => ({
      ...prev,
      [gameId]: {
        ...prev[gameId],
        [field]: value === '' ? null : parseInt(value, 10),
      },
    }));
  };

  const handleToggleEnabled = (gameId) => {
    setRanges((prev) => {
      const current = prev[gameId];
      if (current && current.enabled) {
        // Disable - remove the range
        const { [gameId]: _, ...rest } = prev;
        return rest;
      } else {
        // Enable - add default range
        return {
          ...prev,
          [gameId]: { min: 1, max: 10, enabled: true, useGranular: false },
        };
      }
    });
  };

  const handleToggleGranular = (gameId) => {
    setRanges((prev) => {
      const current = prev[gameId];
      if (!current) return prev;

      return {
        ...prev,
        [gameId]: {
          ...current,
          useGranular: !current.useGranular,
          // Initialize granular values from global values when enabling
          operand1Min: !current.useGranular ? current.min : null,
          operand1Max: !current.useGranular ? current.max : null,
          operand2Min: !current.useGranular ? current.min : null,
          operand2Max: !current.useGranular ? current.max : null,
        },
      };
    });
  };

  const handleStarChange = (gameId, difficulty, value) => {
    const intValue = value === '' ? 1 : Math.max(1, parseInt(value, 10) || 1);
    setStarRewards((prev) => ({
      ...prev,
      [gameId]: {
        ...prev[gameId],
        [difficulty]: intValue,
      },
    }));
  };

  const handleSave = async () => {
    // Validate ranges
    for (const [gameId, range] of Object.entries(ranges)) {
      const gameName = games.find((g) => g.id === gameId)?.name;

      if (range.min >= range.max) {
        alert(`${gameName}: Minimum must be less than maximum`);
        return;
      }

      if (range.useGranular) {
        if (range.operand1Min != null && range.operand1Max != null && range.operand1Min >= range.operand1Max) {
          alert(`${gameName}: Operand 1 minimum must be less than maximum`);
          return;
        }
        if (range.operand2Min != null && range.operand2Max != null && range.operand2Min >= range.operand2Max) {
          alert(`${gameName}: Operand 2 minimum must be less than maximum`);
          return;
        }
      }
    }

    if (!window.confirm('Save these settings? This will affect the questions and star rewards for this player.')) {
      return;
    }

    try {
      setSaving(true);

      // Save number ranges
      const numberRanges = {};
      for (const [gameId, range] of Object.entries(ranges)) {
        numberRanges[gameId] = {
          min_number: range.min,
          max_number: range.max,
        };

        // Add operand-specific ranges if enabled
        if (range.useGranular) {
          if (range.operand1Min != null) numberRanges[gameId].operand1_min = range.operand1Min;
          if (range.operand1Max != null) numberRanges[gameId].operand1_max = range.operand1Max;
          if (range.operand2Min != null) numberRanges[gameId].operand2_min = range.operand2Min;
          if (range.operand2Max != null) numberRanges[gameId].operand2_max = range.operand2Max;
        }
      }

      // Save star rewards
      const starRewardsPayload = {};
      for (const [gameId, config] of Object.entries(starRewards)) {
        starRewardsPayload[gameId] = {
          easy_stars: config.easy,
          medium_stars: config.medium,
          hard_stars: config.hard,
        };
      }

      // Save both in parallel
      await Promise.all([
        Object.keys(numberRanges).length > 0 ? api.setPlayerNumberRanges(playerId, numberRanges) : Promise.resolve(),
        Object.keys(starRewardsPayload).length > 0 ? api.setPlayerStarRewards(playerId, starRewardsPayload) : Promise.resolve(),
      ]);

      alert('Settings saved successfully!');
    } catch (error) {
      console.error('Failed to save settings:', error);
      alert('Failed to save settings');
    } finally {
      setSaving(false);
    }
  };

  const handleResetToDefaults = async () => {
    if (!window.confirm('Reset all ranges to defaults (remove all custom settings)?')) {
      return;
    }

    // Delete all ranges for each game
    try {
      setSaving(true);
      for (const gameId of Object.keys(ranges)) {
        await api.deletePlayerNumberRange(playerId, gameId);
      }
      setRanges({});
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
    <div className="number-range-settings-screen">
      <header className="admin-header">
        <h1>Number Range & Star Reward Settings for {player.player_name}</h1>
        <button onClick={() => navigate(`/admin/player/${playerId}`)}>← Back to Player</button>
      </header>

      <div className="settings-content">
        <section className="info-section">
          <p>
            Set custom number ranges and star rewards for each game to customize difficulty and motivation
            based on the player's age and ability. You can set a simple range for all numbers, or use
            advanced settings to set different ranges for each operand. Star rewards are given for each
            correct answer and vary by difficulty level (Easy/Medium/Hard).
          </p>
          <div className="grade-examples">
            <strong>Grade Examples:</strong>
            <ul>
              <li>Kindergarten/1st Grade: Min 1, Max 10</li>
              <li>2nd Grade: Min 1, Max 20</li>
              <li>3rd Grade: Min 1, Max 100</li>
              <li>4th Grade: Min 1, Max 500</li>
              <li>5th Grade+: Min 1, Max 1000</li>
            </ul>
          </div>
        </section>

        <section className="games-section">
          <h2>Game Number Ranges</h2>
          <div className="games-grid">
            {games.map((game) => {
              const range = ranges[game.id];
              const isEnabled = range && range.enabled;

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
                    <div className="range-inputs">
                      {/* Global Range */}
                      <div className="range-section">
                        <h4>Basic Range (all numbers)</h4>
                        <div className="input-row">
                          <div className="input-group">
                            <label>Min:</label>
                            <input
                              type="number"
                              min="0"
                              value={range.min}
                              onChange={(e) => handleRangeChange(game.id, 'min', e.target.value)}
                            />
                          </div>
                          <div className="input-group">
                            <label>Max:</label>
                            <input
                              type="number"
                              min="1"
                              value={range.max}
                              onChange={(e) => handleRangeChange(game.id, 'max', e.target.value)}
                            />
                          </div>
                        </div>
                      </div>

                      {/* Granular Toggle */}
                      <div className="granular-toggle">
                        <label className="checkbox-label">
                          <input
                            type="checkbox"
                            checked={range.useGranular || false}
                            onChange={() => handleToggleGranular(game.id)}
                          />
                          Use different ranges for each operand
                        </label>
                      </div>

                      {/* Operand-specific ranges */}
                      {range.useGranular && (
                        <div className="operand-ranges">
                          <div className="operand-section">
                            <h4>{game.operand1Label}</h4>
                            <div className="input-row">
                              <div className="input-group">
                                <label>Min:</label>
                                <input
                                  type="number"
                                  min="0"
                                  value={range.operand1Min ?? ''}
                                  onChange={(e) => handleRangeChange(game.id, 'operand1Min', e.target.value)}
                                  placeholder={range.min}
                                />
                              </div>
                              <div className="input-group">
                                <label>Max:</label>
                                <input
                                  type="number"
                                  min="1"
                                  value={range.operand1Max ?? ''}
                                  onChange={(e) => handleRangeChange(game.id, 'operand1Max', e.target.value)}
                                  placeholder={range.max}
                                />
                              </div>
                            </div>
                          </div>

                          <div className="operand-section">
                            <h4>{game.operand2Label}</h4>
                            <div className="input-row">
                              <div className="input-group">
                                <label>Min:</label>
                                <input
                                  type="number"
                                  min="0"
                                  value={range.operand2Min ?? ''}
                                  onChange={(e) => handleRangeChange(game.id, 'operand2Min', e.target.value)}
                                  placeholder={range.min}
                                />
                              </div>
                              <div className="input-group">
                                <label>Max:</label>
                                <input
                                  type="number"
                                  min="1"
                                  value={range.operand2Max ?? ''}
                                  onChange={(e) => handleRangeChange(game.id, 'operand2Max', e.target.value)}
                                  placeholder={range.max}
                                />
                              </div>
                            </div>
                          </div>
                        </div>
                      )}

                      <div className="range-summary">
                        {range.useGranular ? (
                          <>
                            {game.operand1Label.split(' ')[0]}: {range.operand1Min ?? range.min} to {range.operand1Max ?? range.max} |
                            {game.operand2Label.split(' ')[0]}: {range.operand2Min ?? range.min} to {range.operand2Max ?? range.max}
                          </>
                        ) : (
                          <>Numbers: {range.min} to {range.max}</>
                        )}
                      </div>

                      {/* Star Allocation Section */}
                      <div className="star-allocation-section">
                        <h4>Star Rewards per Correct Answer</h4>
                        <div className="star-inputs">
                          <div className="star-row">
                            <label>
                              <span className="difficulty-badge easy">Easy</span>
                            </label>
                            <div className="star-input-group">
                              <input
                                type="number"
                                min="1"
                                value={starRewards[game.id]?.easy || 5}
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
                                value={starRewards[game.id]?.medium || 10}
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
                                value={starRewards[game.id]?.hard || 15}
                                onChange={(e) => handleStarChange(game.id, 'hard', e.target.value)}
                              />
                              <span className="star-icon">⭐</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="range-disabled">
                      Using default game ranges and star rewards (Easy: 5⭐, Medium: 10⭐, Hard: 15⭐)
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </section>

        <section className="actions-section">
          <button className="secondary-button" onClick={handleResetToDefaults} disabled={saving}>
            Reset to Defaults
          </button>
          <button className="primary-button" onClick={handleSave} disabled={saving}>
            {saving ? 'Saving...' : 'Save Settings'}
          </button>
        </section>
      </div>
    </div>
  );
}

export default NumberRangeSettingsScreen;
