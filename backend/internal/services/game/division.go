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

// DivisionGame implements the division facts game
type DivisionGame struct{}

// GenerateQuestion generates a division question
func (g *DivisionGame) GenerateQuestion(playerID, topicID, difficultyLevelID string, seed int64) (*models.Question, error) {
	rng := rand.New(rand.NewSource(seed))

	// Check if player has custom operand-specific ranges
	operandRanges, hasCustom, err := numberranges.GetOperandRangesForQuestionGeneration(playerID, "division-facts")
	if err != nil {
		hasCustom = false
	}

	// If custom range is set, use simplified logic based on the range
	if hasCustom && operandRanges != nil {
		return g.generateCustomRangeQuestion(rng, topicID, difficultyLevelID, operandRanges)
	}
	// Parse topic to get divisor (e.g., "divide-7" -> 7)
	parts := strings.Split(topicID, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid topic ID: %s", topicID)
	}

	divisor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid topic ID: %s", topicID)
	}

	// Use existing rng from custom range check

	// Generate quotient based on difficulty
	var quotient int
	switch difficultyLevelID {
	case "easy":
		// Easy: 1-5
		quotient = rng.Intn(5) + 1 // 1-5
	case "medium", "hard":
		// Medium and Hard: 1-12
		quotient = rng.Intn(12) + 1 // 1-12
	default:
		quotient = rng.Intn(12) + 1
	}

	dividend := quotient * divisor
	correctAnswer := quotient

	question := &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d ÷ %d?", dividend, divisor),
		QuestionData: map[string]interface{}{
			"dividend": dividend,
			"divisor":  divisor,
			"quotient": quotient,
		},
		VisualHintData: map[string]interface{}{
			"hint_type": "grouping",
			"groups":    quotient,
			"per_group": divisor,
			"total":     dividend,
			"color":     "green",
		},
		VerbalHint:    g.generateVerbalHint(dividend, divisor),
		CorrectAnswer: fmt.Sprintf("%d", correctAnswer),
		GeneratedAt:   time.Now(),
	}

	return question, nil
}

// generateVerbalHint creates an age-appropriate verbal hint for division
func (g *DivisionGame) generateVerbalHint(dividend, divisor int) string {
	quotient := dividend / divisor

	// For simple division, use reverse multiplication thinking
	if divisor <= 5 && quotient <= 10 {
		return g.generateReverseMultiplicationHint(dividend, divisor, quotient)
	}

	// For division by 10, 11
	if divisor == 10 {
		return fmt.Sprintf("Think: Just remove the 0 from %d!", dividend)
	}
	if divisor == 11 && quotient <= 9 {
		return fmt.Sprintf("Think: %d ÷ %d = %d (repeat digits)", dividend, divisor, quotient)
	}

	// For larger problems, use friendly numbers
	return g.generateGroupingHint(dividend, divisor, quotient)
}

// generateReverseMultiplicationHint uses the relationship between multiplication and division
func (g *DivisionGame) generateReverseMultiplicationHint(dividend, divisor, quotient int) string {
	if quotient <= 3 {
		// Show the multiplication pattern with ellipsis
		pattern := ""
		for i := 1; i <= quotient; i++ {
			if i > 1 {
				pattern += ", "
			}
			pattern += fmt.Sprintf("%d", i*divisor)
		}
		return fmt.Sprintf("Count by %ds: %s... What do you multiply %d by to get %d?",
			divisor, pattern, divisor, dividend)
	}

	if quotient == 4 {
		return fmt.Sprintf("Think: %d × ? = %d (Tip: Count by %ds: %d, %d, %d, ...)",
			divisor, dividend, divisor, divisor, divisor*2, divisor*3)
	}

	if quotient == 5 {
		return fmt.Sprintf("Think: How many %ds make %d? (Use your %s times table!)",
			divisor, dividend, g.numberToWord(divisor))
	}

	return fmt.Sprintf("Think: ? × %d = %d", divisor, dividend)
}

// generateGroupingHint helps visualize the division as grouping
func (g *DivisionGame) generateGroupingHint(dividend, divisor, quotient int) string {
	if quotient <= 10 {
		return fmt.Sprintf("Think: If you split %d into groups of %d, how many groups?", dividend, divisor)
	}

	return fmt.Sprintf("Tip: What times %d gives you %d?", divisor, dividend)
}

// numberToWord converts small numbers to words for kid-friendly hints
func (g *DivisionGame) numberToWord(num int) string {
	words := map[int]string{
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
		10: "ten",
		11: "eleven",
		12: "twelve",
	}
	if word, ok := words[num]; ok {
		return word
	}
	return strconv.Itoa(num)
}

// ValidateAnswer checks if the submitted answer is correct
func (g *DivisionGame) ValidateAnswer(correctAnswer, submittedAnswer string) (*models.AnswerValidationResult, error) {
	isCorrect := strings.TrimSpace(submittedAnswer) == correctAnswer

	result := &models.AnswerValidationResult{
		IsCorrect:       isCorrect,
		CorrectAnswer:   correctAnswer,
		FeedbackMessage: g.getFeedbackMessage(isCorrect),
	}

	return result, nil
}

func (g *DivisionGame) getFeedbackMessage(isCorrect bool) string {
	if isCorrect {
		messages := []string{"Amazing!", "Excellent!", "Great job!", "You got it!", "Fantastic!", "Brilliant!"}
		return messages[rand.Intn(len(messages))]
	}
	messages := []string{"Not quite, but close!", "Let's try another one!", "Good try!", "Almost there!"}
	return messages[rand.Intn(len(messages))]
}

// generateCustomRangeQuestion generates a question within custom operand ranges
// For division: operand1 = dividend range, operand2 = divisor range
func (g *DivisionGame) generateCustomRangeQuestion(rng *rand.Rand, topicID, difficultyLevelID string, ranges *models.OperandRanges) (*models.Question, error) {
	// Ensure we have valid ranges
	dividendMin, dividendMax := ranges.Operand1Min, ranges.Operand1Max
	divisorMin, divisorMax := ranges.Operand2Min, ranges.Operand2Max

	if dividendMax <= dividendMin {
		dividendMax = dividendMin + 1
	}
	if divisorMax <= divisorMin {
		divisorMax = divisorMin + 1
	}

	// Ensure divisor is at least 1
	if divisorMin < 1 {
		divisorMin = 1
	}

	// Generate divisor within range
	divisor := rng.Intn(divisorMax-divisorMin+1) + divisorMin
	if divisor == 0 {
		divisor = 1
	}

	// Generate quotient such that dividend is within range
	maxQuotient := dividendMax / divisor
	if maxQuotient < 1 {
		maxQuotient = 1
	}
	quotient := rng.Intn(maxQuotient) + 1

	dividend := quotient * divisor

	// Ensure dividend is within range
	if dividend < dividendMin {
		quotient = (dividendMin / divisor) + 1
		dividend = quotient * divisor
	}

	correctAnswer := quotient
	hintText := g.generateVerbalHint(dividend, divisor)

	return &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d ÷ %d?", dividend, divisor),
		QuestionData: map[string]interface{}{
			"dividend": dividend,
			"divisor":  divisor,
			"quotient": quotient,
		},
		VisualHintData: map[string]interface{}{
			"hint_type": "grouping",
			"groups":    quotient,
			"per_group": divisor,
			"total":     dividend,
			"color":     "green",
		},
		VerbalHint:    hintText,
		CorrectAnswer: fmt.Sprintf("%d", correctAnswer),
		GeneratedAt:   time.Now(),
	}, nil
}

// GetGameDefinition returns the game metadata
func (g *DivisionGame) GetGameDefinition() *models.Game {
	return &models.Game{
		GameID:      "division-facts",
		DisplayName: "Division Detective",
		Description: "Solve division puzzles and become a math detective!",
		IconName:    "division-icon.svg",
		Version:     "1.0.0",
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
}
