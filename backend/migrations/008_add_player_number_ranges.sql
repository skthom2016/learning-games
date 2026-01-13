-- Player Number Ranges table
-- Stores custom min/max number ranges per player per game
-- This allows grade-based customization (e.g., 1st grade uses numbers 1-10, 5th grade uses 1-1000)

CREATE TABLE player_number_ranges (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    game_id VARCHAR(50) NOT NULL REFERENCES games(game_id),
    min_number INTEGER NOT NULL CHECK (min_number >= 0),
    max_number INTEGER NOT NULL CHECK (max_number > min_number),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(player_id, game_id)
);

-- Create index for lookups
CREATE INDEX idx_player_number_ranges_player ON player_number_ranges(player_id);
CREATE INDEX idx_player_number_ranges_game ON player_number_ranges(game_id);

-- Comment for documentation
COMMENT ON TABLE player_number_ranges IS 'Custom number range settings per player per game for grade-based difficulty customization';
