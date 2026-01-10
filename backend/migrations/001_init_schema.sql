-- Initial database schema for Learning Game Platform
-- See ARCHITECTURE.md Section 7 for complete data model specification

-- Create ENUM types
CREATE TYPE mastery_state AS ENUM ('UNKNOWN', 'WEAK', 'LEARNING', 'STRONG', 'MASTERED');
CREATE TYPE reward_type AS ENUM ('STAR', 'STICKER', 'BADGE', 'CHARACTER_UNLOCK');

-- Players table
CREATE TABLE players (
    player_id UUID PRIMARY KEY,
    player_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_played_at TIMESTAMP,
    total_session_count INTEGER NOT NULL DEFAULT 0,
    preferred_character_id VARCHAR(50)
);

-- Games table
CREATE TABLE games (
    game_id VARCHAR(50) PRIMARY KEY,
    display_name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    icon_name VARCHAR(100) NOT NULL,
    version VARCHAR(20) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Topics table
CREATE TABLE topics (
    topic_id VARCHAR(50) NOT NULL,
    game_id VARCHAR(50) NOT NULL REFERENCES games(game_id),
    display_name VARCHAR(100) NOT NULL,
    display_order INTEGER NOT NULL,
    prerequisite_topic_ids TEXT[], -- Array of topic IDs
    estimated_pool_size INTEGER NOT NULL,
    PRIMARY KEY (game_id, topic_id)
);

-- Difficulty Levels table
CREATE TABLE difficulty_levels (
    level_id VARCHAR(50) NOT NULL,
    game_id VARCHAR(50) NOT NULL REFERENCES games(game_id),
    display_name VARCHAR(100) NOT NULL,
    numeric_tier INTEGER NOT NULL CHECK (numeric_tier >= 1 AND numeric_tier <= 5),
    description VARCHAR(500),
    PRIMARY KEY (game_id, level_id)
);

-- Player Game Profiles table
CREATE TABLE player_game_profiles (
    profile_id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    game_id VARCHAR(50) NOT NULL REFERENCES games(game_id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_played_at TIMESTAMP,
    total_questions_answered INTEGER NOT NULL DEFAULT 0,
    total_correct_answers INTEGER NOT NULL DEFAULT 0,
    current_difficulty_level_id VARCHAR(50) NOT NULL,
    confidence_recovery_mode_active BOOLEAN NOT NULL DEFAULT false,
    confidence_recovery_entered_at TIMESTAMP,
    session_count INTEGER NOT NULL DEFAULT 0,
    UNIQUE(player_id, game_id)
);

-- Mastery Records table
CREATE TABLE mastery_records (
    mastery_record_id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    game_id VARCHAR(50) NOT NULL,
    topic_id VARCHAR(50) NOT NULL,
    mastery_state mastery_state NOT NULL DEFAULT 'UNKNOWN',
    last_assessed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    times_practiced INTEGER NOT NULL DEFAULT 0,
    consecutive_correct INTEGER NOT NULL DEFAULT 0,
    last_correct_at TIMESTAMP,
    last_incorrect_at TIMESTAMP,
    UNIQUE(player_id, game_id, topic_id),
    FOREIGN KEY (game_id, topic_id) REFERENCES topics(game_id, topic_id)
);

-- Attempt Records table (immutable log)
CREATE TABLE attempt_records (
    attempt_id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    game_id VARCHAR(50) NOT NULL,
    topic_id VARCHAR(50) NOT NULL,
    question_id UUID NOT NULL,
    difficulty_level_id VARCHAR(50) NOT NULL,
    difficulty_numeric_tier INTEGER NOT NULL CHECK (difficulty_numeric_tier >= 1 AND difficulty_numeric_tier <= 5),
    was_correct BOOLEAN NOT NULL,
    submitted_answer VARCHAR(200) NOT NULL,
    correct_answer VARCHAR(200) NOT NULL,
    time_to_answer_ms INTEGER,
    hint_used BOOLEAN NOT NULL DEFAULT false,
    confidence_recovery_mode BOOLEAN NOT NULL DEFAULT false,
    attempted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id, topic_id) REFERENCES topics(game_id, topic_id)
);

-- Reward Transactions table (immutable ledger)
CREATE TABLE reward_transactions (
    transaction_id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    reward_type reward_type NOT NULL,
    reward_identifier VARCHAR(100) NOT NULL,
    amount INTEGER NOT NULL,
    earned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    game_id VARCHAR(50),
    topic_id VARCHAR(50),
    attempt_id UUID REFERENCES attempt_records(attempt_id),
    difficulty_tier INTEGER CHECK (difficulty_tier >= 1 AND difficulty_tier <= 5),
    confidence_recovery_mode BOOLEAN NOT NULL DEFAULT false,
    consolidated BOOLEAN NOT NULL DEFAULT false,
    notes VARCHAR(500)
);

-- Reward Reset Events table (audit trail)
CREATE TABLE reward_reset_events (
    reset_id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(player_id) ON DELETE CASCADE,
    reset_by VARCHAR(100) NOT NULL,
    reset_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    stars_consolidated INTEGER NOT NULL CHECK (stars_consolidated > 0),
    real_world_reward VARCHAR(500) NOT NULL,
    transaction_ids_consolidated UUID[] -- Array of transaction IDs
);

-- Create indexes for performance
-- See ARCHITECTURE.md Section 7 for index rationale

-- Attempt records indexes (for recent accuracy calculations)
CREATE INDEX idx_attempts_player_game_topic_time ON attempt_records(player_id, game_id, topic_id, attempted_at DESC);
CREATE INDEX idx_attempts_player_time ON attempt_records(player_id, attempted_at DESC);
CREATE INDEX idx_attempts_game_topic ON attempt_records(game_id, topic_id);

-- Reward transactions indexes (for balance calculations)
CREATE INDEX idx_rewards_player_time ON reward_transactions(player_id, earned_at DESC);
CREATE INDEX idx_rewards_player_type_consolidated ON reward_transactions(player_id, reward_type, consolidated);

-- Mastery records index (for quick lookups)
CREATE INDEX idx_mastery_player_game ON mastery_records(player_id, game_id);

-- Player game profiles index
CREATE INDEX idx_pgp_player ON player_game_profiles(player_id);

-- Comments for documentation
COMMENT ON TABLE players IS 'Child users of the platform';
COMMENT ON TABLE games IS 'Learning games (e.g., multiplication, division)';
COMMENT ON TABLE topics IS 'Learnable sub-units within games';
COMMENT ON TABLE mastery_records IS 'Tracks learning progress per topic per player';
COMMENT ON TABLE attempt_records IS 'Immutable log of every question attempt';
COMMENT ON TABLE reward_transactions IS 'Immutable ledger of all rewards earned';
COMMENT ON TABLE reward_reset_events IS 'Audit trail for admin reward consolidation';
