package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/models"
	"github.com/learning-game/backend/internal/services/numberranges"
)

// SubtractionGame implements the subtraction facts game
type SubtractionGame struct{}

// GenerateQuestion generates a subtraction question
func (g *SubtractionGame) GenerateQuestion(playerID, topicID, difficultyLevelID string, seed int64) (*models.Question, error) {
	rng := rand.New(rand.NewSource(seed))

	// Check if player has custom operand-specific ranges
	operandRanges, hasCustom, err := numberranges.GetOperandRangesForQuestionGeneration(playerID, "subtraction-facts")
	if err != nil {
		hasCustom = false
	}

	// If custom range is set, use simplified logic based on the range
	if hasCustom && operandRanges != nil {
		return g.generateCustomRangeQuestion(rng, topicID, difficultyLevelID, operandRanges)
	}

	// Topic IDs now represent skill levels
	// Format: "subtract-level-N" where N is the level
	var level int
	if len(topicID) > 14 && topicID[:14] == "subtract-level" {
		_, err := fmt.Sscanf(topicID, "subtract-level-%d", &level)
		if err != nil {
			return nil, fmt.Errorf("invalid topic ID: %s", topicID)
		}
	} else {
		return nil, fmt.Errorf("invalid topic ID format: %s (expected subtract-level-N)", topicID)
	}

	var minuend, subtrahend int
	var hintText string

	// Generate numbers based on skill level
	switch level {
	case 1:
		// Level 1: Single digit, subtract small numbers (e.g., 8-2, 7-3)
		minuend = rng.Intn(8) + 2    // 2-9
		subtrahend = rng.Intn(minuend-1) + 1  // 1 to minuend-1
		hintText = g.generateCountBackHint(minuend, subtrahend)

	case 2:
		// Level 2: Single digit, result 0-5 (e.g., 9-4, 7-5)
		subtrahend = rng.Intn(5) + 2  // 2-6
		maxMinuend := subtrahend + 5
		if maxMinuend > 9 { maxMinuend = 9 }
		minuend = rng.Intn(maxMinuend-subtrahend) + subtrahend + 1
		hintText = g.generateCountBackHint(minuend, subtrahend)

	case 3:
		// Level 3: Double digit, no borrowing (e.g., 86-42, 57-23)
		ones1 := rng.Intn(9) + 1    // 1-9
		ones2 := rng.Intn(ones1)    // 0 to ones1-1 (no borrowing needed)
		tens1 := rng.Intn(8) + 2    // 2-9
		tens2 := rng.Intn(tens1)    // 0 to tens1-1
		minuend = tens1*10 + ones1
		subtrahend = tens2*10 + ones2
		hintText = fmt.Sprintf("Think: Subtract ones: %d-%d=%d, then tens: %d-%d=%d. Answer: %d",
			ones1, ones2, ones1-ones2, tens1, tens2, tens1-tens2, minuend-subtrahend)

	case 4:
		// Level 4: Double digit, with borrowing (e.g., 52-37, 81-26)
		subtrahend = rng.Intn(89) + 10  // 10-99
		minuend = rng.Intn(89) + 10     // 10-99
		if minuend < subtrahend {
			minuend, subtrahend = subtrahend, minuend
		}
		hintText = fmt.Sprintf("Think: If the top number is smaller, borrow from the tens place!")

	case 5:
		// Level 5: Triple digit, no borrowing in hundreds (e.g., 456-123)
		ones1 := rng.Intn(9) + 1
		ones2 := rng.Intn(ones1 + 1)
		tens1 := rng.Intn(9) + 1
		tens2 := rng.Intn(tens1 + 1)
		hundreds1 := rng.Intn(8) + 2    // 2-9
		hundreds2 := rng.Intn(hundreds1) // 0 to hundreds1-1
		minuend = hundreds1*100 + tens1*10 + ones1
		subtrahend = hundreds2*100 + tens2*10 + ones2
		hintText = fmt.Sprintf("Think: Subtract each place value: ones, then tens, then hundreds.")

	case 6:
		// Level 6: Triple digit with borrowing (e.g., 654-278, 921-357)
		subtrahend = rng.Intn(899) + 100  // 100-999
		minuend = rng.Intn(899) + 100     // 100-999
		if minuend < subtrahend {
			minuend, subtrahend = subtrahend, minuend
		}
		hintText = fmt.Sprintf("Think: Subtract carefully, borrowing from the next place when needed.")

	case 7:
		// Level 7: Any 3-digit numbers, more complex borrowing
		subtrahend = rng.Intn(899) + 100  // 100-999
		minuend = rng.Intn(899) + 100     // 100-999
		if minuend < subtrahend {
			minuend, subtrahend = subtrahend, minuend
		}
		hintText = fmt.Sprintf("Think: Take your time subtracting each place, borrowing when the top is smaller.")

	default:
		// Default to level 1
		minuend = rng.Intn(8) + 2
		subtrahend = rng.Intn(minuend-1) + 1
		hintText = g.generateCountBackHint(minuend, subtrahend)
	}

	difference := minuend - subtrahend

	question := &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d − %d?", minuend, subtrahend),
		QuestionData: map[string]interface{}{
			"minuend":    minuend,
			"subtrahend": subtrahend,
			"difference": difference,
			"level":      level,
		},
		VisualHintData: map[string]interface{}{
			"hint_type":  "take-away",
			"minuend":    minuend,
			"subtrahend": subtrahend,
			"difference": difference,
			"color":      "orange",
			"level":      level,
		},
		VerbalHint:    hintText,
		CorrectAnswer: fmt.Sprintf("%d", difference),
		GeneratedAt:   time.Now(),
	}

	return question, nil
}

// generateCountBackHint provides hints using counting back
func (g *SubtractionGame) generateCountBackHint(minuend, subtrahend int) string {
	if subtrahend == 0 {
		return fmt.Sprintf("Think: Taking away zero changes nothing. %d − 0 = %d!", minuend, minuend)
	}
	if subtrahend == 1 {
		return fmt.Sprintf("Think: Just count back one from %d. What comes before?", minuend)
	}
	if subtrahend == 2 {
		return fmt.Sprintf("Think: Count back two from %d. %d, %d, then %d!",
			minuend, minuend, minuend-1, minuend-2)
	}
	if subtrahend == 3 {
		return fmt.Sprintf("Think: Count back three from %d. %d, %d, %d, then %d!",
			minuend, minuend, minuend-1, minuend-2, minuend-3)
	}

	return fmt.Sprintf("Think: Start at %d and count back %d times.", minuend, subtrahend)
}

// ValidateAnswer checks if the submitted answer is correct
func (g *SubtractionGame) ValidateAnswer(correctAnswer, submittedAnswer string) (*models.AnswerValidationResult, error) {
	isCorrect := submittedAnswer == correctAnswer

	result := &models.AnswerValidationResult{
		IsCorrect:       isCorrect,
		CorrectAnswer:   correctAnswer,
		FeedbackMessage: g.getFeedbackMessage(isCorrect),
	}

	return result, nil
}

func (g *SubtractionGame) getFeedbackMessage(isCorrect bool) string {
	if isCorrect {
		messages := []string{"Amazing!", "Excellent!", "Great job!", "You got it!", "Fantastic!", "Brilliant!"}
		return messages[rand.Intn(len(messages))]
	}
	messages := []string{"Not quite, but close!", "Let's try another one!", "Good try!", "Almost there!"}
	return messages[rand.Intn(len(messages))]
}

// generateCustomRangeQuestion generates a question within custom operand ranges
// For subtraction: operand1 = minuend range, operand2 = subtrahend range
func (g *SubtractionGame) generateCustomRangeQuestion(rng *rand.Rand, topicID, difficultyLevelID string, ranges *models.OperandRanges) (*models.Question, error) {
	// Ensure we have valid ranges
	minuendMin, minuendMax := ranges.Operand1Min, ranges.Operand1Max
	subtrahendMin, subtrahendMax := ranges.Operand2Min, ranges.Operand2Max

	if minuendMax <= minuendMin {
		minuendMax = minuendMin + 1
	}
	if subtrahendMax <= subtrahendMin {
		subtrahendMax = subtrahendMin + 1
	}

	// Generate minuend within range
	minuend := rng.Intn(minuendMax-minuendMin+1) + minuendMin

	// Generate subtrahend within range, but ensure it's <= minuend for positive results
	maxSubtrahend := subtrahendMax
	if maxSubtrahend > minuend {
		maxSubtrahend = minuend
	}
	if maxSubtrahend < subtrahendMin {
		// If minuend is smaller than subtrahend min, swap to ensure positive result
		minuend, maxSubtrahend = minuendMax, minuend
	}

	subtrahend := rng.Intn(maxSubtrahend-subtrahendMin+1) + subtrahendMin
	if subtrahend > minuend {
		subtrahend = minuend
	}

	difference := minuend - subtrahend
	level := 1 // Default level for custom ranges
	hintText := fmt.Sprintf("Think: Start at %d and count back %d times.", minuend, subtrahend)

	return &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d − %d?", minuend, subtrahend),
		QuestionData: map[string]interface{}{
			"minuend":    minuend,
			"subtrahend": subtrahend,
			"difference": difference,
			"level":      level,
		},
		VisualHintData: map[string]interface{}{
			"hint_type":  "take-away",
			"minuend":    minuend,
			"subtrahend": subtrahend,
			"difference": difference,
			"color":      "orange",
			"level":      level,
		},
		VerbalHint:    hintText,
		CorrectAnswer: fmt.Sprintf("%d", difference),
		GeneratedAt:   time.Now(),
	}, nil
}

// GetGameDefinition returns the game metadata
func (g *SubtractionGame) GetGameDefinition() *models.Game {
	return &models.Game{
		GameID:      "subtraction-facts",
		DisplayName: "Subtraction Safari",
		Description: "Explore subtraction and become a math ranger!",
		IconName:    "subtraction-icon.svg",
		Version:     "1.0.0",
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
}
