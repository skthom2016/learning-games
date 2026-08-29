-- Reset Database - Keep only David and Daniel
-- This script deletes all players except David and Daniel
-- CASCADE delete will remove all related data (attempts, mastery, rewards, etc.)

-- First, ensure David and Daniel exist (create them if they don't)
-- Using proper UUID format (8-4-4-4-12 hexadecimal characters)
INSERT INTO players (player_id, player_name, created_at, last_played_at, total_session_count)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'David', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0),
    ('00000000-0000-0000-0000-000000000002', 'Daniel', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
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
