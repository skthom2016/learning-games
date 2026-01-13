package orchestrator

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
	"github.com/learning-game/backend/internal/services/game"
	"github.com/learning-game/backend/internal/services/learning"
	"github.com/learning-game/backend/internal/services/platform"
	"github.com/learning-game/backend/internal/services/rewards"
)

// GameEngine interface for all game implementations
type GameEngine interface {
	GenerateQuestion(playerID, topicID, difficultyLevelID string, seed int64) (*models.Question, error)
	ValidateAnswer(correctAnswer, submittedAnswer string) (*models.AnswerValidationResult, error)
	GetGameDefinition() *models.Game
}

// GameOrchestrator coordinates all services to manage game flow
type GameOrchestrator struct {
	gameEngines     map[string]GameEngine
	learningEngine  *learning.AdaptiveLearningEngine
	rewardEngine    *rewards.RewardEngine
}

// NewGameOrchestrator creates a new orchestrator instance
func NewGameOrchestrator() *GameOrchestrator {
	engines := make(map[string]GameEngine)
	engines["multiplication-tables"] = &game.MultiplicationGame{}
	engines["division-facts"] = &game.DivisionGame{}
	engines["addition-facts"] = &game.AdditionGame{}
	engines["subtraction-facts"] = &game.SubtractionGame{}

	return &GameOrchestrator{
		gameEngines:     engines,
		learningEngine:  &learning.AdaptiveLearningEngine{},
		rewardEngine:    &rewards.RewardEngine{},
	}
}

// getGameEngine returns the appropriate game engine for the given game ID
func (o *GameOrchestrator) getGameEngine(gameID string) (GameEngine, error) {
	engine, ok := o.gameEngines[gameID]
	if !ok {
		return nil, fmt.Errorf("no game engine found for game_id: %s", gameID)
	}
	return engine, nil
}

// GetNextQuestion generates the next question for a player
func (o *GameOrchestrator) GetNextQuestion(playerID, gameID string) (*models.Question, error) {
	db := database.GetDB()

	// Get or create PlayerGameProfile
	profile, err := platform.StartGameSession(playerID, gameID)
	if err != nil {
		return nil, err
	}

	// Get all mastery records for this player + game
	masteryRecords, err := o.getMasteryRecords(playerID, gameID)
	if err != nil {
		return nil, err
	}

	// If no mastery records exist, create initial ones for all topics
	if len(masteryRecords) == 0 {
		masteryRecords, err = o.initializeMasteryRecords(playerID, gameID)
		if err != nil {
			return nil, err
		}
	}

	// Get recent topics (last 5 attempts)
	recentTopics, err := o.getRecentTopics(playerID, gameID, 5)
	if err != nil {
		return nil, err
	}

	// Select next topic using adaptive learning engine
	selectedTopicID, err := o.learningEngine.SelectNextTopic(
		masteryRecords,
		recentTopics,
		profile.ConfidenceRecoveryModeActive,
	)
	if err != nil {
		return nil, err
	}

	// Find the mastery record for the selected topic
	var selectedMastery *models.MasteryRecord
	for _, mr := range masteryRecords {
		if mr.TopicID == selectedTopicID {
			selectedMastery = mr
			break
		}
	}

	if selectedMastery == nil {
		return nil, fmt.Errorf("mastery record not found for topic %s", selectedTopicID)
	}

	// Get recent accuracy for difficulty selection
	recentAccuracy, err := o.getRecentAccuracy(playerID, gameID, selectedTopicID, 5)
	if err != nil {
		return nil, err
	}

	// Select difficulty level
	difficultyLevelID := o.learningEngine.SelectDifficultyForTopic(selectedMastery, recentAccuracy)

	// Get the appropriate game engine
	gameEngine, err := o.getGameEngine(gameID)
	if err != nil {
		return nil, err
	}

	// Generate question using game engine
	seed := time.Now().UnixNano()
	question, err := gameEngine.GenerateQuestion(playerID, selectedTopicID, difficultyLevelID, seed)
	if err != nil {
		return nil, err
	}

	// Store question in cache (in-memory for now, could be Redis later)
	// This prevents re-asking the same question if they refresh
	err = o.cacheQuestion(playerID, question)
	if err != nil {
		return nil, err
	}

	_ = db
	return question, nil
}

// ProcessAnswerRequest contains the answer submission data
type ProcessAnswerRequest struct {
	PlayerID          string `json:"player_id"`
	GameID            string `json:"game_id"`
	QuestionID        string `json:"question_id"`
	TopicID           string `json:"topic_id"`
	DifficultyLevelID string `json:"difficulty_level_id"`
	SubmittedAnswer   string `json:"submitted_answer"`
	CorrectAnswer     string `json:"correct_answer"`
	TimeToAnswerMs    *int   `json:"time_to_answer_ms"`
	HintUsed          bool   `json:"hint_used"`
}

// ProcessAnswerResponse contains the result of answer processing
type ProcessAnswerResponse struct {
	IsCorrect         bool                 `json:"is_correct"`
	CorrectAnswer     string               `json:"correct_answer"`
	FeedbackMessage   string               `json:"feedback_message"`
	StarsEarned       int                  `json:"stars_earned"`
	NewMasteryState   models.MasteryState  `json:"new_mastery_state"`
	RecoveryModeActive bool                `json:"recovery_mode_active"`
}

// ProcessAnswer handles answer submission and all related updates
func (o *GameOrchestrator) ProcessAnswer(req ProcessAnswerRequest) (*ProcessAnswerResponse, error) {
	db := database.GetDB()

	// Get the appropriate game engine
	gameEngine, err := o.getGameEngine(req.GameID)
	if err != nil {
		return nil, err
	}

	// Validate answer
	validationResult, err := gameEngine.ValidateAnswer(req.CorrectAnswer, req.SubmittedAnswer)
	if err != nil {
		return nil, err
	}

	// Start database transaction (critical for data consistency)
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get difficulty tier for reward calculation
	difficultyTier, err := o.getDifficultyTier(req.GameID, req.DifficultyLevelID)
	if err != nil {
		return nil, err
	}

	// Get player game profile for recovery mode status
	profile, err := o.getPlayerGameProfile(req.PlayerID, req.GameID)
	if err != nil {
		return nil, err
	}

	// Create attempt record (immutable)
	attemptID := uuid.New().String()
	attemptQuery := `INSERT INTO attempt_records
	                 (attempt_id, player_id, game_id, topic_id, question_id, difficulty_level_id,
	                  difficulty_numeric_tier, was_correct, submitted_answer, correct_answer,
	                  time_to_answer_ms, hint_used, confidence_recovery_mode, attempted_at)
	                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err = tx.Exec(attemptQuery,
		attemptID,
		req.PlayerID,
		req.GameID,
		req.TopicID,
		req.QuestionID,
		req.DifficultyLevelID,
		difficultyTier,
		validationResult.IsCorrect,
		req.SubmittedAnswer,
		req.CorrectAnswer,
		req.TimeToAnswerMs,
		req.HintUsed,
		profile.ConfidenceRecoveryModeActive,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	// Get all attempts for this topic to assess mastery
	topicAttempts, err := o.getTopicAttempts(tx, req.PlayerID, req.GameID, req.TopicID)
	if err != nil {
		return nil, err
	}

	// Assess new mastery state
	newMasteryState := o.learningEngine.AssessMastery(topicAttempts)

	// Update or create mastery record
	err = o.updateMasteryRecord(tx, req.PlayerID, req.GameID, req.TopicID, newMasteryState, validationResult.IsCorrect)
	if err != nil {
		return nil, err
	}

	// Get all recent attempts for confidence recovery check
	recentAttempts, err := o.getAllRecentAttempts(tx, req.PlayerID, req.GameID, 7)
	if err != nil {
		return nil, err
	}

	// Check confidence recovery triggers and exits
	shouldEnterRecovery := o.learningEngine.CheckConfidenceRecoveryTrigger(recentAttempts, profile.ConfidenceRecoveryModeActive)
	shouldExitRecovery := o.learningEngine.CheckConfidenceRecoveryExit(recentAttempts, profile.ConfidenceRecoveryEnteredAt)

	// Update profile recovery mode if needed
	if shouldEnterRecovery {
		now := time.Now()
		updateQuery := `UPDATE player_game_profiles
		                SET confidence_recovery_mode_active = true, confidence_recovery_entered_at = $1
		                WHERE player_id = $2 AND game_id = $3`
		_, err = tx.Exec(updateQuery, now, req.PlayerID, req.GameID)
		if err != nil {
			return nil, err
		}
		profile.ConfidenceRecoveryModeActive = true
	} else if shouldExitRecovery {
		updateQuery := `UPDATE player_game_profiles
		                SET confidence_recovery_mode_active = false, confidence_recovery_entered_at = NULL
		                WHERE player_id = $2 AND game_id = $3`
		_, err = tx.Exec(updateQuery, req.PlayerID, req.GameID)
		if err != nil {
			return nil, err
		}
		profile.ConfidenceRecoveryModeActive = false
	}

	// Update profile stats
	statsQuery := `UPDATE player_game_profiles
	               SET total_questions_answered = total_questions_answered + 1,
	                   total_correct_answers = total_correct_answers + $1
	               WHERE player_id = $2 AND game_id = $3`
	correctIncrement := 0
	if validationResult.IsCorrect {
		correctIncrement = 1
	}
	_, err = tx.Exec(statsQuery, correctIncrement, req.PlayerID, req.GameID)
	if err != nil {
		return nil, err
	}

	// If correct, calculate and award stars
	starsEarned := 0
	if validationResult.IsCorrect {
		starsEarned = o.rewardEngine.CalculateStarReward(newMasteryState, difficultyTier, profile.ConfidenceRecoveryModeActive)

		// Create reward transaction
		transactionID := uuid.New().String()
		rewardQuery := `INSERT INTO reward_transactions
		                (transaction_id, player_id, reward_type, reward_identifier, amount, earned_at,
		                 game_id, topic_id, attempt_id, difficulty_tier, confidence_recovery_mode, consolidated)
		                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

		_, err = tx.Exec(rewardQuery,
			transactionID,
			req.PlayerID,
			models.RewardTypeStar,
			"star",
			starsEarned,
			time.Now(),
			req.GameID,
			req.TopicID,
			attemptID,
			difficultyTier,
			profile.ConfidenceRecoveryModeActive,
			false,
		)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &ProcessAnswerResponse{
		IsCorrect:          validationResult.IsCorrect,
		CorrectAnswer:      validationResult.CorrectAnswer,
		FeedbackMessage:    validationResult.FeedbackMessage,
		StarsEarned:        starsEarned,
		NewMasteryState:    newMasteryState,
		RecoveryModeActive: profile.ConfidenceRecoveryModeActive,
	}, nil
}

// Helper functions

func (o *GameOrchestrator) getMasteryRecords(playerID, gameID string) ([]*models.MasteryRecord, error) {
	db := database.GetDB()

	query := `SELECT mastery_record_id, player_id, game_id, topic_id, mastery_state, last_assessed_at,
	          times_practiced, consecutive_correct, last_correct_at, last_incorrect_at
	          FROM mastery_records WHERE player_id = $1 AND game_id = $2`

	rows, err := db.Query(query, playerID, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*models.MasteryRecord
	for rows.Next() {
		var mr models.MasteryRecord
		err := rows.Scan(&mr.MasteryRecordID, &mr.PlayerID, &mr.GameID, &mr.TopicID, &mr.MasteryState,
			&mr.LastAssessedAt, &mr.TimesPracticed, &mr.ConsecutiveCorrect, &mr.LastCorrectAt, &mr.LastIncorrectAt)
		if err != nil {
			return nil, err
		}
		records = append(records, &mr)
	}

	return records, nil
}

func (o *GameOrchestrator) initializeMasteryRecords(playerID, gameID string) ([]*models.MasteryRecord, error) {
	db := database.GetDB()

	// Get all topics for this game
	topicQuery := `SELECT topic_id FROM topics WHERE game_id = $1 ORDER BY display_order`
	rows, err := db.Query(topicQuery, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topicIDs []string
	for rows.Next() {
		var topicID string
		if err := rows.Scan(&topicID); err != nil {
			return nil, err
		}
		topicIDs = append(topicIDs, topicID)
	}

	// Create mastery record for each topic
	var records []*models.MasteryRecord
	for _, topicID := range topicIDs {
		masteryRecordID := uuid.New().String()
		insertQuery := `INSERT INTO mastery_records
		                (mastery_record_id, player_id, game_id, topic_id, mastery_state, last_assessed_at,
		                 times_practiced, consecutive_correct)
		                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := db.Exec(insertQuery, masteryRecordID, playerID, gameID, topicID,
			models.MasteryStateUnknown, time.Now(), 0, 0)
		if err != nil {
			return nil, err
		}

		records = append(records, &models.MasteryRecord{
			MasteryRecordID:    masteryRecordID,
			PlayerID:           playerID,
			GameID:             gameID,
			TopicID:            topicID,
			MasteryState:       models.MasteryStateUnknown,
			LastAssessedAt:     time.Now(),
			TimesPracticed:     0,
			ConsecutiveCorrect: 0,
		})
	}

	return records, nil
}

func (o *GameOrchestrator) getRecentTopics(playerID, gameID string, limit int) ([]string, error) {
	db := database.GetDB()

	query := `SELECT topic_id FROM attempt_records
	          WHERE player_id = $1 AND game_id = $2
	          ORDER BY attempted_at DESC LIMIT $3`

	rows, err := db.Query(query, playerID, gameID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var topicID string
		if err := rows.Scan(&topicID); err != nil {
			return nil, err
		}
		topics = append(topics, topicID)
	}

	return topics, nil
}

func (o *GameOrchestrator) getRecentAccuracy(playerID, gameID, topicID string, limit int) (float64, error) {
	db := database.GetDB()

	query := `SELECT was_correct FROM attempt_records
	          WHERE player_id = $1 AND game_id = $2 AND topic_id = $3
	          ORDER BY attempted_at DESC LIMIT $4`

	rows, err := db.Query(query, playerID, gameID, topicID, limit)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	correct := 0
	total := 0
	for rows.Next() {
		var wasCorrect bool
		if err := rows.Scan(&wasCorrect); err != nil {
			return 0, err
		}
		total++
		if wasCorrect {
			correct++
		}
	}

	if total == 0 {
		return 0.5, nil // Default to 50% if no history
	}

	return float64(correct) / float64(total), nil
}

func (o *GameOrchestrator) cacheQuestion(playerID string, question *models.Question) error {
	// For MVP, we'll skip caching. In production, use Redis or in-memory cache
	return nil
}

func (o *GameOrchestrator) getDifficultyTier(gameID, difficultyLevelID string) (int, error) {
	db := database.GetDB()

	var tier int
	query := `SELECT numeric_tier FROM difficulty_levels WHERE game_id = $1 AND level_id = $2`
	err := db.QueryRow(query, gameID, difficultyLevelID).Scan(&tier)
	if err != nil {
		return 0, err
	}

	return tier, nil
}

func (o *GameOrchestrator) getPlayerGameProfile(playerID, gameID string) (*models.PlayerGameProfile, error) {
	db := database.GetDB()

	var profile models.PlayerGameProfile
	query := `SELECT profile_id, player_id, game_id, created_at, last_played_at, total_questions_answered,
	          total_correct_answers, current_difficulty_level_id, confidence_recovery_mode_active,
	          confidence_recovery_entered_at, session_count
	          FROM player_game_profiles WHERE player_id = $1 AND game_id = $2`

	err := db.QueryRow(query, playerID, gameID).Scan(
		&profile.ProfileID, &profile.PlayerID, &profile.GameID, &profile.CreatedAt, &profile.LastPlayedAt,
		&profile.TotalQuestionsAnswered, &profile.TotalCorrectAnswers, &profile.CurrentDifficultyLevelID,
		&profile.ConfidenceRecoveryModeActive, &profile.ConfidenceRecoveryEnteredAt, &profile.SessionCount,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (o *GameOrchestrator) getTopicAttempts(tx interface{}, playerID, gameID, topicID string) ([]*models.AttemptRecord, error) {
	// Type assert to handle both *sql.DB and *sql.Tx
	type Queryable interface {
		Query(query string, args ...interface{}) (interface{}, error)
	}

	db := database.GetDB()
	query := `SELECT attempt_id, player_id, game_id, topic_id, question_id, difficulty_level_id,
	          difficulty_numeric_tier, was_correct, submitted_answer, correct_answer,
	          time_to_answer_ms, hint_used, confidence_recovery_mode, attempted_at
	          FROM attempt_records
	          WHERE player_id = $1 AND game_id = $2 AND topic_id = $3
	          ORDER BY attempted_at ASC`

	rows, err := db.Query(query, playerID, gameID, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.AttemptRecord
	for rows.Next() {
		var att models.AttemptRecord
		err := rows.Scan(&att.AttemptID, &att.PlayerID, &att.GameID, &att.TopicID, &att.QuestionID,
			&att.DifficultyLevelID, &att.DifficultyNumericTier, &att.WasCorrect, &att.SubmittedAnswer,
			&att.CorrectAnswer, &att.TimeToAnswerMs, &att.HintUsed, &att.ConfidenceRecoveryMode, &att.AttemptedAt)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, &att)
	}

	return attempts, nil
}

func (o *GameOrchestrator) getAllRecentAttempts(tx interface{}, playerID, gameID string, limit int) ([]*models.AttemptRecord, error) {
	db := database.GetDB()
	query := `SELECT attempt_id, player_id, game_id, topic_id, question_id, difficulty_level_id,
	          difficulty_numeric_tier, was_correct, submitted_answer, correct_answer,
	          time_to_answer_ms, hint_used, confidence_recovery_mode, attempted_at
	          FROM attempt_records
	          WHERE player_id = $1 AND game_id = $2
	          ORDER BY attempted_at DESC LIMIT $3`

	rows, err := db.Query(query, playerID, gameID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.AttemptRecord
	for rows.Next() {
		var att models.AttemptRecord
		err := rows.Scan(&att.AttemptID, &att.PlayerID, &att.GameID, &att.TopicID, &att.QuestionID,
			&att.DifficultyLevelID, &att.DifficultyNumericTier, &att.WasCorrect, &att.SubmittedAnswer,
			&att.CorrectAnswer, &att.TimeToAnswerMs, &att.HintUsed, &att.ConfidenceRecoveryMode, &att.AttemptedAt)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, &att)
	}

	return attempts, nil
}

func (o *GameOrchestrator) updateMasteryRecord(tx interface{}, playerID, gameID, topicID string, newState models.MasteryState, wasCorrect bool) error {
	db := database.GetDB()

	// Check if mastery record exists
	var existingID string
	checkQuery := `SELECT mastery_record_id FROM mastery_records
	               WHERE player_id = $1 AND game_id = $2 AND topic_id = $3`
	err := db.QueryRow(checkQuery, playerID, gameID, topicID).Scan(&existingID)

	now := time.Now()

	if err != nil {
		// Record doesn't exist, create it
		masteryRecordID := uuid.New().String()
		insertQuery := `INSERT INTO mastery_records
		                (mastery_record_id, player_id, game_id, topic_id, mastery_state, last_assessed_at,
		                 times_practiced, consecutive_correct, last_correct_at, last_incorrect_at)
		                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		consecutiveCorrect := 0
		var lastCorrectAt, lastIncorrectAt *time.Time
		if wasCorrect {
			consecutiveCorrect = 1
			lastCorrectAt = &now
		} else {
			lastIncorrectAt = &now
		}

		_, err = db.Exec(insertQuery, masteryRecordID, playerID, gameID, topicID, newState, now,
			1, consecutiveCorrect, lastCorrectAt, lastIncorrectAt)
		return err
	}

	// Update existing record
	updateQuery := `UPDATE mastery_records
	                SET mastery_state = $1, last_assessed_at = $2, times_practiced = times_practiced + 1`

	if wasCorrect {
		updateQuery += `, consecutive_correct = consecutive_correct + 1, last_correct_at = $3
		                WHERE mastery_record_id = $4`
		_, err = db.Exec(updateQuery, newState, now, now, existingID)
	} else {
		updateQuery += `, consecutive_correct = 0, last_incorrect_at = $3
		                WHERE mastery_record_id = $4`
		_, err = db.Exec(updateQuery, newState, now, now, existingID)
	}

	return err
}
