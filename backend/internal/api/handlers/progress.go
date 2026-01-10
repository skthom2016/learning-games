package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
)

// GetPlayerProgress handles GET /api/progress/:player_id/:game_id
func GetPlayerProgress(c *gin.Context) {
	playerID := c.Param("player_id")
	gameID := c.Param("game_id")

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
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	accuracy := 0.0
	if profile.TotalQuestionsAnswered > 0 {
		accuracy = float64(profile.TotalCorrectAnswers) / float64(profile.TotalQuestionsAnswered)
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id":                   profile.PlayerID,
		"game_id":                     profile.GameID,
		"total_questions":             profile.TotalQuestionsAnswered,
		"total_correct":               profile.TotalCorrectAnswers,
		"accuracy":                    accuracy,
		"session_count":               profile.SessionCount,
		"current_difficulty":          profile.CurrentDifficultyLevelID,
		"confidence_recovery_active":  profile.ConfidenceRecoveryModeActive,
	})
}

// GetMasteryHeatmap handles GET /api/mastery/:player_id/:game_id
func GetMasteryHeatmap(c *gin.Context) {
	playerID := c.Param("player_id")
	gameID := c.Param("game_id")

	db := database.GetDB()

	// First get all topics for this game
	topicsQuery := `SELECT topic_id, display_name, display_order
	                FROM topics WHERE game_id = $1
	                ORDER BY display_order`

	topicsRows, err := db.Query(topicsQuery, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer topicsRows.Close()

	// Build a map of topic_id -> display_name
	topicNames := make(map[string]string)
	var topicIDs []string
	for topicsRows.Next() {
		var topicID, displayName string
		var displayOrder int
		err := topicsRows.Scan(&topicID, &displayName, &displayOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		topicNames[topicID] = displayName
		topicIDs = append(topicIDs, topicID)
	}

	// Now get mastery records
	query := `SELECT mastery_record_id, player_id, game_id, topic_id, mastery_state, last_assessed_at,
	          times_practiced, consecutive_correct, last_correct_at, last_incorrect_at
	          FROM mastery_records WHERE player_id = $1 AND game_id = $2
	          ORDER BY topic_id`

	rows, err := db.Query(query, playerID, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Build a map of topic_id -> mastery record
	masteryMap := make(map[string]*models.MasteryRecord)
	for rows.Next() {
		var mr models.MasteryRecord
		err := rows.Scan(&mr.MasteryRecordID, &mr.PlayerID, &mr.GameID, &mr.TopicID, &mr.MasteryState,
			&mr.LastAssessedAt, &mr.TimesPracticed, &mr.ConsecutiveCorrect, &mr.LastCorrectAt, &mr.LastIncorrectAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		masteryMap[mr.TopicID] = &mr
	}

	// Build response with all topics, filling in UNKNOWN for topics without mastery records
	var topics []gin.H
	for _, topicID := range topicIDs {
		mr, hasMastery := masteryMap[topicID]
		masteryState := models.MasteryStateUnknown
		timesPracticed := 0
		consecutiveCorrect := 0
		var lastAssessedAt interface{} = nil

		if hasMastery {
			masteryState = mr.MasteryState
			timesPracticed = mr.TimesPracticed
			consecutiveCorrect = mr.ConsecutiveCorrect
			lastAssessedAt = mr.LastAssessedAt
		}

		topics = append(topics, gin.H{
			"topic_id":           topicID,
			"display_name":       topicNames[topicID],
			"mastery_state":      masteryState,
			"times_practiced":    timesPracticed,
			"consecutive_correct": consecutiveCorrect,
			"last_assessed_at":   lastAssessedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"game_id":   gameID,
		"topics":    topics,
	})
}
