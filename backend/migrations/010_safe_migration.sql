-- SAFE MIGRATION: Fix reward_transactions foreign key constraint
-- This migration uses a safe approach: backup -> modify -> verify

-- Step 1: Create backup table
CREATE TABLE IF NOT EXISTS reward_transactions_backup_20250112 AS
SELECT * FROM reward_transactions;

-- Verify backup was created
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'reward_transactions_backup_20250112') THEN
        RAISE EXCEPTION 'Backup table was not created!';
    END IF;

    IF (SELECT COUNT(*) FROM reward_transactions_backup_20250112) != (SELECT COUNT(*) FROM reward_transactions) THEN
        RAISE EXCEPTION 'Backup count mismatch! Backup has % rows, original has % rows',
            (SELECT COUNT(*) FROM reward_transactions_backup_20250112),
            (SELECT COUNT(*) FROM reward_transactions);
    END IF;

    RAISE NOTICE 'Backup created successfully with % records', (SELECT COUNT(*) FROM reward_transactions_backup_20250112);
END $$;

-- Step 2: Drop the old foreign key constraint (if exists)
DO $$
DECLARE
    fk_name TEXT;
BEGIN
    SELECT conname INTO fk_name
    FROM pg_constraint
    WHERE conrelid = 'reward_transactions'::regclass
      AND confrelid = 'attempt_records'::regclass
      AND contype = 'f';

    IF fk_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE reward_transactions DROP CONSTRAINT %I', fk_name);
        RAISE NOTICE 'Dropped old constraint: %', fk_name;
    ELSE
        RAISE NOTICE 'No old constraint found - this is okay if the table is new';
    END IF;
END $$;

-- Step 3: Add the new foreign key with ON DELETE SET NULL
DO $$
BEGIN
    -- Check if constraint already exists
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'reward_transactions'::regclass
        AND conname = 'reward_transactions_attempt_id_fkey'
    ) THEN
        ALTER TABLE reward_transactions
        ADD CONSTRAINT reward_transactions_attempt_id_fkey
        FOREIGN KEY (attempt_id) REFERENCES attempt_records(attempt_id) ON DELETE SET NULL;

        RAISE NOTICE 'New foreign key constraint added successfully';
    ELSE
        RAISE NOTICE 'Constraint already exists - skipping';
    END IF;
END $$;

-- Step 4: Verification - count the original records
DO $$
DECLARE
    backup_count INT;
    original_count INT;
BEGIN
    SELECT COUNT(*) INTO backup_count FROM reward_transactions_backup_20250112;
    SELECT COUNT(*) INTO original_count FROM reward_transactions;

    RAISE NOTICE '=== MIGRATION VERIFICATION ===';
    RAISE NOTICE 'Backup records: %', backup_count;
    RAISE NOTICE 'Current records: %', original_count;

    IF backup_count != original_count THEN
        RAISE WARNING 'Record count changed! Was: %, Now: %', backup_count, original_count;
    ELSE
        RAISE NOTICE '✓ Migration completed successfully - data integrity verified';
    END IF;
END $$;

-- NOTE: The backup table reward_transactions_backup_20250112 is kept for safety
-- It can be dropped after verifying everything works:
-- DROP TABLE reward_transactions_backup_20250112;
