package starrewards

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
)

// DefaultStarRewards defines the system default star rewards per difficulty level
var DefaultStarRewards = map[string]int{
	"easy":   5,
	"medium": 10,
	"hard":   15,
}

// StarRewardService handles player-specific star reward configuration
type StarRewardService struct{}

// GetPlayerStarRewards retrieves all custom star reward configs for a player
func (s *StarRewardService) GetPlayerStarRewards(playerID string) (map[string]models.PlayerStarReward, error) {
	db := database.GetDB()

	query := `SELECT id, player_id, game_id, easy_stars, medium_stars, hard_stars, created_at, updated_at
	          FROM player_star_rewards
	          WHERE player_id = $1`

	rows, err := db.Query(query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get star rewards: %w", err)
	}
	defer rows.Close()

	rewards := make(map[string]models.PlayerStarReward)
	for rows.Next() {
		var reward models.PlayerStarReward
		err := rows.Scan(
			&reward.ID,
			&reward.PlayerID,
			&reward.GameID,
			&reward.EasyStars,
			&reward.MediumStars,
			&reward.HardStars,
			&reward.CreatedAt,
			&reward.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan star reward: %w", err)
		}
		rewards[reward.GameID] = reward
	}

	return rewards, nil
}

// GetStarRewardForPlayer returns the star reward for a specific difficulty level
// Falls back to system defaults if no custom config exists
func (s *StarRewardService) GetStarRewardForPlayer(playerID, gameID, difficultyLevelID string) int {
	db := database.GetDB()

	query := `SELECT easy_stars, medium_stars, hard_stars
	          FROM player_star_rewards
	          WHERE player_id = $1 AND game_id = $2`

	var easyStars, mediumStars, hardStars int
	err := db.QueryRow(query, playerID, gameID).Scan(&easyStars, &mediumStars, &hardStars)

	if err == sql.ErrNoRows {
		// Fall back to defaults
		if stars, ok := DefaultStarRewards[difficultyLevelID]; ok {
			return stars
		}
		return DefaultStarRewards["easy"] // Ultimate fallback
	}

	if err != nil {
		// On error, return default
		return DefaultStarRewards["easy"]
	}

	// Return player-specific config
	switch difficultyLevelID {
	case "easy":
		return easyStars
	case "medium":
		return mediumStars
	case "hard":
		return hardStars
	default:
		return DefaultStarRewards["easy"]
	}
}

// SetPlayerStarRewards upserts star reward configs for a player
func (s *StarRewardService) SetPlayerStarRewards(playerID string, request models.PlayerStarRewardsRequest) error {
	db := database.GetDB()

	// Use a transaction for atomic updates
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	upsertQuery := `INSERT INTO player_star_rewards
	                (id, player_id, game_id, easy_stars, medium_stars, hard_stars, created_at, updated_at)
	                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	                ON CONFLICT (player_id, game_id)
	                DO UPDATE SET
	                    easy_stars = EXCLUDED.easy_stars,
	                    medium_stars = EXCLUDED.medium_stars,
	                    hard_stars = EXCLUDED.hard_stars,
	                    updated_at = EXCLUDED.updated_at`

	now := time.Now()

	for gameID, config := range request.StarRewards {
		// Validate star values
		if config.EasyStars <= 0 || config.MediumStars <= 0 || config.HardStars <= 0 {
			return fmt.Errorf("star values must be positive for game %s", gameID)
		}

		// Check if record exists
		var existingID string
		checkQuery := `SELECT id FROM player_star_rewards WHERE player_id = $1 AND game_id = $2`
		err := tx.QueryRow(checkQuery, playerID, gameID).Scan(&existingID)

		if err == sql.ErrNoRows {
			// Insert new record
			id := uuid.New().String()
			_, err = tx.Exec(upsertQuery,
				id,
				playerID,
				gameID,
				config.EasyStars,
				config.MediumStars,
				config.HardStars,
				now,
				now,
			)
		} else if err == nil {
			// Update existing record
			updateQuery := `UPDATE player_star_rewards
			                SET easy_stars = $1, medium_stars = $2, hard_stars = $3, updated_at = $4
			                WHERE player_id = $5 AND game_id = $6`
			_, err = tx.Exec(updateQuery,
				config.EasyStars,
				config.MediumStars,
				config.HardStars,
				now,
				playerID,
				gameID,
			)
		}

		if err != nil {
			return fmt.Errorf("failed to set star rewards for game %s: %w", gameID, err)
		}
	}

	return tx.Commit()
}

// DeletePlayerStarReward removes custom star config for a player (reverts to defaults)
func (s *StarRewardService) DeletePlayerStarReward(playerID, gameID string) error {
	db := database.GetDB()

	query := `DELETE FROM player_star_rewards WHERE player_id = $1 AND game_id = $2`

	result, err := db.Exec(query, playerID, gameID)
	if err != nil {
		return fmt.Errorf("failed to delete star reward: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no custom star reward found for player %s and game %s", playerID, gameID)
	}

	return nil
}

// GetAllGamesWithDefaults returns a map of all games with their default star rewards
func (s *StarRewardService) GetAllGamesWithDefaults() (map[string]models.StarRewardConfig, error) {
	db := database.GetDB()

	query := `SELECT game_id FROM games WHERE is_active = true ORDER BY display_name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %w", err)
	}
	defer rows.Close()

	games := make(map[string]models.StarRewardConfig)
	for rows.Next() {
		var gameID string
		if err := rows.Scan(&gameID); err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		games[gameID] = models.StarRewardConfig{
			EasyStars:   DefaultStarRewards["easy"],
			MediumStars: DefaultStarRewards["medium"],
			HardStars:   DefaultStarRewards["hard"],
		}
	}

	return games, nil
}
