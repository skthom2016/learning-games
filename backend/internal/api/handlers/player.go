package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/services/platform"
)

// CreatePlayer handles POST /api/players
// TODO: Implement player creation logic
func CreatePlayer(c *gin.Context) {
	var req struct {
		PlayerName string `json:"player_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Call platform service to create player
	player, err := platform.CreatePlayer(req.PlayerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, player)
}

// ListPlayers handles GET /api/players
// TODO: Implement listing all players
func ListPlayers(c *gin.Context) {
	// TODO: Call platform service to list all players
	players, err := platform.ListPlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

// GetPlayer handles GET /api/players/:player_id
// TODO: Implement getting single player by ID
func GetPlayer(c *gin.Context) {
	playerID := c.Param("player_id")

	// TODO: Call platform service to get player by ID
	player, err := platform.GetPlayerByID(playerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

// DeletePlayer handles DELETE /api/players/:player_id
func DeletePlayer(c *gin.Context) {
	playerID := c.Param("player_id")

	err := platform.DeletePlayer(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}
