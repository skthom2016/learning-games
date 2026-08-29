# Database Migration Guide

This guide explains how to apply database migrations automatically when deploying to existing systems.

## How It Works

The system now includes an **automatic migration runner** that:
1. Checks which migrations have already been applied
2. Runs any pending migrations in order
3. Tracks applied migrations in a `schema_migrations` table
4. Runs automatically on every `docker-compose up`

## Key Features

✅ **Safe for existing databases** - Only runs new migrations, never re-runs old ones
✅ **Automatic on startup** - No manual intervention needed
✅ **Transactional** - Each migration runs in a transaction (rolls back on error)
✅ **Zero downtime** - Works with running databases

## Migration Files

Current migrations include:
- `001_initial_schema.sql` - Core tables (players, games, topics, etc.)
- `002_seed_games.sql` - Seed data for games
- `003-012_*.sql` - Previous features
- `013_add_topic_difficulty_progress.sql` - **NEW** Topic difficulty tracking
- `014_add_player_star_rewards.sql` - **NEW** Custom star rewards per player

## Deploying to Another Computer

### Step 1: Transfer the Code

Copy the entire codebase to the target computer, including:
```
division/
├── backend/
│   ├── migrations/          # All migration files
│   ├── internal/database/   # Migration runner
│   └── ...
├── frontend/
└── docker-compose.yml
```

### Step 2: Start the Services

```bash
# Navigate to the project directory
cd division

# Start all services
docker-compose up --build -d

# Watch the logs to see migrations being applied
docker-compose logs -f backend
```

**Expected output:**
```
learning-game-backend | Starting database migration check...
learning-game-backend | Applying migration: 013_add_topic_difficulty_progress.sql
learning-game-backend | Successfully applied migration: 013_add_topic_difficulty_progress.sql
learning-game-backend | Applying migration: 014_add_player_star_rewards.sql
learning-game-backend | Successfully applied migration: 014_add_player_star_rewards.sql
learning-game-backend | Successfully applied 2 migration(s)
learning-game-backend | Server starting on port 8080...
```

### Step 3: Verify

```bash
# Check that migrations were applied
docker exec -it learning-game-db psql -U gameuser -d learning_game -c "SELECT * FROM schema_migrations ORDER BY id;"

# Verify new tables exist
docker exec -it learning-game-db psql -U gameuser -d learning_game -c "\dt"
```

## For Existing Databases

If the target computer **already has the application running** with data:

1. **Backup first** (always!):
   ```bash
   docker exec learning-game-db pg_dump -U gameuser learning_game > backup_$(date +%Y%m%d).sql
   ```

2. **Pull the latest code**:
   ```bash
   git pull origin main
   # or copy the updated codebase
   ```

3. **Restart services**:
   ```bash
   docker-compose down
   docker-compose up --build -d
   ```

4. **Migrations run automatically** - The backend will:
   - Detect existing data (won't re-run old migrations)
   - Apply only the new migrations (013, 014)
   - Keep all existing player data intact

## Migration Safety

✅ **No data loss** - Migrations only ADD tables, never DELETE or MODIFY existing data
✅ **Atomic** - Each migration runs in a transaction
✅ **Idempotent** - Safe to restart if interrupted
✅ **Trackable** - All migrations are logged in `schema_migrations` table

## Troubleshooting

### Problem: Migration fails with error

**Solution:**
```bash
# Check the error logs
docker-compose logs backend

# Rollback: The failed migration is automatically rolled back
# Fix the issue in the migration file and restart
docker-compose restart backend
```

### Problem: Migration already applied manually

**Solution:**
```bash
# Mark it as applied without running it
docker exec -it learning-game-db psql -U gameuser -d learning_game
INSERT INTO schema_migrations (filename) VALUES ('013_add_topic_difficulty_progress.sql');
INSERT INTO schema_migrations (filename) VALUES ('014_add_player_star_rewards.sql');
\q
```

### Problem: Need to reset and re-apply all migrations

**Solution:**
```bash
# WARNING: This deletes ALL data!
docker-compose down -v
docker-compose up --build -d
```

### Problem: Want to see which migrations are applied

**Solution:**
```bash
docker exec -it learning-game-db psql -U gameuser -d learning_game -c \
  "SELECT id, filename, applied_at FROM schema_migrations ORDER BY id;"
```

## Adding New Migrations

When you add a new feature that requires database changes:

1. **Create migration file** in `backend/migrations/`:
   ```bash
   # Name it with the next number in sequence
   touch backend/migrations/015_add_new_feature.sql
   ```

2. **Write SQL**:
   ```sql
   -- 015_add_new_feature.sql
   CREATE TABLE new_feature (
       id UUID PRIMARY KEY,
       player_id UUID REFERENCES players(player_id) ON DELETE CASCADE,
       -- ... other columns
   );

   CREATE INDEX idx_new_feature_player ON new_feature(player_id);
   ```

3. **Deploy** - Just restart the backend:
   ```bash
   docker-compose restart backend
   # Migration runs automatically
   ```

## Production Deployment Checklist

- [ ] Backup database before deployment
- [ ] Test migrations on staging environment first
- [ ] Review migration files for correctness
- [ ] Ensure migrations are backward compatible (don't break old code)
- [ ] Monitor logs during deployment
- [ ] Verify new tables exist after deployment
- [ ] Test application functionality
- [ ] Keep backup for 30 days

## Technical Details

**Migration Tracking Table:**
```sql
CREATE TABLE schema_migrations (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Migration Runner:**
- Location: `backend/internal/database/migrations.go`
- Runs on: Backend startup (before API starts)
- Order: Alphabetical by filename (001, 002, 003, ...)
- Execution: Each migration in its own transaction

**Environment Variable:**
- `MIGRATIONS_PATH`: Path to migrations folder (default: `./migrations`)
- Can be customized in `docker-compose.yml` if needed

## Summary

The migration system is now **fully automatic**. When you deploy the updated codebase to another computer:

1. Copy the code
2. Run `docker-compose up --build -d`
3. Migrations apply automatically
4. No manual intervention needed
5. Existing data is preserved

✅ Safe, automatic, and production-ready!
