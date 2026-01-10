package api

import (
	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/api/handlers"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api")
	{
		// Player Management
		// TODO: Implement player CRUD operations
		v1.POST("/players", handlers.CreatePlayer)
		v1.GET("/players", handlers.ListPlayers)
		v1.GET("/players/:id", handlers.GetPlayer)

		// Session Management
		// TODO: Implement session start/end logic
		v1.POST("/sessions/start", handlers.StartSession)
		v1.POST("/sessions/end", handlers.EndSession)

		// Game Play
		// TODO: Implement question generation and answer validation
		v1.GET("/questions/next", handlers.GetNextQuestion)
		v1.POST("/answers/submit", handlers.SubmitAnswer)

		// Progress & Analytics
		// TODO: Implement progress tracking and mastery heatmap
		v1.GET("/progress/:player_id/:game_id", handlers.GetPlayerProgress)
		v1.GET("/mastery/:player_id/:game_id", handlers.GetMasteryHeatmap)

		// Rewards
		// TODO: Implement reward balance and transaction queries
		v1.GET("/rewards/:player_id", handlers.GetRewardBalance)
		v1.GET("/rewards/:player_id/transactions", handlers.GetRewardTransactions)

		// Admin Routes
		admin := v1.Group("/admin")
		{
			// TODO: Implement admin reward reset functionality
			admin.POST("/rewards/reset", handlers.ResetRewards)
			admin.GET("/analytics/:player_id", handlers.GetPlayerAnalytics)
		}
	}
}
