package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learning-game/backend/internal/services/orchestrator"
)

var gameOrchestrator = orchestrator.NewGameOrchestrator()

// GetNextQuestion handles GET /api/questions/next
func GetNextQuestion(c *gin.Context) {
	playerID := c.Query("player_id")
	gameID := c.Query("game_id")

	if playerID == "" || gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player_id and game_id are required"})
		return
	}

	question, err := gameOrchestrator.GetNextQuestion(playerID, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

// SubmitAnswer handles POST /api/answers/submit
func SubmitAnswer(c *gin.Context) {
	var req struct {
		PlayerID          string `json:"player_id" binding:"required"`
		GameID            string `json:"game_id" binding:"required"`
		QuestionID        string `json:"question_id" binding:"required"`
		TopicID           string `json:"topic_id" binding:"required"`
		DifficultyLevelID string `json:"difficulty_level_id" binding:"required"`
		CorrectAnswer     string `json:"correct_answer" binding:"required"`
		SubmittedAnswer   string `json:"submitted_answer" binding:"required"`
		HintUsed          bool   `json:"hint_used"`
		TimeToAnswerMs    *int   `json:"time_to_answer_ms"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := gameOrchestrator.ProcessAnswer(orchestrator.ProcessAnswerRequest{
		PlayerID:          req.PlayerID,
		GameID:            req.GameID,
		QuestionID:        req.QuestionID,
		TopicID:           req.TopicID,
		DifficultyLevelID: req.DifficultyLevelID,
		SubmittedAnswer:   req.SubmittedAnswer,
		CorrectAnswer:     req.CorrectAnswer,
		TimeToAnswerMs:    req.TimeToAnswerMs,
		HintUsed:          req.HintUsed,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
