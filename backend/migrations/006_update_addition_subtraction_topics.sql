-- Update addition and subtraction games with skill-based level progression
-- This replaces the old add-0/add-1 and subtract-0/subtract-1 topics with level-based topics

-- First, delete old topics for addition
DELETE FROM topics WHERE game_id = 'addition-facts';

-- Insert new level-based topics for addition
INSERT INTO topics (topic_id, game_id, display_name, display_order, prerequisite_topic_ids, estimated_pool_size)
VALUES
    ('add-level-1', 'addition-facts', 'Level 1: Single Digit (Sums to 10)', 1, '{}', 50),
    ('add-level-2', 'addition-facts', 'Level 2: Single Digit (Sums to 18)', 2, ARRAY['add-level-1'], 50),
    ('add-level-3', 'addition-facts', 'Level 3: Double Digit (No Carry)', 3, ARRAY['add-level-2'], 100),
    ('add-level-4', 'addition-facts', 'Level 4: Double Digit (With Carry)', 4, ARRAY['add-level-3'], 200),
    ('add-level-5', 'addition-facts', 'Level 5: Triple Digit (No Carry)', 5, ARRAY['add-level-4'], 200),
    ('add-level-6', 'addition-facts', 'Level 6: Triple Digit (With Carry)', 6, ARRAY['add-level-5'], 500),
    ('add-level-7', 'addition-facts', 'Level 7: Triple Digit (Complex)', 7, ARRAY['add-level-6'], 500);

-- Delete old topics for subtraction
DELETE FROM topics WHERE game_id = 'subtraction-facts';

-- Insert new level-based topics for subtraction
INSERT INTO topics (topic_id, game_id, display_name, display_order, prerequisite_topic_ids, estimated_pool_size)
VALUES
    ('subtract-level-1', 'subtraction-facts', 'Level 1: Single Digit (Simple)', 1, '{}', 50),
    ('subtract-level-2', 'subtraction-facts', 'Level 2: Single Digit (Results 0-5)', 2, ARRAY['subtract-level-1'], 50),
    ('subtract-level-3', 'subtraction-facts', 'Level 3: Double Digit (No Borrow)', 3, ARRAY['subtract-level-2'], 100),
    ('subtract-level-4', 'subtraction-facts', 'Level 4: Double Digit (With Borrow)', 4, ARRAY['subtract-level-3'], 200),
    ('subtract-level-5', 'subtraction-facts', 'Level 5: Triple Digit (No Borrow)', 5, ARRAY['subtract-level-4'], 200),
    ('subtract-level-6', 'subtraction-facts', 'Level 6: Triple Digit (With Borrow)', 6, ARRAY['subtract-level-5'], 500),
    ('subtract-level-7', 'subtraction-facts', 'Level 7: Triple Digit (Complex)', 7, ARRAY['subtract-level-6'], 500);

-- Note: Progression is now performance-based through the adaptive learning system
-- Players must demonstrate mastery (through consecutive correct answers) to unlock higher levels
-- Each topic has a large pool size to allow for varied practice
