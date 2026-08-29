package numberranges

import (
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/database"
	"github.com/learning-game/backend/internal/models"
)

// GetPlayerNumberRanges retrieves all number ranges for a player
func GetPlayerNumberRanges(playerID string) ([]*models.PlayerNumberRange, error) {
	db := database.GetDB()

	query := `SELECT id, player_id, game_id, min_number, max_number,
	          operand1_min, operand1_max, operand2_min, operand2_max,
	          created_at, updated_at
	          FROM player_number_ranges WHERE player_id = $1 ORDER BY game_id`

	rows, err := db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranges []*models.PlayerNumberRange
	for rows.Next() {
		var r models.PlayerNumberRange
		err := rows.Scan(&r.ID, &r.PlayerID, &r.GameID, &r.MinNumber, &r.MaxNumber,
			&r.Operand1Min, &r.Operand1Max, &r.Operand2Min, &r.Operand2Max,
			&r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, &r)
	}

	return ranges, nil
}

// GetPlayerNumberRange retrieves a specific number range for a player and game
func GetPlayerNumberRange(playerID, gameID string) (*models.PlayerNumberRange, error) {
	db := database.GetDB()

	query := `SELECT id, player_id, game_id, min_number, max_number,
	          operand1_min, operand1_max, operand2_min, operand2_max,
	          created_at, updated_at
	          FROM player_number_ranges WHERE player_id = $1 AND game_id = $2`

	var r models.PlayerNumberRange
	err := db.QueryRow(query, playerID, gameID).Scan(
		&r.ID, &r.PlayerID, &r.GameID, &r.MinNumber, &r.MaxNumber,
		&r.Operand1Min, &r.Operand1Max, &r.Operand2Min, &r.Operand2Max,
		&r.CreatedAt, &r.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

// SetPlayerNumberRangeWithOperands creates or updates a number range with operand-specific ranges
func SetPlayerNumberRangeWithOperands(playerID, gameID string, config models.NumberRangeConfig) (*models.PlayerNumberRange, error) {
	db := database.GetDB()
	now := time.Now()

	// Create new range object
	newRange := &models.PlayerNumberRange{
		ID:          uuid.New().String(),
		PlayerID:    playerID,
		GameID:      gameID,
		MinNumber:   config.MinNumber,
		MaxNumber:   config.MaxNumber,
		Operand1Min: config.Operand1Min,
		Operand1Max: config.Operand1Max,
		Operand2Min: config.Operand2Min,
		Operand2Max: config.Operand2Max,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `INSERT INTO player_number_ranges
	          (id, player_id, game_id, min_number, max_number,
	           operand1_min, operand1_max, operand2_min, operand2_max,
	           created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	          ON CONFLICT (player_id, game_id) DO UPDATE
	          SET min_number = EXCLUDED.min_number,
	              max_number = EXCLUDED.max_number,
	              operand1_min = EXCLUDED.operand1_min,
	              operand1_max = EXCLUDED.operand1_max,
	              operand2_min = EXCLUDED.operand2_min,
	              operand2_max = EXCLUDED.operand2_max,
	              updated_at = EXCLUDED.updated_at`

	_, err := db.Exec(query, newRange.ID, newRange.PlayerID, newRange.GameID,
		newRange.MinNumber, newRange.MaxNumber,
		newRange.Operand1Min, newRange.Operand1Max, newRange.Operand2Min, newRange.Operand2Max,
		newRange.CreatedAt, newRange.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return newRange, nil
}

// SetPlayerNumberRange creates or updates a number range for a player and game (legacy)
func SetPlayerNumberRange(playerID, gameID string, minNumber, maxNumber int) (*models.PlayerNumberRange, error) {
	config := models.NumberRangeConfig{
		MinNumber: minNumber,
		MaxNumber: maxNumber,
	}
	return SetPlayerNumberRangeWithOperands(playerID, gameID, config)
}

// DeletePlayerNumberRange removes a custom number range for a player and game
func DeletePlayerNumberRange(playerID, gameID string) error {
	db := database.GetDB()

	query := `DELETE FROM player_number_ranges WHERE player_id = $1 AND game_id = $2`
	_, err := db.Exec(query, playerID, gameID)
	return err
}

// GetNumberRangeForQuestionGeneration returns the min/max numbers to use for question generation.
// If no custom range is set, returns default values.
func GetNumberRangeForQuestionGeneration(playerID, gameID string) (min, max int, hasCustom bool, err error) {
	rangeConfig, err := GetPlayerNumberRange(playerID, gameID)
	if err != nil {
		// No custom range set, return defaults
		return 0, 0, false, nil
	}

	return rangeConfig.MinNumber, rangeConfig.MaxNumber, true, nil
}

// GetOperandRangesForQuestionGeneration returns separate ranges for each operand.
// If operand-specific ranges are not set, falls back to min_number/max_number.
// If no custom range is set at all, returns (nil, false, nil).
func GetOperandRangesForQuestionGeneration(playerID, gameID string) (*models.OperandRanges, bool, error) {
	rangeConfig, err := GetPlayerNumberRange(playerID, gameID)
	if err != nil {
		// No custom range set
		return nil, false, nil
	}

	ranges := &models.OperandRanges{}

	// Use operand-specific ranges if set, otherwise fall back to min_number/max_number
	if rangeConfig.Operand1Min != nil {
		ranges.Operand1Min = *rangeConfig.Operand1Min
	} else {
		ranges.Operand1Min = rangeConfig.MinNumber
	}

	if rangeConfig.Operand1Max != nil {
		ranges.Operand1Max = *rangeConfig.Operand1Max
	} else {
		ranges.Operand1Max = rangeConfig.MaxNumber
	}

	if rangeConfig.Operand2Min != nil {
		ranges.Operand2Min = *rangeConfig.Operand2Min
	} else {
		ranges.Operand2Min = rangeConfig.MinNumber
	}

	if rangeConfig.Operand2Max != nil {
		ranges.Operand2Max = *rangeConfig.Operand2Max
	} else {
		ranges.Operand2Max = rangeConfig.MaxNumber
	}

	return ranges, true, nil
}
