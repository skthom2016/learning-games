package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/models"
	"github.com/learning-game/backend/internal/services/numberranges"
)

// GetPlayerNumberRanges handles GET /api/players/:player_id/number-ranges
// Returns all number range configurations for a player
func GetPlayerNumberRanges(c *gin.Context) {
	playerID := c.Param("player_id")

	ranges, err := numberranges.GetPlayerNumberRanges(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number_ranges": ranges,
	})
}

// SetPlayerNumberRanges handles PUT /api/players/:player_id/number-ranges
// Creates or updates number range configurations for a player
func SetPlayerNumberRanges(c *gin.Context) {
	playerID := c.Param("player_id")

	var req models.PlayerNumberRangesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate ranges for each game
	for gameID, config := range req.NumberRanges {
		// Validate global min < max
		if config.MinNumber >= config.MaxNumber {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "For game " + gameID + ": min_number must be less than max_number",
			})
			return
		}

		// Validate operand1 range if set
		if config.Operand1Min != nil && config.Operand1Max != nil {
			if *config.Operand1Min >= *config.Operand1Max {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "For game " + gameID + ": operand1_min must be less than operand1_max",
				})
				return
			}
		}

		// Validate operand2 range if set
		if config.Operand2Min != nil && config.Operand2Max != nil {
			if *config.Operand2Min >= *config.Operand2Max {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "For game " + gameID + ": operand2_min must be less than operand2_max",
				})
				return
			}
		}
	}

	// Set each range using the new function that supports operand-specific ranges
	var results []*models.PlayerNumberRange
	for gameID, config := range req.NumberRanges {
		rangeConfig, err := numberranges.SetPlayerNumberRangeWithOperands(playerID, gameID, config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, rangeConfig)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Number ranges updated successfully",
		"number_ranges": results,
	})
}

// DeletePlayerNumberRange handles DELETE /api/players/:player_id/number-ranges/:game_id
// Removes a custom number range for a specific game
func DeletePlayerNumberRange(c *gin.Context) {
	playerID := c.Param("player_id")
	gameID := c.Param("game_id")

	err := numberranges.DeletePlayerNumberRange(playerID, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Number range deleted successfully",
	})
}
