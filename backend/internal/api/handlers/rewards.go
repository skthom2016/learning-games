package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
	"github.com/learning-game/backend/internal/services/rewards"
)

var rewardEngine = &rewards.RewardEngine{}

// GetRewardBalance handles GET /api/rewards/:player_id
func GetRewardBalance(c *gin.Context) {
	playerID := c.Param("player_id")

	total, unconsolidated, err := rewardEngine.GetRewardBalance(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id":            playerID,
		"total_stars":          total,
		"unconsolidated_stars": unconsolidated,
	})
}

// GetRewardTransactions handles GET /api/rewards/:player_id/transactions
func GetRewardTransactions(c *gin.Context) {
	playerID := c.Param("player_id")

	db := database.GetDB()

	query := `SELECT transaction_id, player_id, reward_type, reward_identifier, amount, earned_at,
	          game_id, topic_id, attempt_id, difficulty_tier, confidence_recovery_mode, consolidated, notes
	          FROM reward_transactions WHERE player_id = $1
	          ORDER BY earned_at DESC LIMIT 100`

	rows, err := db.Query(query, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var transactions []models.RewardTransaction
	for rows.Next() {
		var rt models.RewardTransaction
		err := rows.Scan(&rt.TransactionID, &rt.PlayerID, &rt.RewardType, &rt.RewardIdentifier, &rt.Amount,
			&rt.EarnedAt, &rt.GameID, &rt.TopicID, &rt.AttemptID, &rt.DifficultyTier,
			&rt.ConfidenceRecoveryMode, &rt.Consolidated, &rt.Notes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, rt)
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id":    playerID,
		"transactions": transactions,
	})
}

// ResetRewards handles POST /api/admin/rewards/reset
func ResetRewards(c *gin.Context) {
	var req struct {
		PlayerID            string `json:"player_id" binding:"required"`
		StarsToConsolidate  int    `json:"stars_to_consolidate" binding:"required,min=1"`
		RealWorldReward     string `json:"real_world_reward" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := rewardEngine.ConsolidateRewards(req.PlayerID, req.StarsToConsolidate, req.RealWorldReward)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get new balance
	total, unconsolidated, err := rewardEngine.GetRewardBalance(req.PlayerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stars_consolidated":       req.StarsToConsolidate,
		"new_total_balance":        total,
		"new_unconsolidated_balance": unconsolidated,
		"real_world_reward":        req.RealWorldReward,
	})
}
