-- Migration: Add granular operand-specific number ranges
-- This allows setting different min/max for each operand type per game:
-- - Addition: addend1 min/max, addend2 min/max
-- - Subtraction: minuend min/max, subtrahend min/max
-- - Multiplication: factor1 min/max, factor2 min/max
-- - Division: dividend min/max, divisor min/max

-- Add new columns for operand-specific ranges
-- For backward compatibility, these are nullable and we fall back to min_number/max_number if not set

ALTER TABLE player_number_ranges
ADD COLUMN IF NOT EXISTS operand1_min INTEGER CHECK (operand1_min >= 0),
ADD COLUMN IF NOT EXISTS operand1_max INTEGER CHECK (operand1_max >= 1),
ADD COLUMN IF NOT EXISTS operand2_min INTEGER CHECK (operand2_min >= 0),
ADD COLUMN IF NOT EXISTS operand2_max INTEGER CHECK (operand2_max >= 1);

-- Add constraint to ensure operand1_min < operand1_max when both are set
ALTER TABLE player_number_ranges
ADD CONSTRAINT check_operand1_range
CHECK (operand1_min IS NULL OR operand1_max IS NULL OR operand1_max > operand1_min);

-- Add constraint to ensure operand2_min < operand2_max when both are set
ALTER TABLE player_number_ranges
ADD CONSTRAINT check_operand2_range
CHECK (operand2_min IS NULL OR operand2_max IS NULL OR operand2_max > operand2_min);

-- Comment for documentation
COMMENT ON COLUMN player_number_ranges.operand1_min IS 'Min value for first operand (addend1/minuend/factor1/dividend)';
COMMENT ON COLUMN player_number_ranges.operand1_max IS 'Max value for first operand (addend1/minuend/factor1/dividend)';
COMMENT ON COLUMN player_number_ranges.operand2_min IS 'Min value for second operand (addend2/subtrahend/factor2/divisor)';
COMMENT ON COLUMN player_number_ranges.operand2_max IS 'Max value for second operand (addend2/subtrahend/factor2/divisor)';
