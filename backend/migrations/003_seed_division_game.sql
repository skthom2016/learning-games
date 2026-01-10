-- Seed data for Division Facts game
-- Division game structured similar to multiplication but with division problems

-- Insert division game
INSERT INTO games (game_id, display_name, description, icon_name, version, is_active)
VALUES (
    'division-facts',
    'Division Detective',
    'Solve division puzzles and become a math detective!',
    'division-icon.svg',
    '1.0.0',
    true
);

-- Insert difficulty levels for division game
INSERT INTO difficulty_levels (level_id, game_id, display_name, numeric_tier, description)
VALUES
    ('easy', 'division-facts', 'Easy', 1, 'Simple division (quotients 1-5), with visual hints always shown'),
    ('medium', 'division-facts', 'Medium', 2, 'All division facts (quotients 1-12), visual hints on request'),
    ('hard', 'division-facts', 'Hard', 3, 'All division facts, randomized, no hints');

-- Insert topics (division by 2-12)
INSERT INTO topics (topic_id, game_id, display_name, display_order, prerequisite_topic_ids, estimated_pool_size)
VALUES
    ('divide-2', 'division-facts', 'Divide by 2', 1, '{}', 10),
    ('divide-3', 'division-facts', 'Divide by 3', 2, '{}', 10),
    ('divide-4', 'division-facts', 'Divide by 4', 3, ARRAY['divide-2'], 10),
    ('divide-5', 'division-facts', 'Divide by 5', 4, '{}', 10),
    ('divide-6', 'division-facts', 'Divide by 6', 5, ARRAY['divide-2', 'divide-3'], 10),
    ('divide-7', 'division-facts', 'Divide by 7', 6, '{}', 10),
    ('divide-8', 'division-facts', 'Divide by 8', 7, ARRAY['divide-2', 'divide-4'], 10),
    ('divide-9', 'division-facts', 'Divide by 9', 8, ARRAY['divide-3'], 10),
    ('divide-10', 'division-facts', 'Divide by 10', 9, '{}', 10),
    ('divide-11', 'division-facts', 'Divide by 11', 10, '{}', 10),
    ('divide-12', 'division-facts', 'Divide by 12', 11, ARRAY['divide-2', 'divide-3', 'divide-4'], 10);

-- Note: This gives us 11 topics total (divide-2 through divide-12)
-- Each topic has approximately 10 unique questions (e.g., 2÷2 through 20÷2 for divide-by-2)
-- Division problems are generated as: (quotient × divisor) ÷ divisor = quotient
-- Easy: quotients 1-5, Medium/Hard: quotients 1-12
