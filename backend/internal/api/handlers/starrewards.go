package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/models"
	"github.com/learning-game/backend/internal/services/starrewards"
)

var starRewardService = &starrewards.StarRewardService{}

// GetPlayerStarRewards handles GET /api/players/:player_id/star-rewards
// Returns all star reward configurations for a player
func GetPlayerStarRewards(c *gin.Context) {
	playerID := c.Param("player_id")

	// Get custom configs
	customRewards, err := starRewardService.GetPlayerStarRewards(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get all games with defaults
	allGames, err := starRewardService.GetAllGamesWithDefaults()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Merge: use custom configs where available, defaults otherwise
	result := make(map[string]models.StarRewardConfig)
	for gameID, defaultConfig := range allGames {
		if customReward, hasCustom := customRewards[gameID]; hasCustom {
			result[gameID] = models.StarRewardConfig{
				EasyStars:   customReward.EasyStars,
				MediumStars: customReward.MediumStars,
				HardStars:   customReward.HardStars,
			}
		} else {
			result[gameID] = defaultConfig
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"star_rewards": result,
	})
}

// SetPlayerStarRewards handles PUT /api/players/:player_id/star-rewards
// Creates or updates star reward configurations for a player
func SetPlayerStarRewards(c *gin.Context) {
	playerID := c.Param("player_id")

	var req models.PlayerStarRewardsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate star values for each game
	for gameID, config := range req.StarRewards {
		if config.EasyStars <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "For game " + gameID + ": easy_stars must be positive",
			})
			return
		}
		if config.MediumStars <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "For game " + gameID + ": medium_stars must be positive",
			})
			return
		}
		if config.HardStars <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "For game " + gameID + ": hard_stars must be positive",
			})
			return
		}
	}

	err := starRewardService.SetPlayerStarRewards(playerID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Star rewards updated successfully",
	})
}

// DeletePlayerStarReward handles DELETE /api/players/:player_id/star-rewards/:game_id
// Removes a custom star reward for a specific game (reverts to defaults)
func DeletePlayerStarReward(c *gin.Context) {
	playerID := c.Param("player_id")
	gameID := c.Param("game_id")

	err := starRewardService.DeletePlayerStarReward(playerID, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Star reward deleted successfully",
	})
}
