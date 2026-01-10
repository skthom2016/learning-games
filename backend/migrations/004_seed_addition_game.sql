-- Seed data for Addition Facts game
-- Addition game for ages 6-9

-- Insert addition game
INSERT INTO games (game_id, display_name, description, icon_name, version, is_active)
VALUES (
    'addition-facts',
    'Addition Adventure',
    'Master addition and become a math explorer!',
    'addition-icon.svg',
    '1.0.0',
    true
);

-- Insert difficulty levels for addition game
INSERT INTO difficulty_levels (level_id, game_id, display_name, numeric_tier, description)
VALUES
    ('easy', 'addition-facts', 'Easy', 1, 'Simple addition (sums up to 10), with visual hints always shown'),
    ('medium', 'addition-facts', 'Medium', 2, 'Addition sums up to 20, visual hints on request'),
    ('hard', 'addition-facts', 'Hard', 3, 'Addition sums up to 50, randomized, no hints');

-- Insert topics (adding 0-10)
INSERT INTO topics (topic_id, game_id, display_name, display_order, prerequisite_topic_ids, estimated_pool_size)
VALUES
    ('add-0', 'addition-facts', 'Adding Zero', 1, '{}', 10),
    ('add-1', 'addition-facts', 'Adding One', 2, '{}', 10),
    ('add-2', 'addition-facts', 'Adding Two', 3, '{}', 10),
    ('add-3', 'addition-facts', 'Adding Three', 4, ARRAY['add-1', 'add-2'], 10),
    ('add-4', 'addition-facts', 'Adding Four', 5, ARRAY['add-2'], 10),
    ('add-5', 'addition-facts', 'Adding Five', 6, '{}', 10),
    ('add-6', 'addition-facts', 'Adding Six', 7, ARRAY['add-3'], 10),
    ('add-7', 'addition-facts', 'Adding Seven', 8, '{}', 10),
    ('add-8', 'addition-facts', 'Adding Eight', 9, ARRAY['add-4'], 10),
    ('add-9', 'addition-facts', 'Adding Nine', 10, ARRAY['add-3'], 10),
    ('add-10', 'addition-facts', 'Adding Ten', 11, '{}', 10);

-- Note: This gives us 11 topics total (add-0 through add-10)
-- Each topic has approximately 10 unique questions based on difficulty
-- Easy: sums up to 10, Medium: sums up to 20, Hard: sums up to 50
-- Addition problems are generated as: addend1 + addend2 = sum
