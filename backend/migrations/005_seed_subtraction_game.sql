-- Seed data for Subtraction Facts game
-- Subtraction game for ages 6-9

-- Insert subtraction game
INSERT INTO games (game_id, display_name, description, icon_name, version, is_active)
VALUES (
    'subtraction-facts',
    'Subtraction Safari',
    'Explore subtraction and become a math ranger!',
    'subtraction-icon.svg',
    '1.0.0',
    true
);

-- Insert difficulty levels for subtraction game
INSERT INTO difficulty_levels (level_id, game_id, display_name, numeric_tier, description)
VALUES
    ('easy', 'subtraction-facts', 'Easy', 1, 'Simple subtraction (minuends up to 10), with visual hints always shown'),
    ('medium', 'subtraction-facts', 'Medium', 2, 'Subtraction with minuends up to 20, visual hints on request'),
    ('hard', 'subtraction-facts', 'Hard', 3, 'Subtraction with minuends up to 50, randomized, no hints');

-- Insert topics (subtracting 0-10)
INSERT INTO topics (topic_id, game_id, display_name, display_order, prerequisite_topic_ids, estimated_pool_size)
VALUES
    ('subtract-0', 'subtraction-facts', 'Subtracting Zero', 1, '{}', 10),
    ('subtract-1', 'subtraction-facts', 'Subtracting One', 2, '{}', 10),
    ('subtract-2', 'subtraction-facts', 'Subtracting Two', 3, '{}', 10),
    ('subtract-3', 'subtraction-facts', 'Subtracting Three', 4, ARRAY['subtract-1', 'subtract-2'], 10),
    ('subtract-4', 'subtraction-facts', 'Subtracting Four', 5, ARRAY['subtract-2'], 10),
    ('subtract-5', 'subtraction-facts', 'Subtracting Five', 6, '{}', 10),
    ('subtract-6', 'subtraction-facts', 'Subtracting Six', 7, ARRAY['subtract-3'], 10),
    ('subtract-7', 'subtraction-facts', 'Subtracting Seven', 8, '{}', 10),
    ('subtract-8', 'subtraction-facts', 'Subtracting Eight', 9, ARRAY['subtract-4'], 10),
    ('subtract-9', 'subtraction-facts', 'Subtracting Nine', 10, ARRAY['subtract-3'], 10),
    ('subtract-10', 'subtraction-facts', 'Subtracting Ten', 11, '{}', 10);

-- Note: This gives us 11 topics total (subtract-0 through subtract-10)
-- Each topic has approximately 10 unique questions based on difficulty
-- Easy: minuends up to 10, Medium: minuends up to 20, Hard: minuends up to 50
-- Subtraction problems are generated as: minuend - subtrahend = difference
