-- Fix foreign key constraint on reward_transactions.attempt_id
-- Add ON DELETE CASCADE to allow player deletion to work properly

-- First, drop the existing foreign key constraint
DO $$ DECLARE
    fk_name TEXT;
BEGIN
    -- Get the name of the foreign key constraint
    SELECT conname INTO fk_name
    FROM pg_constraint
    WHERE conrelid = 'reward_transactions'::regclass
      AND confrelid = 'attempt_records'::regclass
      AND contype = 'f';

    -- Drop the constraint if it exists
    IF fk_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE reward_transactions DROP CONSTRAINT %I', fk_name);
    END IF;
END $$;

-- Add the foreign key with ON DELETE CASCADE
ALTER TABLE reward_transactions
ADD CONSTRAINT reward_transactions_attempt_id_fkey
FOREIGN KEY (attempt_id) REFERENCES attempt_records(attempt_id) ON DELETE SET NULL;

-- Using ON DELETE SET NULL instead of CASCADE to preserve reward transaction history
-- When an attempt is deleted, the reward transaction remains but attempt_id becomes NULL
