package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/models"
	"github.com/learning-game/backend/internal/services/numberranges"
)

// MultiplicationGame implements the multiplication tables game
type MultiplicationGame struct{}

// GenerateQuestion generates a multiplication question
func (g *MultiplicationGame) GenerateQuestion(playerID, topicID, difficultyLevelID string, seed int64) (*models.Question, error) {
	// Check if player has custom operand-specific ranges
	operandRanges, hasCustom, err := numberranges.GetOperandRangesForQuestionGeneration(playerID, "multiplication-tables")
	if err != nil {
		hasCustom = false
	}

	// If custom range is set, use simplified logic based on the range
	if hasCustom && operandRanges != nil {
		rng := rand.New(rand.NewSource(seed))
		return g.generateCustomRangeQuestion(rng, topicID, difficultyLevelID, operandRanges)
	}
	// Parse topic to get multiplier (e.g., "times-7" -> 7)
	parts := strings.Split(topicID, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid topic ID: %s", topicID)
	}

	multiplier, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid topic ID: %s", topicID)
	}

	rng := rand.New(rand.NewSource(seed))

	// Generate operand1 based on difficulty
	var operand1 int
	switch difficultyLevelID {
	case "easy":
		// Easy: 2-5
		operand1 = rng.Intn(4) + 2 // 2-5
	case "medium", "hard":
		// Medium and Hard: 2-12
		operand1 = rng.Intn(11) + 2 // 2-12
	default:
		operand1 = rng.Intn(11) + 2
	}

	operand2 := multiplier
	correctAnswer := operand1 * operand2

	question := &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d × %d?", operand1, operand2),
		QuestionData: map[string]interface{}{
			"operand1": operand1,
			"operand2": operand2,
		},
		VisualHintData: map[string]interface{}{
			"hint_type": "dot-array",
			"rows":      operand1,
			"cols":      operand2,
			"color":     "blue",
		},
		VerbalHint: g.generateVerbalHint(operand1, operand2),
		CorrectAnswer: fmt.Sprintf("%d", correctAnswer),
		GeneratedAt:   time.Now(),
	}

	return question, nil
}

// generateVerbalHint creates an age-appropriate verbal hint for multiplication
func (g *MultiplicationGame) generateVerbalHint(operand1, operand2 int) string {
	// Ensure operand1 <= operand2 for simpler logic
	if operand1 > operand2 {
		operand1, operand2 = operand2, operand1
	}

	// For easy problems, use skip counting
	if operand1 <= 5 && operand2 <= 10 {
		return g.generateSkipCountHint(operand1, operand2)
	}

	// For problems with 10, 11, or easier patterns
	if operand2 == 10 {
		return fmt.Sprintf("Think: Just add a 0 to %d!", operand1)
	}
	if operand2 == 11 {
		if operand1 <= 9 {
			return fmt.Sprintf("Think: %d repeated twice = %d%d", operand1, operand1, operand1)
		}
	}

	// For larger problems, use distributive property
	return g.generateDistributiveHint(operand1, operand2)
}

// generateSkipCountHint creates a skip-counting hint (doesn't reveal the answer)
func (g *MultiplicationGame) generateSkipCountHint(operand1, operand2 int) string {
	// Only show a pattern that leads to the answer, not the final answer
	// Show them the approach, not the result

	if operand1 <= 3 {
		// For small multipliers, show the pattern with ellipsis
		pattern := ""
		for i := 1; i <= operand1; i++ {
			if i > 0 {
				pattern += ", "
			}
			pattern += fmt.Sprintf("%d", i * operand2)
		}
		return fmt.Sprintf("Count by %ds: %s...", operand2, pattern)
	}

	// For larger multipliers, describe the approach without showing the full sequence
	if operand1 == 4 {
		return fmt.Sprintf("Think: Count by %ds four times (start: %d, %d, ...)", operand2, operand2, operand2*2)
	}
	if operand1 == 5 {
		return fmt.Sprintf("Think: Count by %ds five times (you can do it!) ", operand2)
	}

	// For larger, suggest a strategy
	return fmt.Sprintf("Tip: Count by %ds, %d times", operand2, operand1)
}

// generateDistributiveHint breaks down the problem using known facts
func (g *MultiplicationGame) generateDistributiveHint(operand1, operand2 int) string {
	// Find a friendly number to break down
	hints := []string{}

	// Try breaking into 10s
	if operand2 > 10 {
		remaining := operand2 - 10
		part1 := operand1 * 10
		part2 := operand1 * remaining
		hints = append(hints, fmt.Sprintf("Think: %d × %d = (%d × 10) + (%d × %d) = %d + %d",
			operand1, operand2, operand1, operand1, remaining, part1, part2))
	}

	// Try breaking using doubles
	if operand1 == 2 || operand2 == 2 {
		bigger := operand1
		if operand2 > operand1 {
			bigger = operand2
		}
		hints = append(hints, fmt.Sprintf("Think: %d doubled = %d + %d", bigger, bigger, bigger))
	}

	// Try breaking using distributive property with 5
	if operand1 > 5 && operand1 <= 10 {
		part1 := operand1 - 5
		total := operand2 * 5
		hints = append(hints, fmt.Sprintf("Think: %d × %d = (5 × %d) + (%d × %d) = %d + %d",
			operand1, operand2, operand2, part1, operand2, total, part1*operand2))
	}

	if len(hints) > 0 {
		return hints[0]
	}

	// Fallback: suggest breaking it down
	return fmt.Sprintf("Tip: Break it into smaller parts you know!")
}

// ValidateAnswer checks if the submitted answer is correct
func (g *MultiplicationGame) ValidateAnswer(correctAnswer, submittedAnswer string) (*models.AnswerValidationResult, error) {
	isCorrect := strings.TrimSpace(submittedAnswer) == correctAnswer

	result := &models.AnswerValidationResult{
		IsCorrect:       isCorrect,
		CorrectAnswer:   correctAnswer,
		FeedbackMessage: g.getFeedbackMessage(isCorrect),
	}

	return result, nil
}

func (g *MultiplicationGame) getFeedbackMessage(isCorrect bool) string {
	if isCorrect {
		messages := []string{"Amazing!", "Excellent!", "Great job!", "You got it!", "Fantastic!", "Brilliant!"}
		return messages[rand.Intn(len(messages))]
	}
	messages := []string{"Not quite, but close!", "Let's try another one!", "Good try!", "Almost there!"}
	return messages[rand.Intn(len(messages))]
}

// generateCustomRangeQuestion generates a question within custom operand ranges
// operandRanges contains separate min/max for each factor
func (g *MultiplicationGame) generateCustomRangeQuestion(rng *rand.Rand, topicID, difficultyLevelID string, ranges *models.OperandRanges) (*models.Question, error) {
	// Ensure we have valid ranges
	op1Min, op1Max := ranges.Operand1Min, ranges.Operand1Max
	op2Min, op2Max := ranges.Operand2Min, ranges.Operand2Max

	if op1Max <= op1Min {
		op1Max = op1Min + 1
	}
	if op2Max <= op2Min {
		op2Max = op2Min + 1
	}

	// Generate operands within their respective custom ranges
	operand1 := rng.Intn(op1Max-op1Min+1) + op1Min
	operand2 := rng.Intn(op2Max-op2Min+1) + op2Min

	correctAnswer := operand1 * operand2

	return &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d × %d?", operand1, operand2),
		QuestionData: map[string]interface{}{
			"operand1": operand1,
			"operand2": operand2,
			"product":  correctAnswer,
		},
		VisualHintData: map[string]interface{}{
			"hint_type": "dot-array",
			"rows":      operand1,
			"cols":      operand2,
			"color":     "blue",
		},
		VerbalHint:    g.generateVerbalHint(operand1, operand2),
		CorrectAnswer: fmt.Sprintf("%d", correctAnswer),
		GeneratedAt:   time.Now(),
	}, nil
}

// GetGameDefinition returns the game metadata
func (g *MultiplicationGame) GetGameDefinition() *models.Game {
	return &models.Game{
		GameID:      "multiplication-tables",
		DisplayName: "Multiplication Master",
		Description: "Practice times tables with fun characters!",
		IconName:    "multiplication-icon.svg",
		Version:     "1.0.0",
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
}
