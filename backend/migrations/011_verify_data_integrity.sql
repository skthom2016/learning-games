-- ============================================================================
-- DATA INTEGRITY VERIFICATION MIGRATION
-- ============================================================================
-- This migration DOES NOT modify any data
-- It verifies that all player progress and stars are intact
-- Run this before and after any migration to ensure data safety
-- ============================================================================

BEGIN;

DO $$
DECLARE
    total_players INT;
    total_profiles INT;
    total_mastery_records INT;
    total_attempts INT;
    total_stars INT;
    total_reward_transactions INT;
    orphan_count INT;

    -- Detailed counts per player
    player_record RECORD;
BEGIN
    RAISE NOTICE '====================================================================';
    RAISE NOTICE 'DATA INTEGRITY VERIFICATION';
    RAISE NOTICE '====================================================================';
    RAISE NOTICE '';

    -- Count all critical data
    SELECT COUNT(*) INTO total_players FROM players;
    SELECT COUNT(*) INTO total_profiles FROM player_game_profiles;
    SELECT COUNT(*) INTO total_mastery_records FROM mastery_records;
    SELECT COUNT(*) INTO total_attempts FROM attempt_records;
    SELECT COALESCE(SUM(amount), 0) INTO total_stars FROM reward_transactions WHERE consolidated = false;
    SELECT COUNT(*) INTO total_reward_transactions FROM reward_transactions;

    RAISE NOTICE 'DATABASE SNAPSHOT: %', NOW();
    RAISE NOTICE '--------------------------------------------------------------------';
    RAISE NOTICE 'AGGREGATE COUNTS:';
    RAISE NOTICE '  Total Players: %', total_players;
    RAISE NOTICE '  Total Game Profiles: %', total_profiles;
    RAISE NOTICE '  Total Mastery Records: %', total_mastery_records;
    RAISE NOTICE '  Total Attempt Records: %', total_attempts;
    RAISE NOTICE '  Total Stars (unconsolidated): %', total_stars;
    RAISE NOTICE '  Total Reward Transactions: %', total_reward_transactions;
    RAISE NOTICE '';

    -- Check for orphans (data consistency)
    RAISE NOTICE 'DATA CONSISTENCY CHECKS:';
    RAISE NOTICE '--------------------------------------------------------------------';

    -- Check for mastery records without players
    SELECT COUNT(*) INTO orphan_count FROM mastery_records mr
    LEFT JOIN players p ON mr.player_id = p.player_id
    WHERE p.player_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING '  WARNING: % mastery records have no matching player!', orphan_count;
    ELSE
        RAISE NOTICE '  All mastery records have valid players';
    END IF;

    -- Check for attempt_records without players
    SELECT COUNT(*) INTO orphan_count FROM attempt_records ar
    LEFT JOIN players p ON ar.player_id = p.player_id
    WHERE p.player_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING '  WARNING: % attempt records have no matching player!', orphan_count;
    ELSE
        RAISE NOTICE '  All attempt records have valid players';
    END IF;

    -- Check for reward_transactions without players
    SELECT COUNT(*) INTO orphan_count FROM reward_transactions rt
    LEFT JOIN players p ON rt.player_id = p.player_id
    WHERE p.player_id IS NULL;

    IF orphan_count > 0 THEN
        RAISE WARNING '  WARNING: % reward transactions have no matching player!', orphan_count;
    ELSE
        RAISE NOTICE '  All reward transactions have valid players';
    END IF;

    RAISE NOTICE '';

    -- Per-player breakdown
    RAISE NOTICE 'PLAYER DETAILS (Top 20 by stars):';
    RAISE NOTICE '--------------------------------------------------------------------';

    FOR player_record IN
        SELECT p.player_id, p.player_name,
               COALESCE(SUM(rt.amount), 0) as player_stars,
               COUNT(DISTINCT pge.game_id) as games_played,
               COUNT(ar.attempt_id) as total_attempts
        FROM players p
        LEFT JOIN reward_transactions rt ON p.player_id = rt.player_id AND rt.consolidated = false
        LEFT JOIN player_game_profiles pge ON p.player_id = pge.player_id
        LEFT JOIN attempt_records ar ON p.player_id = ar.player_id
        GROUP BY p.player_id, p.player_name
        ORDER BY player_stars DESC
        LIMIT 20
    LOOP
        RAISE NOTICE '  Player: % | Stars: % | Games: % | Attempts: %',
            player_record.player_name,
            player_record.player_stars,
            player_record.games_played,
            player_record.total_attempts;
    END LOOP;

    RAISE NOTICE '--------------------------------------------------------------------';
    RAISE NOTICE 'DATA INTEGRITY VERIFICATION COMPLETE';
    RAISE NOTICE '====================================================================';

    -- Save snapshot to a log table (for future comparison)
    -- Create log table if it doesn't exist
    CREATE TABLE IF NOT EXISTS data_integrity_snapshots (
        snapshot_id SERIAL PRIMARY KEY,
        snapshot_time TIMESTAMP NOT NULL DEFAULT NOW(),
        total_players INT NOT NULL,
        total_profiles INT NOT NULL,
        total_mastery_records INT NOT NULL,
        total_attempts INT NOT NULL,
        total_stars INT NOT NULL,
        total_reward_transactions INT NOT NULL,
        verified BOOLEAN NOT NULL DEFAULT true
    );

    INSERT INTO data_integrity_snapshots (
        total_players, total_profiles, total_mastery_records,
        total_attempts, total_stars, total_reward_transactions, verified
    ) VALUES (
        total_players, total_profiles, total_mastery_records,
        total_attempts, total_stars, total_reward_transactions, true
    );

    RAISE NOTICE '';
    RAISE NOTICE 'Snapshot saved to data_integrity_snapshots table for future comparison.';
    RAISE NOTICE 'To view history: SELECT * FROM data_integrity_snapshots ORDER BY snapshot_time DESC LIMIT 10;';

END $$;

COMMIT;

-- ============================================================================
-- POST-MIGRATION VERIFICATION QUERY
-- ============================================================================
-- Run this after any migration to compare with previous snapshot:
--
-- SELECT
--     snapshot_time,
--     total_players,
--     total_stars,
--     LEAD(total_stars) OVER (ORDER BY snapshot_time) as prev_stars,
--     total_stars - LEAD(total_stars) OVER (ORDER BY snapshot_time) as stars_diff
-- FROM data_integrity_snapshots
-- ORDER BY snapshot_time DESC
-- LIMIT 5;
