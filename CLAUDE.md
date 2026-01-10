# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A kid-friendly adaptive learning game platform for children aged 6-9, focused on multiplication practice. The system uses adaptive learning to automatically focus on weak areas, with a reward system to motivate learning.

**Tech Stack:**
- Frontend: React 18 with React Router
- Backend: Golang 1.21 with Gin framework
- Database: PostgreSQL 15
- Deployment: Docker Compose

## Common Development Commands

### Starting the Full Stack
```bash
docker-compose up --build
```
This starts PostgreSQL (port 5432), backend API (port 8080), and React frontend (port 3000).

### Backend Development
```bash
cd backend
go run cmd/server/main.go
```
Requires environment variables: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `SERVER_PORT`.

### Frontend Development
```bash
cd frontend
npm install
npm start
```

### Testing
```bash
# Backend
cd backend
go test ./...

# Frontend
cd frontend
npm test
```

### Database Operations
```bash
# Connect to database
docker exec -it learning-game-db psql -U gameuser -d learning_game

# Reset all data
docker-compose down -v
docker-compose up --build
```

## Architecture Overview

### Backend Service Layer Architecture

The backend follows a strict service-oriented architecture with clear separation of concerns:

1. **Platform Service** (`backend/internal/services/platform/platform.go`)
   - Player CRUD operations
   - Session management (StartGameSession, EndGameSession)
   - Creates/updates PlayerGameProfile records

2. **Game Engine** (`backend/internal/services/game/multiplication.go`)
   - Generates questions based on topic and difficulty
   - Parses topic IDs (e.g., "times-7" → multiplier 7)
   - Validates answers and provides feedback
   - Generates visual hints (dot array representation)

3. **Adaptive Learning Engine** (`backend/internal/services/learning/adaptive.go`)
   - **Mastery States**: UNKNOWN → WEAK → LEARNING → STRONG → MASTERED
   - **Topic Selection**: Uses weighted random selection favoring weak areas
   - **Confidence Recovery Mode**: Activates when 5+ of last 7 answers are incorrect across 2+ topics
   - **Mastery Assessment**: Recalculates state after each attempt using last 10 attempts

4. **Reward Engine** (`backend/internal/services/rewards/rewards.go`)
   - Star calculation based on mastery state + difficulty tier matrix
   - Creates RewardTransaction records
   - Supports reward consolidation (admin feature)

5. **Game Orchestrator** (`backend/internal/services/orchestrator/orchestrator.go`)
   - **CRITICAL**: Ties all services together
   - `GetNextQuestion()`: Coordinates adaptive learning + game engine
   - `ProcessAnswer()`: Wraps answer processing in database transaction
   - Must update: AttemptRecord, MasteryRecord, PlayerGameProfile, RewardTransaction (if correct)

### Frontend Structure

- **Child UI** (`frontend/src/screens/child/`)
  - `PlayerSelectionScreen.js`: Choose/create player
  - `GameSelectionScreen.js`: Shows available games and star count
  - `GamePlayScreen.js`: Displays questions, handles input, shows hints
  - `FeedbackScreen.js`: Celebration/encouragement after each answer
  - `ProgressOverviewScreen.js`: Shows mastery heatmap

- **Admin UI** (`frontend/src/screens/admin/`)
  - `AdminDashboard.js`: Overview of all players
  - `PlayerProgressView.js`: Detailed progress per player with heatmap
  - `RewardResetScreen.js`: Consolidate stars (real-world reward redemption)

### Key Data Flow

1. **Game Start**: `StartGameSession()` → creates/loads `PlayerGameProfile`
2. **Question Generation**: `GetNextQuestion()` → adaptive learning selects topic → game engine generates question
3. **Answer Submission**: `ProcessAnswer()` → validates → updates mastery → checks confidence recovery → awards stars (if correct)
4. **Session End**: `EndGameSession()` → returns session summary

## Database Schema

Key tables:
- `players`: Player records
- `player_game_profiles`: Per-player-per-game progress tracking
- `mastery_records`: Per-player-per-game-per-topic mastery state
- `attempt_records`: Immutable record of each answer attempt
- `reward_transactions`: Star earnings with consolidation flag
- `reward_reset_events`: Admin-initiated star consolidations

Database migrations in `backend/migrations/` run automatically on Docker startup.

## Critical Implementation Notes

### Adaptive Learning Weights
Topic selection uses these base weights:
- UNKNOWN: 100
- WEAK: 200 (highest priority)
- LEARNING: 80
- STRONG: 30
- MASTERED: 5 (lowest priority)

### Star Reward Matrix
Stars awarded based on mastery state + difficulty tier (1-4):
- UNKNOWN/WEAK: [5, 10, 15, 20]
- LEARNING: [3, 10, 15, 20]
- STRONG: [2, 8, 12, 18]
- MASTERED: [1, 5, 8, 12]

### Confidence Recovery Mode
- **Trigger**: 5+ incorrect in last 7 attempts across 2+ topics
- **Effect**: Locks difficulty to easy, heavily weights weak topics
- **Exit**: 3 consecutive correct OR 5 of 7 correct OR 15 minutes elapsed

### Database Transaction Critical Section
`ProcessAnswer()` in the orchestrator MUST wrap all updates in a single transaction:
1. Create AttemptRecord (immutable)
2. Update MasteryRecord (recalculate mastery state)
3. Update PlayerGameProfile (counts, confidence recovery mode)
4. If correct: Create RewardTransaction
5. Commit only if all succeed

## API Endpoints

All endpoints are prefixed with `/api`:
- `POST /players` - Create player
- `GET /players` - List all players
- `GET /players/:id` - Get player by ID
- `POST /sessions/start` - Start game session
- `POST /sessions/end` - End game session
- `GET /questions/next?player_id=&game_id=` - Get next question
- `POST /answers/submit` - Submit answer
- `GET /progress/:player_id/:game_id` - Get player progress
- `GET /mastery/:player_id/:game_id` - Get mastery heatmap
- `GET /rewards/:player_id` - Get reward balance
- `POST /admin/rewards/reset` - Consolidate rewards (admin)

## Development Principles

1. **Adaptive Learning is Core**: All game features must support the adaptive engine's mastery tracking
2. **No Punishment**: The UI never penalizes wrong answers—only encouragement
3. **Child-Safe Design**: Large buttons, bright colors, no timers, intuitive navigation
4. **Data Immutability**: AttemptRecords are never modified—only create new records
5. **Transaction Safety**: Any multi-table update must use database transactions

## Common Issues

### Port Conflicts
If ports 3000, 8080, or 5432 are in use, modify `docker-compose.yml`.

### Hot Reload
Backend uses Air for hot reload in Docker. Frontend uses React Scripts dev server.

### Database Connection
Backend depends on database health check. If backend fails to start, check `docker-compose logs db`.
