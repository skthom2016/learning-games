# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Kid-friendly adaptive learning game platform (ages 6-9) for math practice (multiplication, division, addition, subtraction). Built with React frontend, Go backend, PostgreSQL database, deployed via Docker Compose.

## Essential Commands

### Development Workflow

```bash
# Start all services (backend, frontend, database)
docker-compose up --build

# Access points
# - Child UI: http://localhost:3000
# - Admin UI: http://localhost:3000/admin
# - Backend API: http://localhost:8080
# - Health check: http://localhost:8080/health

# Stop services
docker-compose down

# Reset database (delete all data)
docker-compose down -v
docker-compose up --build
```

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run locally (requires PostgreSQL on localhost:5432)
export DB_HOST=localhost DB_PORT=5432 DB_USER=gameuser DB_PASSWORD=gamepass123 DB_NAME=learning_game SERVER_PORT=8080
go run cmd/server/main.go

# Hot reload is enabled in Docker via Air (.air.toml)
# Any .go file changes auto-rebuild (excludes *_test.go)

# Run tests (when implemented)
go test ./...

# Test specific package
go test ./internal/services/learning

# Test with verbose output
go test -v ./internal/services/game
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run locally (connects to localhost:8080 backend)
npm start

# Run tests (when implemented)
npm test

# Build production bundle
npm build
```

### API Testing

```bash
# Run comprehensive API test suite (Windows)
test-api.bat

# Test API manually
curl http://localhost:8080/health
curl http://localhost:8080/api/players
curl -X POST http://localhost:8080/api/players -H "Content-Type: application/json" -d '{"player_name": "TestPlayer"}'
```

### Database Operations

```bash
# Connect to PostgreSQL container
docker exec -it learning-game-db psql -U gameuser -d learning_game

# View applied migrations
docker exec -it learning-game-db psql -U gameuser -d learning_game -c "SELECT * FROM schema_migrations ORDER BY id;"

# Backup database
docker exec learning-game-db pg_dump -U gameuser learning_game > backup.sql

# Restore database
docker exec -i learning-game-db psql -U gameuser learning_game < backup.sql

# Manual migration (usually not needed - migrations run automatically)
docker exec -i learning-game-db psql -U gameuser -d learning_game < backend/migrations/014_add_player_star_rewards.sql
```

**Automatic Migrations:**
- Migrations now run automatically on backend startup
- Tracked in `schema_migrations` table
- Only pending migrations are applied (never re-runs old migrations)
- Safe for existing databases with data
- Each migration runs in a transaction (rolls back on error)
- See `MIGRATION_GUIDE.md` for detailed documentation

## Architecture Overview

### Backend Service Layer (Go)

The backend follows a **layered service-oriented architecture** with central orchestration:

```
API Handlers (thin controllers)
    ↓
GameOrchestrator (central coordinator)
    ↓
┌──────────────┬─────────────────┬──────────────┬──────────────┐
│ GameEngine   │ LearningEngine  │ RewardEngine │ PlatformSvc  │
│ (question    │ (adaptive AI)   │ (star calc)  │ (player mgmt)│
│  generation) │                 │              │              │
└──────────────┴─────────────────┴──────────────┴──────────────┘
    ↓
Database Layer (PostgreSQL)
```

**Key Components:**

1. **GameOrchestrator** (`backend/internal/services/orchestrator/orchestrator.go`)
   - Central hub coordinating all game flow
   - Manages complete request lifecycle: GetNextQuestion → ProcessAnswer
   - Creates singleton game engine instances per game type
   - Orchestrates: topic selection → difficulty selection → question generation → answer validation → mastery update → reward calculation

2. **AdaptiveLearningEngine** (`backend/internal/services/learning/adaptive.go`)
   - `SelectNextTopic()`: Weighted random selection based on mastery state (WEAK=200, UNKNOWN=100, LEARNING=80, STRONG=30, MASTERED=5)
   - `SelectDifficultyForTopic()`: Maps mastery state to difficulty tier
   - `AssessMastery()`: Calculates mastery state from attempt history (accuracy thresholds over last 10 attempts)
   - `CheckConfidenceRecoveryTrigger()`: Activates when 5+ errors in last 7 attempts
   - `CheckConfidenceRecoveryExit()`: Exits after 3 consecutive correct OR 5/7 correct OR 15 min elapsed

3. **RewardEngine** (`backend/internal/services/rewards/rewards.go`)
   - `CalculateStarReward()`: Matrix-based calculation (mastery state × difficulty tier)
   - Recovery mode bonus: 2× star multiplier for emotional support
   - `CreateRewardTransaction()`: Immutable ledger logging
   - `ConsolidateRewards()`: FIFO consolidation for admin reward redemption

4. **Game Engines** (`backend/internal/services/game/`)
   - Implement `GameEngine` interface: GenerateQuestion(), ValidateAnswer(), GetGameDefinition()
   - Each game (multiplication, division, addition, subtraction) has separate logic
   - Support custom number ranges via NumberRangesService
   - Generate visual hints (grouping diagrams) and verbal hints (pedagogical suggestions)

5. **Platform Service** (`backend/internal/services/platform/platform.go`)
   - Player CRUD: CreatePlayer, ListPlayers, GetPlayerByID, DeletePlayer
   - Session management: StartGameSession, EndGameSession
   - Cascading deletes respect foreign key constraints

### Database Schema Patterns

**Immutable Event Sourcing:**
- `attempt_records`: Append-only log of every question attempt (never updated)
- `reward_transactions`: Immutable ledger of all stars earned (consolidated flag for redemption)
- `reward_reset_events`: Audit trail of admin consolidations

**Mutable State:**
- `mastery_records`: Current mastery state per topic (UNKNOWN → WEAK → LEARNING → STRONG → MASTERED)
- `player_game_profiles`: Aggregate stats and confidence recovery mode state
- `players`: Player metadata

**Key Relationships:**
```
players (1) ──┬── (N) player_game_profiles
              ├── (N) mastery_records
              ├── (N) attempt_records
              ├── (N) reward_transactions
              └── (N) player_number_ranges

games (1) ──┬── (N) topics
            └── (N) difficulty_levels
```

**Migration Strategy:**
- Migrations in `backend/migrations/` numbered sequentially (001, 002, 003, ...)
- Run automatically on EVERY backend startup (via `backend/internal/database/migrations.go`)
- Tracked in `schema_migrations` table (only pending migrations are applied)
- Each migration runs in a transaction (atomic, rolls back on error)
- Use PostgreSQL ENUMs for type safety (mastery_state, reward_type)
- Indexes on hot query paths (attempts by player+game+topic+time, rewards by player+time)
- Safe for existing databases with data (never re-runs applied migrations)

### Frontend Architecture (React)

**Single-Page Application** with React Router v6:

```
AuthProvider (React Context for admin login)
    ↓
Router
    ↓
┌─────────────────────────┬──────────────────────────┐
│ Public Routes (Child)   │ Protected Routes (Admin) │
│ - PlayerSelection       │ - AdminLogin             │
│ - GameSelection         │ - AdminDashboard         │
│ - GamePlay              │ - PlayerProgressView     │
│ - Feedback              │ - RewardReset            │
│ - ProgressOverview      │ - NumberRangeSettings    │
└─────────────────────────┴──────────────────────────┘
```

**Key Patterns:**
- `frontend/src/api/client.js`: Centralized axios instance with interceptors
- `frontend/src/contexts/AuthContext.js`: Admin authentication state management
- `frontend/src/components/ProtectedRoute.js`: Route guard for admin pages
- All screens are stateful components making direct API calls

**Screen Flow (Child):**
1. PlayerSelectionScreen → select player
2. GameSelectionScreen → choose game type
3. GamePlayScreen → GET /questions/next → POST /answers/submit
4. FeedbackScreen → show results + stars + mastery progress
5. Loop back to GamePlayScreen

**Screen Flow (Admin):**
1. AdminLoginScreen → authenticate
2. AdminDashboard → player list + quick stats
3. PlayerProgressView → detailed mastery heatmap, attempt history
4. NumberRangeSettingsScreen → customize difficulty ranges per player
5. RewardReset → consolidate stars (mark as redeemed)

## Critical Architectural Patterns

### 1. Confidence Recovery System

**Purpose:** Emotional support when child is struggling

**Trigger** (`adaptive.go:185-214`):
- Requires: ≥10 total attempts, ≥5 incorrect in last 7, ≥2 different topics
- Sets `player_game_profiles.confidence_recovery_mode_active = true`

**Active Behavior:**
- Topic selection: WEAK topics get 3× weight (focus on gaps)
- Difficulty: Forces "easy" tier
- Rewards: 2× star multiplier (motivation boost)

**Exit Conditions** (`adaptive.go:217-261`):
- 3 consecutive correct answers, OR
- 5 out of last 7 correct, OR
- 15 minutes elapsed

### 2. Weighted Topic Selection

Algorithm (`adaptive.go:14-102`):

```
1. If all topics UNKNOWN → return first (sequential start)
2. Calculate weight per topic:
   - Base: WEAK=200, UNKNOWN=100, LEARNING=80, STRONG=30, MASTERED=5
   - If practiced recently: ×0.5
   - If in recovery mode:
     - WEAK topics: ×3.0 (focus)
     - Others: ×0.5 (avoid overload)
3. Normalize to probability distribution
4. Weighted random selection
```

**Rationale:** Balances focus on weak areas with topic variety to prevent monotony.

### 3. Mastery State Transitions

Based on accuracy over last 10 attempts (`adaptive.go:128-182`):

```
≥90% accuracy + ≥20 total + 10 consecutive correct → MASTERED
<40% accuracy                                      → WEAK
40-70% accuracy                                    → LEARNING
70-90% accuracy                                    → STRONG
≥90% accuracy (not mastered yet)                   → STRONG
```

### 4. Reward Matrix

Stars awarded by mastery state × difficulty tier (`rewards.go:14-47`):

```
             Tier1  Tier2  Tier3  Tier4+
UNKNOWN       5      10     15     20
WEAK          5      10     15     20
LEARNING      3      10     15     20
STRONG        2       8     12     18
MASTERED      1       5      8     12
```

- Recovery mode: 2× multiplier
- Design intent: Motivate new/weak topics, prevent grinding mastered topics

### 5. Immutable Ledger Pattern

**Why:**
- Auditability: Every reward/attempt has full context
- Reversibility: Can replay history by transaction ID
- Compliance: Audit trail for admin actions

**Implementation:**
- `attempt_records`: Never UPDATE, only INSERT
- `reward_transactions`: Never DELETE, mark `consolidated=true`
- `reward_reset_events`: Links consolidated transactions to real-world rewards

## Data Flow: Complete Game Loop

```
Child clicks "Next Question"
    ↓
GET /questions/next?player_id=X&game_id=Y
    ↓
GameOrchestrator.GetNextQuestion()
    ├─ StartGameSession() → Create/load PlayerGameProfile
    ├─ Get MasteryRecords for player+game
    ├─ If first time → Initialize UNKNOWN mastery for all topics
    ├─ SelectNextTopic() → Weighted selection
    ├─ SelectDifficultyForTopic() → Map mastery→difficulty
    ├─ GenerateQuestion() → Create question with hints
    └─ Return Question (WITHOUT correct answer)
    ↓
Frontend displays question + visual/verbal hints
    ↓
Child submits answer
    ↓
POST /answers/submit
    ↓
GameOrchestrator.ProcessAnswer()
    ├─ BEGIN TRANSACTION
    ├─ ValidateAnswer()
    ├─ INSERT AttemptRecord (immutable log)
    ├─ Get all attempts for this topic
    ├─ AssessMastery() → Calculate new state
    ├─ UPDATE/INSERT MasteryRecord
    ├─ CheckConfidenceRecoveryTrigger()
    ├─ If correct:
    │   ├─ CalculateStarReward()
    │   └─ INSERT RewardTransaction (immutable ledger)
    ├─ UPDATE PlayerGameProfile stats
    ├─ COMMIT TRANSACTION
    └─ Return ProcessAnswerResponse
    ↓
FeedbackScreen shows:
    ├─ Correctness feedback
    ├─ Stars earned
    ├─ Mastery progress
    └─ [Next] button → Loop
```

## Common Development Tasks

### Adding a New Game Type

1. Create game engine in `backend/internal/services/game/newgame.go`:
   - Implement `GameEngine` interface
   - Define topic IDs (e.g., "multiply-7", "divide-12")
   - Implement GenerateQuestion() with visual/verbal hints
   - Implement ValidateAnswer()

2. Add migration in `backend/migrations/00X_seed_newgame.sql`:
   - INSERT INTO games (game_id, game_name, ...)
   - INSERT INTO topics (topic_id, game_id, topic_name, ...)
   - INSERT INTO difficulty_levels (difficulty_id, game_id, ...)

3. Register in orchestrator (`orchestrator.go:NewGameOrchestrator`):
   - Add case for new game_id
   - Instantiate game engine

4. Add frontend UI:
   - Create game card in GameSelectionScreen
   - Update GamePlayScreen routing

### Modifying Adaptive Learning Algorithm

**Files to change:**
- `backend/internal/services/learning/adaptive.go`

**Key functions:**
- `SelectNextTopic()`: Adjust base weights or modifiers
- `SelectDifficultyForTopic()`: Change mastery→difficulty mapping
- `AssessMastery()`: Modify accuracy thresholds or transition logic
- `CheckConfidenceRecoveryTrigger()`: Adjust trigger conditions
- `CheckConfidenceRecoveryExit()`: Modify exit criteria

**Testing:**
- Create test attempts with specific patterns
- Verify mastery state transitions
- Check confidence recovery activation/deactivation

### Adjusting Reward Matrix

**File:** `backend/internal/services/rewards/rewards.go:14-47`

**Matrix structure:**
```go
var starMatrix = map[MasteryState]map[int]int{
    MasteryStateUnknown:  {1: 5, 2: 10, 3: 15, 4: 20},
    MasteryStateWeak:     {1: 5, 2: 10, 3: 15, 4: 20},
    // ... modify values here
}
```

**Recovery mode multiplier:** Line 44-47 (currently 2×)

### Adding Database Fields

1. Create migration: `backend/migrations/00X_add_field.sql`
2. Update model structs: `backend/internal/models/models.go`
3. Update service layer queries
4. Test with fresh database: `docker-compose down -v && docker-compose up --build`

### Debugging Common Issues

**Backend not responding:**
```bash
docker-compose logs backend
# Check Air rebuild logs
# Verify database connection
curl http://localhost:8080/health
```

**Frontend API calls failing:**
- Check `REACT_APP_API_URL` in docker-compose.yml (should be localhost:8080 for browser)
- Open browser DevTools → Network tab
- Check `frontend/src/api/client.js` axios interceptors

**Database migration issues:**
```bash
# Migrations run automatically on EVERY backend startup
# Check migration logs:
docker-compose logs backend | grep migration

# View which migrations have been applied:
docker exec -it learning-game-db psql -U gameuser -d learning_game -c "SELECT * FROM schema_migrations;"

# To re-run ALL migrations (WARNING: deletes all data):
docker-compose down -v  # Delete volume
docker-compose up --build
```

**Hot reload not working:**
- Backend: Check `.air.toml` include/exclude patterns
- Frontend: Verify volume mount in docker-compose.yml
- Windows: May need CHOKIDAR_USEPOLLING=true for React

## Important Notes

### Database Credentials
- **Development only**: `gameuser:gamepass123`
- **Never commit production credentials**
- Change in `docker-compose.yml` environment variables + backend connection

### Migration Best Practices
- Always number sequentially (`001_`, `002_`, ...)
- Include rollback comments for reversibility
- Migrations run automatically on backend startup (tracked in `schema_migrations` table)
- Test with fresh database: `docker-compose down -v && docker-compose up --build`
- Test with existing database: `docker-compose restart backend` (only new migrations run)
- Backup before production migrations: `pg_dump -U gameuser learning_game > backup.sql`
- See `MIGRATION_GUIDE.md` for deployment instructions

### Foreign Key Cascading
- `ON DELETE CASCADE` used for child records (attempts, rewards, mastery)
- Deleting a player deletes ALL associated data
- Admin UI should warn before deletion

### Transaction Boundaries
- `ProcessAnswer()` wraps all operations in a single transaction
- Ensures atomicity: attempt + mastery update + reward creation
- Rollback on any error prevents partial state

### Custom Number Ranges
- Feature for grade-level customization
- Admin sets per player, per game
- If set, overrides default difficulty ranges in question generation
- Check `numberranges.GetNumberRangeForQuestionGeneration()` in game engines

### Recovery Mode Edge Cases
- Can stay active indefinitely if child struggles
- 15-minute timeout prevents infinite easy mode
- Design intent: short-term emotional support, not permanent crutch
