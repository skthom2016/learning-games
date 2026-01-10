# Implementation TODO List

This document contains a comprehensive checklist for implementing the Learning Game Platform MVP. Each item references specific files and sections from `ARCHITECTURE.md`.

**IMPORTANT**: Read `ARCHITECTURE.md` in full before starting implementation.

---

## 📋 Progress Overview

- ✅ **Completed**: Project structure, boilerplate code, database schema
- 🟡 **In Progress**: Implementation (this checklist)
- ⬜ **Not Started**: Testing, deployment

---

## Backend Implementation

### 1. Database Layer (`backend/internal/database/`)

#### 1.1 Database Connection ✅ (Already Scaffolded)

- [x] Database initialization in `database.go`
- [x] Connection pooling setup

#### 1.2 Player Operations (`backend/internal/services/platform/platform.go`)

Reference: ARCHITECTURE.md Section 7 - Player Entity

- [ ] **CreatePlayer**:
  - Insert into `players` table
  - Generate UUID for `player_id`
  - Set default values (`total_session_count = 0`)
  - Return created player

- [ ] **ListPlayers**:
  - Query all players `ORDER BY created_at DESC`
  - Return array of players

- [ ] **GetPlayerByID**:
  - Query player by `player_id`
  - Return error if not found

- [ ] **UpdatePlayerLastPlayed**:
  - Update `last_played_at` timestamp
  - Increment `total_session_count`

**Test**: Create 2-3 test players, verify they appear in list

---

### 2. Game Engine (`backend/internal/services/game/multiplication.go`)

Reference: ARCHITECTURE.md Section 3 - Game Abstraction Model

- [ ] **Parse Topic ID**:
  - Extract multiplier from topic (e.g., "times-7" → 7)
  - Handle invalid format gracefully

- [ ] **GenerateQuestion**:
  - Implement difficulty logic:
    - Easy: operands 2-5
    - Medium: operands 2-12
    - Hard: operands 2-12 (same as medium for MVP)
  - Use seed for reproducible randomness
  - Generate visual hint data (rows × cols)
  - Calculate and store correct answer
  - Return Question struct

- [ ] **ValidateAnswer**:
  - Compare submitted answer with correct answer
  - Handle string/int conversion
  - Return AnswerValidationResult with:
    - `is_correct` boolean
    - `correct_answer` string
    - `feedback_message` (encouraging!)

- [ ] **GetGameDefinition**:
  - Return multiplication game metadata
  - Include all 11 topics (times-2 through times-12)
  - Include 3 difficulty levels

**Test**: Generate 10 questions, validate various answers

---

### 3. Adaptive Learning Engine (`backend/internal/services/learning/adaptive.go`)

Reference: ARCHITECTURE.md Section 4 - Learning Intelligence Engine

#### 3.1 Question Selection

- [ ] **SelectNextTopic**:
  - Calculate weights for each topic:
    - UNKNOWN: 100
    - WEAK: 200
    - LEARNING: 80
    - STRONG: 30
    - MASTERED: 5
  - Apply modifiers:
    - Not practiced in last 10: × 1.5
    - Practiced in last 3: × 0.5
    - Confidence recovery mode: WEAK × 3, others × 0.5
  - Normalize to probabilities
  - Use weighted random selection
  - Return selected topic_id

- [ ] **SelectDifficultyForTopic**:
  - UNKNOWN/WEAK: "easy"
  - LEARNING: "easy" if accuracy < 50%, else "medium"
  - STRONG: 70% "medium", 30% "hard" (random)
  - MASTERED: "hard"

**Test**: Create mastery records with various states, verify selection favors WEAK topics

#### 3.2 Mastery Assessment

- [ ] **AssessMastery**:
  - Get last 10 attempts for topic
  - Calculate accuracy
  - Apply state rules:
    - < 3 attempts: UNKNOWN
    - < 40% accuracy: WEAK
    - 40-69%: LEARNING
    - 70-89%: STRONG
    - ≥ 90% + 10 consecutive correct: MASTERED
  - Return MasteryState

- [ ] **UpdateMasteryRecord**:
  - After each attempt, recalculate mastery
  - Update `mastery_state`, `times_practiced`, `consecutive_correct`
  - Update `last_correct_at` or `last_incorrect_at`
  - Update `last_assessed_at` timestamp

**Test**: Create 15 attempts for topic (mix of correct/incorrect), verify state transitions

#### 3.3 Confidence Recovery Mode

- [ ] **CheckConfidenceRecoveryTrigger**:
  - Get last 7 attempts
  - Count incorrect (need ≥ 5)
  - Count unique topics (need ≥ 2)
  - Return true if conditions met and not already in recovery

- [ ] **CheckConfidenceRecoveryExit**:
  - Check for 3 consecutive correct
  - OR 5 of last 7 correct
  - OR 15 minutes elapsed since entry
  - Return true if any condition met

- [ ] **ActivateConfidenceRecoveryMode**:
  - Update PlayerGameProfile: `confidence_recovery_mode_active = true`
  - Set `confidence_recovery_entered_at = now()`
  - Log event for analytics

- [ ] **DeactivateConfidenceRecoveryMode**:
  - Update PlayerGameProfile: `confidence_recovery_mode_active = false`
  - Clear `confidence_recovery_entered_at`
  - Log event

**Test**: Create 7 attempts with 5 incorrect across 2 topics, verify recovery triggers

---

### 4. Reward Engine (`backend/internal/services/rewards/rewards.go`)

Reference: ARCHITECTURE.md Section 5 - Reward System Architecture

#### 4.1 Star Calculation

- [ ] **CalculateStarReward**:
  - Implement star matrix (see Section 5):
    ```
    UNKNOWN:  [5, 10, 15, 20]
    WEAK:     [5, 10, 15, 20]
    LEARNING: [3, 10, 15, 20]
    STRONG:   [2,  8, 12, 18]
    MASTERED: [1,  5,  8, 12]
    ```
  - Index by difficulty_tier (1-4)
  - Return star count

**Test**: Verify all matrix values, test edge cases (tier 0, tier 5)

#### 4.2 Reward Transactions

- [ ] **CreateRewardTransaction**:
  - Generate transaction_id (UUID)
  - Insert into `reward_transactions` table
  - Include all metadata (game_id, topic_id, attempt_id, difficulty_tier)
  - Flag `confidence_recovery_mode` if applicable
  - Set `consolidated = false`

- [ ] **GetRewardBalance**:
  - Query `reward_transactions` for player
  - Filter `reward_type = 'STAR'`
  - Sum all stars (total)
  - Sum where `consolidated = false` (unconsolidated)
  - Return both values

**Test**: Create 10 transactions, verify balance calculation

#### 4.3 Reward Consolidation (Admin)

- [ ] **ConsolidateRewards**:
  - Validate: `stars_to_consolidate <= unconsolidated_balance`
  - Start database transaction
  - Create RewardResetEvent record
  - Get oldest unconsolidated star transactions (ORDER BY earned_at ASC)
  - Update transactions: SET `consolidated = true`
  - Commit transaction
  - Return new balance

**Test**: Create 1000 stars, consolidate 500, verify balance updates correctly

---

### 5. Game Orchestrator (New File: `backend/internal/services/orchestrator/orchestrator.go`)

Reference: ARCHITECTURE.md Section 2 - GameOrchestrator

**Create New File**: This ties all services together

- [ ] **GetNextQuestion**:
  - Get PlayerGameProfile (or create if first time)
  - Get all MasteryRecords for player + game
  - Call AdaptiveLearningEngine.SelectNextTopic()
  - Call AdaptiveLearningEngine.SelectDifficultyForTopic()
  - Call GameEngine.GenerateQuestion()
  - Cache question in session (to prevent re-asking)
  - Return question to API handler

- [ ] **ProcessAnswer**:
  - Validate answer using GameEngine.ValidateAnswer()
  - **Start database transaction** (this is critical!)
  - Create AttemptRecord (immutable)
  - Update MasteryRecord (recalculate state)
  - Check confidence recovery triggers/exits
  - If correct:
    - Calculate stars using RewardEngine
    - Create RewardTransaction
  - Commit transaction
  - Return result (correct/incorrect, stars_earned, new_mastery_state)

**Test**: Submit 20 answers, verify all database records created correctly

---

### 6. API Handlers (`backend/internal/api/handlers/`)

Reference: ARCHITECTURE.md Section 8 - Backend API

Complete all TODOs in the following files:

#### 6.1 `player.go`

- [ ] CreatePlayer: Call platform.CreatePlayer(), return 201
- [ ] ListPlayers: Call platform.ListPlayers(), return 200
- [ ] GetPlayer: Call platform.GetPlayerByID(), return 200 or 404

#### 6.2 `session.go`

- [ ] StartSession: Create/load PlayerGameProfile, increment session_count
- [ ] EndSession: Update last_played_at, return session summary

#### 6.3 `gameplay.go`

- [ ] GetNextQuestion: Call orchestrator.GetNextQuestion()
- [ ] SubmitAnswer: Call orchestrator.ProcessAnswer()

#### 6.4 `progress.go`

- [ ] GetPlayerProgress: Query PlayerGameProfile stats
- [ ] GetMasteryHeatmap: Query MasteryRecords for all topics

#### 6.5 `rewards.go`

- [ ] GetRewardBalance: Call rewards.GetRewardBalance()
- [ ] GetRewardTransactions: Query reward_transactions table
- [ ] ResetRewards: Call admin.ConsolidateRewards()

#### 6.6 `admin.go`

- [ ] GetPlayerAnalytics: Aggregate stats across games

**Test**: Use `curl` or Postman to test all endpoints

---

## Frontend Implementation

### 7. API Client (`frontend/src/api/client.js`)

Reference: All TODO comments in file

- [ ] **Uncomment all API calls**:
  - Replace mock returns with actual axios calls
  - Test each endpoint individually

**Test**: Use browser DevTools Network tab to verify API calls

---

### 8. Child UI Screens (`frontend/src/screens/child/`)

#### 8.1 `PlayerSelectionScreen.js`

Reference: ARCHITECTURE.md Section 6 - Screen 1

- [ ] Fetch players from API (uncomment api.listPlayers())
- [ ] Display player cards with name + avatar
- [ ] Implement add player modal with form validation
- [ ] On player click, navigate to /games with playerId

**Test**: Add 3 players, verify they appear, click to navigate

#### 8.2 `GameSelectionScreen.js`

Reference: ARCHITECTURE.md Section 6 - Screen 2

- [ ] Fetch player data and star count
- [ ] Display greeting with player name
- [ ] On "Let's Play", call api.startSession()
- [ ] Navigate to /play with playerId + gameId

**Test**: Select player, verify star count shown, start game

#### 8.3 `GamePlayScreen.js`

Reference: ARCHITECTURE.md Section 6 - Screen 4

- [ ] Fetch next question (uncomment api.getNextQuestion())
- [ ] Display question text
- [ ] Implement number input with validation
- [ ] Implement "Show Hint" button (conditional on difficulty)
- [ ] Render visual hint (dot array CSS grid)
- [ ] Track time to answer (Date.now() diff)
- [ ] On submit, call api.submitAnswer()
- [ ] Navigate to /feedback with result

**Test**: Answer 10 questions, verify hints work, check timing

#### 8.4 `FeedbackScreen.js`

Reference: ARCHITECTURE.md Section 6 - Screen 5

- [ ] Display celebration (if correct) or encouragement (if incorrect)
- [ ] Show stars earned (if correct)
- [ ] Animate star earning (CSS keyframes)
- [ ] Auto-advance to next question after 2.5s

**Test**: Submit correct and incorrect answers, verify feedback

#### 8.5 `ProgressOverviewScreen.js`

Reference: ARCHITECTURE.md Section 6 - Screen 7

- [ ] Fetch mastery heatmap (uncomment api.getMasteryHeatmap())
- [ ] Display star count
- [ ] Show topic list with mastery icons
- [ ] Keep it simple (no percentages)

**Test**: View progress, verify topics shown with correct states

---

### 9. Admin UI Screens (`frontend/src/screens/admin/`)

#### 9.1 `AdminDashboard.js`

Reference: ARCHITECTURE.md Section 6 - Screen 8

- [ ] Fetch all players
- [ ] For each player, fetch star count and mastery summary
- [ ] Display player cards
- [ ] Implement navigation to progress and reset screens

**Test**: View dashboard, verify all players shown

#### 9.2 `PlayerProgressView.js`

Reference: ARCHITECTURE.md Section 6 - Screen 9

- [ ] Fetch player progress stats
- [ ] Fetch mastery heatmap
- [ ] Display summary (total questions, accuracy, sessions)
- [ ] Render mastery heatmap with color coding:
  - Green: MASTERED, STRONG
  - Yellow: LEARNING
  - Red: WEAK
  - Gray: UNKNOWN
- [ ] Auto-generate recommendations

**Test**: View progress for player, verify heatmap colors match states

#### 9.3 `RewardResetScreen.js`

Reference: ARCHITECTURE.md Section 6 - Screen 11

- [ ] Fetch reward balance
- [ ] Validate stars to consolidate ≤ balance
- [ ] Require reward description
- [ ] Confirm before submitting
- [ ] Call api.resetRewards()
- [ ] Display consolidation history

**Test**: Reset 500 stars, verify balance updates, check history

---

## Testing

### 10. Backend Tests

- [ ] Create `backend/internal/services/learning/adaptive_test.go`
- [ ] Test mastery state transitions
- [ ] Test question selection weights
- [ ] Test confidence recovery triggers

- [ ] Create `backend/internal/services/rewards/rewards_test.go`
- [ ] Test star matrix values
- [ ] Test reward consolidation

- [ ] Create `backend/internal/services/game/multiplication_test.go`
- [ ] Test question generation
- [ ] Test answer validation

**Run**: `cd backend && go test ./...`

---

### 11. Frontend Tests

- [ ] Create tests for API client
- [ ] Create tests for screen components (basic rendering)

**Run**: `cd frontend && npm test`

---

## Integration Testing

### 12. End-to-End User Flow

- [ ] **Child Flow**:
  1. Select player
  2. Start game
  3. Answer 20 questions (mix correct/incorrect)
  4. View progress
  5. Verify mastery states updated

- [ ] **Admin Flow**:
  1. View dashboard
  2. Check player progress
  3. Reset 500 stars
  4. Verify balance updated

- [ ] **Adaptive Behavior**:
  1. Create player
  2. Answer 5 questions incorrectly on "times-7"
  3. Verify next questions favor "times-7" (WEAK state)
  4. Answer 10 correctly
  5. Verify state changes to STRONG

- [ ] **Confidence Recovery**:
  1. Answer 5 of 7 incorrectly across 2+ topics
  2. Verify difficulty locks to easy
  3. Answer 3 correctly
  4. Verify recovery mode exits

---

## Styling & Polish

### 13. Frontend UI Polish

- [ ] **Child UI**:
  - Large buttons (min 60px height)
  - Bright colors
  - Large fonts (min 1.2rem)
  - Rounded corners (border-radius: 10px+)
  - Touch-friendly spacing

- [ ] **Animations**:
  - Star earn animation (CSS keyframes)
  - Button hover effects
  - Smooth transitions

- [ ] **Responsive**:
  - Test on 1024×768 (tablet)
  - Test on 1920×1080 (desktop)

**Test**: Have a child (or adult) use the UI, observe confusion points

---

## Documentation

### 14. Code Documentation

- [ ] Add comments to all public functions (Golang)
- [ ] Add JSDoc comments to React components
- [ ] Document API endpoints in `routes.go`

### 15. Update TODO.md

- [ ] Mark completed items as done
- [ ] Add any discovered TODOs

---

## Deployment Verification

### 16. Docker Compose

- [ ] Test fresh install:
  ```bash
  docker-compose down -v
  docker-compose up --build
  ```
- [ ] Verify all services start
- [ ] Verify database migrations run
- [ ] Verify seed data loads

- [ ] Test persistence:
  1. Create player
  2. Restart Docker: `docker-compose restart`
  3. Verify player still exists

---

## MVP Acceptance Criteria

### 17. Final Checklist

Reference: ARCHITECTURE.md Section 8 - MVP Acceptance Criteria

**Functional**:
- [ ] 2 players can be created
- [ ] Each player can practice multiplication independently
- [ ] System asks more questions from weak topics
- [ ] Stars are earned based on difficulty + mastery
- [ ] Parent can view which topics child struggles with
- [ ] Parent can reset star balance
- [ ] All data persists across restarts

**Quality**:
- [ ] No crashes during 30-minute play session
- [ ] Questions load within 200ms
- [ ] UI readable on 1024×768 screen
- [ ] Database migrations run cleanly

**User Validation**:
- [ ] Child (age 6-9) can navigate UI without help (after 1 demo)
- [ ] Child reports game is "fun"
- [ ] Parent understands progress view

---

## Estimated Time per Section

Based on junior developer skill level:

| Section | Estimated Time |
|---------|----------------|
| Database layer (player ops) | 4 hours |
| Game engine | 6 hours |
| Adaptive learning | 10 hours |
| Reward system | 6 hours |
| Game orchestrator | 6 hours |
| API handlers | 8 hours |
| Frontend API client | 2 hours |
| Child UI screens | 16 hours |
| Admin UI screens | 12 hours |
| Testing | 16 hours |
| Polish & documentation | 8 hours |
| **TOTAL** | **94 hours** (~2.5 weeks full-time) |

---

## Tips for Success

1. **Read ARCHITECTURE.md first** — Don't skip this!
2. **Work in order** — Backend before frontend
3. **Test incrementally** — Don't build everything then test
4. **Use the database** — Check `psql` to verify data
5. **Check browser console** — Catch frontend errors early
6. **Ask questions** — If architecture is unclear, clarify before implementing
7. **Keep it simple** — MVP means minimal, don't add features

---

## Getting Help

- **Architecture questions**: Refer to ARCHITECTURE.md sections
- **Implementation questions**: Check TODO comments in code
- **Bugs**: Check Docker logs (`docker-compose logs [service]`)
- **Database issues**: Connect via psql and inspect data

---

## Next Steps After MVP

Once all items above are checked:

1. Deploy to test family (1-2 children)
2. Observe gameplay for 1-2 weeks
3. Collect feedback
4. Plan Phase 2 features (see ARCHITECTURE.md Section 8)

**Good luck! 🚀**
