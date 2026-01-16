# Deployment Quick Start

## Deploy to Another Computer (with existing database)

### 1. Backup First (Always!)
```bash
docker exec learning-game-db pg_dump -U gameuser learning_game > backup_$(date +%Y%m%d).sql
```

### 2. Copy Updated Code
Transfer the entire `division/` directory to the target computer.

### 3. Deploy
```bash
cd division
docker-compose down
docker-compose up --build -d
```

### 4. Watch Migrations Apply
```bash
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

### 5. Verify
```bash
# Check applied migrations
docker exec -it learning-game-db psql -U gameuser -d learning_game -c "SELECT * FROM schema_migrations ORDER BY id;"

# Check new tables exist
docker exec -it learning-game-db psql -U gameuser -d learning_game -c "\d topic_difficulty_progress"
docker exec -it learning-game-db psql -U gameuser -d learning_game -c "\d player_star_rewards"
```

## What Gets Applied?

### New Tables Added:

1. **`topic_difficulty_progress`** - Tracks difficulty progression per topic
   - Columns: `player_id`, `game_id`, `topic_id`, `current_difficulty_level_id`
   - Progressive difficulty: easy → medium → hard
   - 10 consecutive correct = level up
   - 2 consecutive wrong = level down

2. **`player_star_rewards`** - Custom star rewards per player per game
   - Columns: `player_id`, `game_id`, `easy_stars`, `medium_stars`, `hard_stars`
   - Admin can customize rewards based on player age/ability
   - Default: 5 stars (easy), 10 stars (medium), 15 stars (hard)

### Existing Data
✅ All existing player data is preserved
✅ All game history, rewards, and progress remains intact
✅ No data is deleted or modified

## Troubleshooting

### Problem: "Migration already exists" error
**Cause:** Migration was manually applied before
**Solution:** Mark it as applied:
```bash
docker exec -it learning-game-db psql -U gameuser -d learning_game
INSERT INTO schema_migrations (filename) VALUES ('013_add_topic_difficulty_progress.sql');
INSERT INTO schema_migrations (filename) VALUES ('014_add_player_star_rewards.sql');
\q
```

### Problem: "Table already exists" error
**Cause:** Tables were created manually
**Solution:** Same as above - just mark migrations as applied

### Problem: Backend won't start
**Check logs:**
```bash
docker-compose logs backend
```

**Common causes:**
- Database not ready (wait 10 seconds and retry)
- Migration syntax error (check logs for SQL error)
- Permission issue (check database credentials)

## Rollback (if needed)

### Undo migrations (keep data):
```bash
# Connect to database
docker exec -it learning-game-db psql -U gameuser -d learning_game

# Drop new tables
DROP TABLE IF EXISTS player_star_rewards;
DROP TABLE IF EXISTS topic_difficulty_progress;

# Remove from migration tracking
DELETE FROM schema_migrations WHERE filename IN (
    '013_add_topic_difficulty_progress.sql',
    '014_add_player_star_rewards.sql'
);
\q
```

### Full rollback (deletes all data):
```bash
docker-compose down -v
docker-compose up --build -d
```

## Key Files Changed

- ✅ `backend/internal/database/migrations.go` - Migration runner
- ✅ `backend/cmd/server/main.go` - Runs migrations on startup
- ✅ `backend/Dockerfile.prod` - Includes migrations in production image
- ✅ `backend/migrations/013_add_topic_difficulty_progress.sql` - New migration
- ✅ `backend/migrations/014_add_player_star_rewards.sql` - New migration

## Production Checklist

- [ ] Backup database
- [ ] Test on staging/development first
- [ ] Review migration files
- [ ] Deploy during low-traffic period
- [ ] Monitor logs during deployment
- [ ] Verify new tables exist
- [ ] Test application functionality
- [ ] Keep backup for 30 days

## Need More Details?

See `MIGRATION_GUIDE.md` for comprehensive documentation.
