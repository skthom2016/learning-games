package rewards

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
)

// RewardEngine handles reward calculation and distribution
type RewardEngine struct{}

// CalculateStarReward calculates stars earned for a correct answer
// Based on ARCHITECTURE.md Section 5 - Star Earning Matrix
//
// Matrix:
// | Mastery State | Tier 1 | Tier 2 | Tier 3 | Tier 4+ |
// |---------------|--------|--------|--------|---------|
// | UNKNOWN       |   5    |   10   |   15   |   20    |
// | WEAK          |   5    |   10   |   15   |   20    |
// | LEARNING      |   3    |   10   |   15   |   20    |
// | STRONG        |   2    |    8   |   12   |   18    |
// | MASTERED      |   1    |    5   |    8   |   12    |
func (e *RewardEngine) CalculateStarReward(masteryState models.MasteryState, difficultyTier int, inRecoveryMode bool) int {
	starMatrix := map[models.MasteryState][]int{
		models.MasteryStateUnknown:  {5, 10, 15, 20},
		models.MasteryStateWeak:     {5, 10, 15, 20},
		models.MasteryStateLearning: {3, 10, 15, 20},
		models.MasteryStateStrong:   {2, 8, 12, 18},
		models.MasteryStateMastered: {1, 5, 8, 12},
	}

	tierIndex := difficultyTier - 1
	if tierIndex < 0 {
		tierIndex = 0
	}
	if tierIndex >= 4 {
		tierIndex = 3
	}

	if stars, ok := starMatrix[masteryState]; ok {
		return stars[tierIndex]
	}

	return 5 // Default
}

// CreateRewardTransaction creates a reward transaction record
func (e *RewardEngine) CreateRewardTransaction(playerID, gameID, topicID, attemptID string, stars int, difficultyTier int, inRecoveryMode bool) error {
	db := database.GetDB()

	transaction := &models.RewardTransaction{
		TransactionID:          uuid.New().String(),
		PlayerID:               playerID,
		RewardType:             models.RewardTypeStar,
		RewardIdentifier:       "star",
		Amount:                 stars,
		EarnedAt:               time.Now(),
		GameID:                 &gameID,
		TopicID:                &topicID,
		AttemptID:              &attemptID,
		DifficultyTier:         &difficultyTier,
		ConfidenceRecoveryMode: inRecoveryMode,
		Consolidated:           false,
	}

	query := `INSERT INTO reward_transactions
	          (transaction_id, player_id, reward_type, reward_identifier, amount, earned_at,
	           game_id, topic_id, attempt_id, difficulty_tier, confidence_recovery_mode, consolidated)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := db.Exec(query,
		transaction.TransactionID,
		transaction.PlayerID,
		transaction.RewardType,
		transaction.RewardIdentifier,
		transaction.Amount,
		transaction.EarnedAt,
		transaction.GameID,
		transaction.TopicID,
		transaction.AttemptID,
		transaction.DifficultyTier,
		transaction.ConfidenceRecoveryMode,
		transaction.Consolidated,
	)

	return err
}

// GetRewardBalance calculates current reward balance for a player
func (e *RewardEngine) GetRewardBalance(playerID string) (total int, unconsolidated int, err error) {
	db := database.GetDB()

	query := `SELECT
	          COALESCE(SUM(amount) FILTER (WHERE reward_type = 'STAR'), 0) as total_stars,
	          COALESCE(SUM(amount) FILTER (WHERE reward_type = 'STAR' AND consolidated = false), 0) as unconsolidated_stars
	          FROM reward_transactions WHERE player_id = $1`

	err = db.QueryRow(query, playerID).Scan(&total, &unconsolidated)
	if err != nil {
		return 0, 0, err
	}

	return total, unconsolidated, nil
}

// ConsolidateRewards marks rewards as consolidated (converted to real-world reward)
// See ARCHITECTURE.md Section 5 - Reward Consolidation
func (e *RewardEngine) ConsolidateRewards(playerID string, starsToConsolidate int, realWorldReward string) error {
	db := database.GetDB()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Get unconsolidated balance
	var unconsolidatedBalance int
	balanceQuery := `SELECT COALESCE(SUM(amount), 0)
	                 FROM reward_transactions
	                 WHERE player_id = $1 AND reward_type = 'STAR' AND consolidated = false`
	err = tx.QueryRow(balanceQuery, playerID).Scan(&unconsolidatedBalance)
	if err != nil {
		return err
	}

	// 2. Validate sufficient balance
	if starsToConsolidate > unconsolidatedBalance {
		return fmt.Errorf("insufficient balance: want %d, have %d", starsToConsolidate, unconsolidatedBalance)
	}

	// 3. Get transaction IDs to consolidate (oldest first)
	transactionIDs := []string{}
	runningTotal := 0

	transQuery := `SELECT transaction_id, amount
	               FROM reward_transactions
	               WHERE player_id = $1 AND reward_type = 'STAR' AND consolidated = false
	               ORDER BY earned_at ASC`
	rows, err := tx.Query(transQuery, playerID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var transID string
		var amount int
		if err := rows.Scan(&transID, &amount); err != nil {
			rows.Close()
			return err
		}

		transactionIDs = append(transactionIDs, transID)
		runningTotal += amount

		if runningTotal >= starsToConsolidate {
			break
		}
	}

	// Explicitly close rows before continuing with transaction
	if err := rows.Close(); err != nil {
		return err
	}

	// 4. Create RewardResetEvent
	resetID := uuid.New().String()
	eventQuery := `INSERT INTO reward_reset_events
	               (reset_id, player_id, reset_by, stars_consolidated, real_world_reward, reset_at)
	               VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(eventQuery, resetID, playerID, "admin", starsToConsolidate, realWorldReward, time.Now())
	if err != nil {
		return err
	}

	// 5. Mark transactions as consolidated
	for _, transID := range transactionIDs {
		updateQuery := `UPDATE reward_transactions
		                SET consolidated = true
		                WHERE transaction_id = $1`
		_, err = tx.Exec(updateQuery, transID)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}
