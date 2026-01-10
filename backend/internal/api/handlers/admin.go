package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/database"
)

// GetPlayerAnalytics handles GET /api/admin/analytics/:player_id
func GetPlayerAnalytics(c *gin.Context) {
	playerID := c.Param("player_id")

	db := database.GetDB()

	// Get all game profiles for this player
	profileQuery := `SELECT game_id, total_questions_answered, total_correct_answers, session_count,
	                 confidence_recovery_mode_active
	                 FROM player_game_profiles WHERE player_id = $1`

	rows, err := db.Query(profileQuery, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var games []gin.H
	totalQuestions := 0
	totalCorrect := 0

	for rows.Next() {
		var gameID string
		var questions, correct, sessions int
		var recoveryActive bool

		err := rows.Scan(&gameID, &questions, &correct, &sessions, &recoveryActive)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalQuestions += questions
		totalCorrect += correct

		accuracy := 0.0
		if questions > 0 {
			accuracy = float64(correct) / float64(questions)
		}

		games = append(games, gin.H{
			"game_id":         gameID,
			"total_questions": questions,
			"total_correct":   correct,
			"accuracy":        accuracy,
			"sessions":        sessions,
		})
	}

	overallAccuracy := 0.0
	if totalQuestions > 0 {
		overallAccuracy = float64(totalCorrect) / float64(totalQuestions)
	}

	// Get reward consolidation history
	consolidationQuery := `SELECT stars_consolidated, real_world_reward, reset_at
	                       FROM reward_reset_events WHERE player_id = $1
	                       ORDER BY reset_at DESC LIMIT 10`

	consolidationRows, err := db.Query(consolidationQuery, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer consolidationRows.Close()

	var consolidations []gin.H
	for consolidationRows.Next() {
		var stars int
		var reward string
		var resetAt interface{}

		err := consolidationRows.Scan(&stars, &reward, &resetAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		consolidations = append(consolidations, gin.H{
			"stars_consolidated": stars,
			"real_world_reward":  reward,
			"reset_at":           resetAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id":         playerID,
		"games":             games,
		"overall_accuracy":  overallAccuracy,
		"total_questions":   totalQuestions,
		"total_correct":     totalCorrect,
		"consolidations":    consolidations,
	})
}
