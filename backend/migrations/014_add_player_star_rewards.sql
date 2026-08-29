-- Add player_star_rewards table for custom star rewards per difficulty level per player
-- Admins can configure different star amounts for easy/medium/hard based on player's age/ability

CREATE TABLE player_star_rewards (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    game_id VARCHAR(50) NOT NULL REFERENCES games(game_id),
    easy_stars INTEGER NOT NULL DEFAULT 5 CHECK (easy_stars > 0),
    medium_stars INTEGER NOT NULL DEFAULT 10 CHECK (medium_stars > 0),
    hard_stars INTEGER NOT NULL DEFAULT 15 CHECK (hard_stars > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(player_id, game_id)
);

-- Index for looking up star rewards by player
CREATE INDEX idx_star_rewards_player ON player_star_rewards(player_id);

-- Comment for documentation
COMMENT ON TABLE player_star_rewards IS 'Custom star rewards per difficulty level per player per game';
COMMENT ON COLUMN player_star_rewards.easy_stars IS 'Stars awarded for correct answers at easy difficulty';
COMMENT ON COLUMN player_star_rewards.medium_stars IS 'Stars awarded for correct answers at medium difficulty';
COMMENT ON COLUMN player_star_rewards.hard_stars IS 'Stars awarded for correct answers at hard difficulty';
