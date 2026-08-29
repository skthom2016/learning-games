package progress

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
)

const (
	// LevelUpThreshold is the number of consecutive correct answers needed to level up
	LevelUpThreshold = 10
	// LevelDownThreshold is the number of consecutive wrong answers needed to level down
	LevelDownThreshold = 2
)

// DifficultyProgressService handles difficulty progression logic
type DifficultyProgressService struct{}

// DifficultyLevelOrder defines the progression order of difficulty levels
var DifficultyLevelOrder = []string{"easy", "medium", "hard"}

// GetOrCreateProgress retrieves existing progress or creates a new one starting at easy
func (s *DifficultyProgressService) GetOrCreateProgress(playerID, gameID, topicID string) (*models.TopicDifficultyProgress, error) {
	db := database.GetDB()

	// Try to get existing progress
	progress := &models.TopicDifficultyProgress{}
	query := `SELECT id, player_id, game_id, topic_id, current_difficulty_level_id,
	                 consecutive_correct, consecutive_wrong, updated_at
	          FROM topic_difficulty_progress
	          WHERE player_id = $1 AND game_id = $2 AND topic_id = $3`

	err := db.QueryRow(query, playerID, gameID, topicID).Scan(
		&progress.ID,
		&progress.PlayerID,
		&progress.GameID,
		&progress.TopicID,
		&progress.CurrentDifficultyLevelID,
		&progress.ConsecutiveCorrect,
		&progress.ConsecutiveWrong,
		&progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create new progress record starting at easy
		return s.createProgress(playerID, gameID, topicID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	return progress, nil
}

// createProgress creates a new progress record starting at easy difficulty
func (s *DifficultyProgressService) createProgress(playerID, gameID, topicID string) (*models.TopicDifficultyProgress, error) {
	db := database.GetDB()

	id := uuid.New().String()
	now := time.Now()

	query := `INSERT INTO topic_difficulty_progress
	          (id, player_id, game_id, topic_id, current_difficulty_level_id,
	           consecutive_correct, consecutive_wrong, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          RETURNING id, player_id, game_id, topic_id, current_difficulty_level_id,
	                    consecutive_correct, consecutive_wrong, updated_at`

	progress := &models.TopicDifficultyProgress{
		ID:                       id,
		PlayerID:                 playerID,
		GameID:                   gameID,
		TopicID:                  topicID,
		CurrentDifficultyLevelID: "easy",
		ConsecutiveCorrect:       0,
		ConsecutiveWrong:         0,
		UpdatedAt:                now,
	}

	err := db.QueryRow(query,
		progress.ID,
		progress.PlayerID,
		progress.GameID,
		progress.TopicID,
		progress.CurrentDifficultyLevelID,
		progress.ConsecutiveCorrect,
		progress.ConsecutiveWrong,
		progress.UpdatedAt,
	).Scan(
		&progress.ID,
		&progress.PlayerID,
		&progress.GameID,
		&progress.TopicID,
		&progress.CurrentDifficultyLevelID,
		&progress.ConsecutiveCorrect,
		&progress.ConsecutiveWrong,
		&progress.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create progress: %w", err)
	}

	return progress, nil
}

// GetCurrentDifficultyLevel returns the current difficulty level for a topic
func (s *DifficultyProgressService) GetCurrentDifficultyLevel(playerID, gameID, topicID string) (string, error) {
	progress, err := s.GetOrCreateProgress(playerID, gameID, topicID)
	if err != nil {
		return "", err
	}
	return progress.CurrentDifficultyLevelID, nil
}

// UpdateProgressAfterAnswer updates consecutive counts and determines level change
// Returns (levelChanged, newLevel, error)
func (s *DifficultyProgressService) UpdateProgressAfterAnswer(playerID, gameID, topicID string, wasCorrect bool) (bool, string, error) {
	progress, err := s.GetOrCreateProgress(playerID, gameID, topicID)
	if err != nil {
		return false, "", err
	}

	oldLevel := progress.CurrentDifficultyLevelID

	if wasCorrect {
		progress.ConsecutiveCorrect++
		progress.ConsecutiveWrong = 0

		// Check if we should level up
		if progress.ConsecutiveCorrect >= LevelUpThreshold {
			if progress.CurrentDifficultyLevelID == "easy" {
				progress.CurrentDifficultyLevelID = "medium"
				progress.ConsecutiveCorrect = 0 // Reset after level up
			} else if progress.CurrentDifficultyLevelID == "medium" {
				progress.CurrentDifficultyLevelID = "hard"
				progress.ConsecutiveCorrect = 0 // Reset after level up
			}
			// If already at hard, stay at hard
		}
	} else {
		progress.ConsecutiveWrong++
		progress.ConsecutiveCorrect = 0

		// Check if we should level down (only if not at easy)
		if progress.ConsecutiveWrong >= LevelDownThreshold && progress.CurrentDifficultyLevelID != "easy" {
			if progress.CurrentDifficultyLevelID == "hard" {
				progress.CurrentDifficultyLevelID = "medium"
				progress.ConsecutiveWrong = 0 // Reset after level down
			} else if progress.CurrentDifficultyLevelID == "medium" {
				progress.CurrentDifficultyLevelID = "easy"
				progress.ConsecutiveWrong = 0 // Reset after level down
			}
		}
	}

	// Save updated progress
	err = s.saveProgress(progress)
	if err != nil {
		return false, "", err
	}

	levelChanged := progress.CurrentDifficultyLevelID != oldLevel
	return levelChanged, progress.CurrentDifficultyLevelID, nil
}

// saveProgress saves the updated progress to the database
func (s *DifficultyProgressService) saveProgress(progress *models.TopicDifficultyProgress) error {
	db := database.GetDB()

	query := `UPDATE topic_difficulty_progress
	          SET current_difficulty_level_id = $1,
	              consecutive_correct = $2,
	              consecutive_wrong = $3,
	              updated_at = $4
	          WHERE id = $5`

	progress.UpdatedAt = time.Now()

	_, err := db.Exec(query,
		progress.CurrentDifficultyLevelID,
		progress.ConsecutiveCorrect,
		progress.ConsecutiveWrong,
		progress.UpdatedAt,
		progress.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to save progress: %w", err)
	}

	return nil
}

// InitializeProgressForTopics initializes progress records for all topics in a game
func (s *DifficultyProgressService) InitializeProgressForTopics(playerID, gameID string, topicIDs []string) error {
	for _, topicID := range topicIDs {
		_, err := s.GetOrCreateProgress(playerID, gameID, topicID)
		if err != nil {
			return fmt.Errorf("failed to initialize progress for topic %s: %w", topicID, err)
		}
	}
	return nil
}

// GetAllProgressForPlayer returns all progress records for a player in a game
func (s *DifficultyProgressService) GetAllProgressForPlayer(playerID, gameID string) ([]models.TopicDifficultyProgress, error) {
	db := database.GetDB()

	query := `SELECT id, player_id, game_id, topic_id, current_difficulty_level_id,
	                 consecutive_correct, consecutive_wrong, updated_at
	          FROM topic_difficulty_progress
	          WHERE player_id = $1 AND game_id = $2
	          ORDER BY topic_id`

	rows, err := db.Query(query, playerID, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}
	defer rows.Close()

	var progressList []models.TopicDifficultyProgress
	for rows.Next() {
		var progress models.TopicDifficultyProgress
		err := rows.Scan(
			&progress.ID,
			&progress.PlayerID,
			&progress.GameID,
			&progress.TopicID,
			&progress.CurrentDifficultyLevelID,
			&progress.ConsecutiveCorrect,
			&progress.ConsecutiveWrong,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}
		progressList = append(progressList, progress)
	}

	return progressList, nil
}
