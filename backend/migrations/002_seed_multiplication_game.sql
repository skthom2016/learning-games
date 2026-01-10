-- Seed data for Multiplication Tables game
-- See ARCHITECTURE.md Section 3 for game definition

-- Insert multiplication game
INSERT INTO games (game_id, display_name, description, icon_name, version, is_active)
VALUES (
    'multiplication-tables',
    'Multiplication Master',
    'Practice times tables with fun characters!',
    'multiplication-icon.svg',
    '1.0.0',
    true
);

-- Insert difficulty levels for multiplication game
INSERT INTO difficulty_levels (level_id, game_id, display_name, numeric_tier, description)
VALUES
    ('easy', 'multiplication-tables', 'Easy', 1, 'Small numbers (2-5), with visual hints always shown'),
    ('medium', 'multiplication-tables', 'Medium', 2, 'All tables (2-12), visual hints on request'),
    ('hard', 'multiplication-tables', 'Hard', 3, 'All tables, randomized, no hints');

-- Insert topics (times tables 2-12)
INSERT INTO topics (topic_id, game_id, display_name, display_order, prerequisite_topic_ids, estimated_pool_size)
VALUES
    ('times-2', 'multiplication-tables', '2× table', 1, '{}', 12),
    ('times-3', 'multiplication-tables', '3× table', 2, '{}', 12),
    ('times-4', 'multiplication-tables', '4× table', 3, ARRAY['times-2'], 12),
    ('times-5', 'multiplication-tables', '5× table', 4, '{}', 12),
    ('times-6', 'multiplication-tables', '6× table', 5, ARRAY['times-2', 'times-3'], 12),
    ('times-7', 'multiplication-tables', '7× table', 6, '{}', 12),
    ('times-8', 'multiplication-tables', '8× table', 7, ARRAY['times-2', 'times-4'], 12),
    ('times-9', 'multiplication-tables', '9× table', 8, ARRAY['times-3'], 12),
    ('times-10', 'multiplication-tables', '10× table', 9, '{}', 12),
    ('times-11', 'multiplication-tables', '11× table', 10, '{}', 12),
    ('times-12', 'multiplication-tables', '12× table', 11, ARRAY['times-2', 'times-3', 'times-4'], 12);

-- Note: This gives us 11 topics total (times-2 through times-12)
-- Each topic has approximately 12 unique questions (e.g., 2×1 through 2×12)
