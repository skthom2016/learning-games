-- Reset Database - Keep only David and Daniel
-- This script deletes all players except David and Daniel
-- CASCADE delete will remove all related data (attempts, mastery, rewards, etc.)

-- First, ensure David and Daniel exist (create them if they don't)
INSERT INTO players (player_id, player_name, created_at, last_played_at, total_session_count)
VALUES
    ('david-uuid-1234', 'David', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
    ('daniel-uuid-5678', 'Daniel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
ON CONFLICT (player_id) DO NOTHING;

-- Delete all other players and their associated data
-- CASCADE will automatically delete:
-- - player_game_profiles
-- - mastery_records
-- - attempt_records
-- - reward_transactions
-- - reward_reset_events
DELETE FROM players
WHERE player_name NOT IN ('David', 'Daniel');

-- Verify the results
SELECT player_id, player_name, created_at, total_session_count
FROM players
ORDER BY player_name;
