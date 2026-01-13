package platform

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
)

// CreatePlayer creates a new player
func CreatePlayer(playerName string) (*models.Player, error) {
	db := database.GetDB()

	player := &models.Player{
		PlayerID:          uuid.New().String(),
		PlayerName:        playerName,
		CreatedAt:         time.Now(),
		TotalSessionCount: 0,
	}

	query := `INSERT INTO players (player_id, player_name, created_at, total_session_count)
	          VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, player.PlayerID, player.PlayerName, player.CreatedAt, player.TotalSessionCount)
	if err != nil {
		return nil, err
	}

	return player, nil
}

// ListPlayers returns all players
func ListPlayers() ([]*models.Player, error) {
	db := database.GetDB()

	query := `SELECT player_id, player_name, created_at, last_played_at, total_session_count, preferred_character_id
	          FROM players ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*models.Player
	for rows.Next() {
		var p models.Player
		err := rows.Scan(&p.PlayerID, &p.PlayerName, &p.CreatedAt, &p.LastPlayedAt, &p.TotalSessionCount, &p.PreferredCharacterID)
		if err != nil {
			return nil, err
		}
		players = append(players, &p)
	}

	return players, nil
}

// GetPlayerByID retrieves a player by ID
func GetPlayerByID(playerID string) (*models.Player, error) {
	db := database.GetDB()

	query := `SELECT player_id, player_name, created_at, last_played_at, total_session_count, preferred_character_id
	          FROM players WHERE player_id = $1`

	var p models.Player
	err := db.QueryRow(query, playerID).Scan(&p.PlayerID, &p.PlayerName, &p.CreatedAt, &p.LastPlayedAt, &p.TotalSessionCount, &p.PreferredCharacterID)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// StartGameSession initializes or loads a PlayerGameProfile
func StartGameSession(playerID, gameID string) (*models.PlayerGameProfile, error) {
	db := database.GetDB()

	// Check if profile exists
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
		// Profile doesn't exist, create it
		profile = models.PlayerGameProfile{
			ProfileID:                    uuid.New().String(),
			PlayerID:                     playerID,
			GameID:                       gameID,
			CreatedAt:                    time.Now(),
			TotalQuestionsAnswered:       0,
			TotalCorrectAnswers:          0,
			CurrentDifficultyLevelID:     "easy",
			ConfidenceRecoveryModeActive: false,
			SessionCount:                 0,
		}

		insertQuery := `INSERT INTO player_game_profiles
		                (profile_id, player_id, game_id, created_at, total_questions_answered, total_correct_answers,
		                 current_difficulty_level_id, confidence_recovery_mode_active, session_count)
		                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err = db.Exec(insertQuery, profile.ProfileID, profile.PlayerID, profile.GameID, profile.CreatedAt,
			profile.TotalQuestionsAnswered, profile.TotalCorrectAnswers, profile.CurrentDifficultyLevelID,
			profile.ConfidenceRecoveryModeActive, profile.SessionCount)
		if err != nil {
			return nil, err
		}
	}

	// Update session count and last played
	now := time.Now()
	updateQuery := `UPDATE player_game_profiles
	                SET session_count = session_count + 1, last_played_at = $1
	                WHERE profile_id = $2`
	_, err = db.Exec(updateQuery, now, profile.ProfileID)
	if err != nil {
		return nil, err
	}

	profile.SessionCount++
	profile.LastPlayedAt = &now

	// Also update player's last_played_at
	playerUpdateQuery := `UPDATE players SET last_played_at = $1, total_session_count = total_session_count + 1 WHERE player_id = $2`
	_, err = db.Exec(playerUpdateQuery, now, playerID)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// EndGameSession updates PlayerGameProfile at session end
func EndGameSession(playerID, gameID, sessionID string) (map[string]interface{}, error) {
	db := database.GetDB()

	// Get profile stats
	var profile models.PlayerGameProfile
	query := `SELECT profile_id, total_questions_answered, total_correct_answers
	          FROM player_game_profiles WHERE player_id = $1 AND game_id = $2`

	err := db.QueryRow(query, playerID, gameID).Scan(&profile.ProfileID, &profile.TotalQuestionsAnswered, &profile.TotalCorrectAnswers)
	if err != nil {
		return nil, err
	}

	// Calculate session summary
	summary := map[string]interface{}{
		"total_questions": profile.TotalQuestionsAnswered,
		"total_correct":   profile.TotalCorrectAnswers,
		"accuracy":        0.0,
	}

	if profile.TotalQuestionsAnswered > 0 {
		summary["accuracy"] = float64(profile.TotalCorrectAnswers) / float64(profile.TotalQuestionsAnswered)
	}

	return summary, nil
}

// DeletePlayer deletes a player and all associated data
func DeletePlayer(playerID string) error {
	db := database.GetDB()

	// Start transaction to ensure all related data is deleted atomically
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete in correct order to respect foreign key constraints
	// 1. Delete reward reset events (references player_id)
	_, err = tx.Exec(`DELETE FROM reward_reset_events WHERE player_id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete reward reset events: %w", err)
	}

	// 2. Delete reward transactions (references player_id)
	_, err = tx.Exec(`DELETE FROM reward_transactions WHERE player_id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete reward transactions: %w", err)
	}

	// 3. Delete attempt records (references player_id)
	_, err = tx.Exec(`DELETE FROM attempt_records WHERE player_id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete attempt records: %w", err)
	}

	// 4. Delete mastery records (references player_id)
	_, err = tx.Exec(`DELETE FROM mastery_records WHERE player_id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete mastery records: %w", err)
	}

	// 5. Delete player game profiles (references player_id)
	_, err = tx.Exec(`DELETE FROM player_game_profiles WHERE player_id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete player game profiles: %w", err)
	}

	// 6. Delete player number ranges (references player_id)
	// Use a subquery to check if table exists first
	var tableExists bool
	err = tx.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = 'player_number_ranges'
	)`).Scan(&tableExists)
	if err == nil && tableExists {
		_, err = tx.Exec(`DELETE FROM player_number_ranges WHERE player_id = $1`, playerID)
		if err != nil {
			return fmt.Errorf("failed to delete player number ranges: %w", err)
		}
	}

	// 7. Finally delete the player
	_, err = tx.Exec(`DELETE FROM players WHERE player_id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
