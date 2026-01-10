package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/services/platform"
)

// StartSession handles POST /api/sessions/start
func StartSession(c *gin.Context) {
	var req struct {
		PlayerID string `json:"player_id" binding:"required"`
		GameID   string `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := platform.StartGameSession(req.PlayerID, req.GameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// EndSession handles POST /api/sessions/end
func EndSession(c *gin.Context) {
	var req struct {
		PlayerID  string `json:"player_id" binding:"required"`
		GameID    string `json:"game_id" binding:"required"`
		SessionID string `json:"session_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := platform.EndGameSession(req.PlayerID, req.GameID, req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
