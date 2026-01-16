package models

import "time"

// Player represents a child using the platform
type Player struct {
	PlayerID            string    `json:"player_id" db:"player_id"`
	PlayerName          string    `json:"player_name" db:"player_name"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	LastPlayedAt        *time.Time `json:"last_played_at,omitempty" db:"last_played_at"`
	TotalSessionCount   int       `json:"total_session_count" db:"total_session_count"`
	PreferredCharacterID *string   `json:"preferred_character_id,omitempty" db:"preferred_character_id"`
}

// Game represents a learning game
type Game struct {
	GameID      string    `json:"game_id" db:"game_id"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Description string    `json:"description" db:"description"`
	IconName    string    `json:"icon_name" db:"icon_name"`
	Version     string    `json:"version" db:"version"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Topic represents a learnable sub-unit within a game
type Topic struct {
	TopicID              string   `json:"topic_id" db:"topic_id"`
	GameID               string   `json:"game_id" db:"game_id"`
	DisplayName          string   `json:"display_name" db:"display_name"`
	DisplayOrder         int      `json:"display_order" db:"display_order"`
	PrerequisiteTopicIDs []string `json:"prerequisite_topic_ids" db:"prerequisite_topic_ids"`
	EstimatedPoolSize    int      `json:"estimated_pool_size" db:"estimated_pool_size"`
}

// DifficultyLevel defines difficulty tiers within a game
type DifficultyLevel struct {
	LevelID     string `json:"level_id" db:"level_id"`
	GameID      string `json:"game_id" db:"game_id"`
	DisplayName string `json:"display_name" db:"display_name"`
	NumericTier int    `json:"numeric_tier" db:"numeric_tier"`
	Description string `json:"description" db:"description"`
}

// PlayerGameProfile tracks player state within a specific game
type PlayerGameProfile struct {
	ProfileID                    string     `json:"profile_id" db:"profile_id"`
	PlayerID                     string     `json:"player_id" db:"player_id"`
	GameID                       string     `json:"game_id" db:"game_id"`
	CreatedAt                    time.Time  `json:"created_at" db:"created_at"`
	LastPlayedAt                 *time.Time `json:"last_played_at,omitempty" db:"last_played_at"`
	TotalQuestionsAnswered       int        `json:"total_questions_answered" db:"total_questions_answered"`
	TotalCorrectAnswers          int        `json:"total_correct_answers" db:"total_correct_answers"`
	CurrentDifficultyLevelID     string     `json:"current_difficulty_level_id" db:"current_difficulty_level_id"`
	ConfidenceRecoveryModeActive bool       `json:"confidence_recovery_mode_active" db:"confidence_recovery_mode_active"`
	ConfidenceRecoveryEnteredAt  *time.Time `json:"confidence_recovery_entered_at,omitempty" db:"confidence_recovery_entered_at"`
	SessionCount                 int        `json:"session_count" db:"session_count"`
}

// MasteryState enum values
type MasteryState string

const (
	MasteryStateUnknown  MasteryState = "UNKNOWN"
	MasteryStateWeak     MasteryState = "WEAK"
	MasteryStateLearning MasteryState = "LEARNING"
	MasteryStateStrong   MasteryState = "STRONG"
	MasteryStateMastered MasteryState = "MASTERED"
)

// MasteryRecord tracks mastery state for a specific topic
type MasteryRecord struct {
	MasteryRecordID   string        `json:"mastery_record_id" db:"mastery_record_id"`
	PlayerID          string        `json:"player_id" db:"player_id"`
	GameID            string        `json:"game_id" db:"game_id"`
	TopicID           string        `json:"topic_id" db:"topic_id"`
	MasteryState      MasteryState  `json:"mastery_state" db:"mastery_state"`
	LastAssessedAt    time.Time     `json:"last_assessed_at" db:"last_assessed_at"`
	TimesPracticed    int           `json:"times_practiced" db:"times_practiced"`
	ConsecutiveCorrect int          `json:"consecutive_correct" db:"consecutive_correct"`
	LastCorrectAt     *time.Time    `json:"last_correct_at,omitempty" db:"last_correct_at"`
	LastIncorrectAt   *time.Time    `json:"last_incorrect_at,omitempty" db:"last_incorrect_at"`
}

// AttemptRecord is an immutable log of every question attempt
type AttemptRecord struct {
	AttemptID              string    `json:"attempt_id" db:"attempt_id"`
	PlayerID               string    `json:"player_id" db:"player_id"`
	GameID                 string    `json:"game_id" db:"game_id"`
	TopicID                string    `json:"topic_id" db:"topic_id"`
	QuestionID             string    `json:"question_id" db:"question_id"`
	DifficultyLevelID      string    `json:"difficulty_level_id" db:"difficulty_level_id"`
	DifficultyNumericTier  int       `json:"difficulty_numeric_tier" db:"difficulty_numeric_tier"`
	WasCorrect             bool      `json:"was_correct" db:"was_correct"`
	SubmittedAnswer        string    `json:"submitted_answer" db:"submitted_answer"`
	CorrectAnswer          string    `json:"correct_answer" db:"correct_answer"`
	TimeToAnswerMs         *int      `json:"time_to_answer_ms,omitempty" db:"time_to_answer_ms"`
	HintUsed               bool      `json:"hint_used" db:"hint_used"`
	ConfidenceRecoveryMode bool      `json:"confidence_recovery_mode" db:"confidence_recovery_mode"`
	AttemptedAt            time.Time `json:"attempted_at" db:"attempted_at"`
}

// RewardType enum values
type RewardType string

const (
	RewardTypeStar            RewardType = "STAR"
	RewardTypeSticker         RewardType = "STICKER"
	RewardTypeBadge           RewardType = "BADGE"
	RewardTypeCharacterUnlock RewardType = "CHARACTER_UNLOCK"
)

// RewardTransaction is an immutable ledger of all rewards earned
type RewardTransaction struct {
	TransactionID          string     `json:"transaction_id" db:"transaction_id"`
	PlayerID               string     `json:"player_id" db:"player_id"`
	RewardType             RewardType `json:"reward_type" db:"reward_type"`
	RewardIdentifier       string     `json:"reward_identifier" db:"reward_identifier"`
	Amount                 int        `json:"amount" db:"amount"`
	EarnedAt               time.Time  `json:"earned_at" db:"earned_at"`
	GameID                 *string    `json:"game_id,omitempty" db:"game_id"`
	TopicID                *string    `json:"topic_id,omitempty" db:"topic_id"`
	AttemptID              *string    `json:"attempt_id,omitempty" db:"attempt_id"`
	DifficultyTier         *int       `json:"difficulty_tier,omitempty" db:"difficulty_tier"`
	ConfidenceRecoveryMode bool       `json:"confidence_recovery_mode" db:"confidence_recovery_mode"`
	Consolidated           bool       `json:"consolidated" db:"consolidated"`
	Notes                  *string    `json:"notes,omitempty" db:"notes"`
}

// RewardResetEvent is an audit trail for reward consolidation
type RewardResetEvent struct {
	ResetID                    string    `json:"reset_id" db:"reset_id"`
	PlayerID                   string    `json:"player_id" db:"player_id"`
	ResetBy                    string    `json:"reset_by" db:"reset_by"`
	ResetAt                    time.Time `json:"reset_at" db:"reset_at"`
	StarsConsolidated          int       `json:"stars_consolidated" db:"stars_consolidated"`
	RealWorldReward            string    `json:"real_world_reward" db:"real_world_reward"`
	TransactionIDsConsolidated []string  `json:"transaction_ids_consolidated" db:"transaction_ids_consolidated"`
}

// Question represents a generated question (not persisted, generated on-the-fly)
type Question struct {
	QuestionID        string                 `json:"question_id"`
	TopicID           string                 `json:"topic_id"`
	DifficultyLevelID string                 `json:"difficulty_level_id"`
	QuestionText      string                 `json:"question_text"`
	QuestionData      map[string]interface{} `json:"question_data"`
	VisualHintData    map[string]interface{} `json:"visual_hint_data,omitempty"`
	VerbalHint        string                 `json:"verbal_hint,omitempty"`
	CorrectAnswer     string                 `json:"-"` // Not sent to frontend
	GeneratedAt       time.Time              `json:"generated_at"`
}

// AnswerValidationResult represents the result of validating an answer
type AnswerValidationResult struct {
	IsCorrect       bool   `json:"is_correct"`
	CorrectAnswer   string `json:"correct_answer"`
	FeedbackMessage string `json:"feedback_message"`
}

// PlayerNumberRange stores custom min/max number ranges per player per game
type PlayerNumberRange struct {
	ID          string    `json:"id" db:"id"`
	PlayerID    string    `json:"player_id" db:"player_id"`
	GameID      string    `json:"game_id" db:"game_id"`
	MinNumber   int       `json:"min_number" db:"min_number"`
	MaxNumber   int       `json:"max_number" db:"max_number"`
	Operand1Min *int      `json:"operand1_min,omitempty" db:"operand1_min"`
	Operand1Max *int      `json:"operand1_max,omitempty" db:"operand1_max"`
	Operand2Min *int      `json:"operand2_min,omitempty" db:"operand2_min"`
	Operand2Max *int      `json:"operand2_max,omitempty" db:"operand2_max"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PlayerNumberRangesRequest is used for updating number ranges for a player
type PlayerNumberRangesRequest struct {
	// Map of game_id -> NumberRangeConfig
	NumberRanges map[string]NumberRangeConfig `json:"number_ranges"`
}

// NumberRangeConfig represents min/max for a single game with optional operand-specific ranges
type NumberRangeConfig struct {
	MinNumber   int  `json:"min_number"`
	MaxNumber   int  `json:"max_number"`
	Operand1Min *int `json:"operand1_min,omitempty"`
	Operand1Max *int `json:"operand1_max,omitempty"`
	Operand2Min *int `json:"operand2_min,omitempty"`
	Operand2Max *int `json:"operand2_max,omitempty"`
}

// OperandRanges provides a convenient structure for game engines to get ranges
type OperandRanges struct {
	Operand1Min int
	Operand1Max int
	Operand2Min int
	Operand2Max int
}

// TopicDifficultyProgress tracks difficulty progression per topic per player
type TopicDifficultyProgress struct {
	ID                       string    `json:"id" db:"id"`
	PlayerID                 string    `json:"player_id" db:"player_id"`
	GameID                   string    `json:"game_id" db:"game_id"`
	TopicID                  string    `json:"topic_id" db:"topic_id"`
	CurrentDifficultyLevelID string    `json:"current_difficulty_level_id" db:"current_difficulty_level_id"`
	ConsecutiveCorrect       int       `json:"consecutive_correct" db:"consecutive_correct"`
	ConsecutiveWrong         int       `json:"consecutive_wrong" db:"consecutive_wrong"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
}

// PlayerStarReward stores custom star rewards per difficulty level per player per game
type PlayerStarReward struct {
	ID          string    `json:"id" db:"id"`
	PlayerID    string    `json:"player_id" db:"player_id"`
	GameID      string    `json:"game_id" db:"game_id"`
	EasyStars   int       `json:"easy_stars" db:"easy_stars"`
	MediumStars int       `json:"medium_stars" db:"medium_stars"`
	HardStars   int       `json:"hard_stars" db:"hard_stars"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PlayerStarRewardsRequest is used for updating star rewards for a player
type PlayerStarRewardsRequest struct {
	// Map of game_id -> StarRewardConfig
	StarRewards map[string]StarRewardConfig `json:"star_rewards"`
}

// StarRewardConfig represents star rewards for a single game
type StarRewardConfig struct {
	EasyStars   int `json:"easy_stars"`
	MediumStars int `json:"medium_stars"`
	HardStars   int `json:"hard_stars"`
}
