-- Add topic_difficulty_progress table for tracking difficulty progression per topic
-- Each topic has 3 difficulty levels: easy, medium, hard
-- Players progress based on: 10 consecutive correct = level up, 2 consecutive wrong = level down

CREATE TABLE topic_difficulty_progress (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    game_id VARCHAR(50) NOT NULL REFERENCES games(game_id),
    topic_id VARCHAR(50) NOT NULL,
    current_difficulty_level_id VARCHAR(50) NOT NULL DEFAULT 'easy',
    consecutive_correct INTEGER NOT NULL DEFAULT 0,
    consecutive_wrong INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(player_id, game_id, topic_id),
    FOREIGN KEY (game_id, topic_id) REFERENCES topics(game_id, topic_id)
);

-- Index for looking up progress by player
CREATE INDEX idx_topic_progress_player ON topic_difficulty_progress(player_id);

-- Comment for documentation
COMMENT ON TABLE topic_difficulty_progress IS 'Tracks difficulty progression per topic per player (easy/medium/hard)';
