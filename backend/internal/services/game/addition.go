package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/learning-game/backend/internal/models"
	"github.com/learning-game/backend/internal/services/numberranges"
)

// AdditionGame implements the addition facts game
type AdditionGame struct{}

// GenerateQuestion generates an addition question
func (g *AdditionGame) GenerateQuestion(playerID, topicID, difficultyLevelID string, seed int64) (*models.Question, error) {
	rng := rand.New(rand.NewSource(seed))

	// Check if player has custom operand-specific ranges
	operandRanges, hasCustom, err := numberranges.GetOperandRangesForQuestionGeneration(playerID, "addition-facts")
	if err != nil {
		hasCustom = false
	}

	// If custom range is set, use simplified logic based on the range
	if hasCustom && operandRanges != nil {
		return g.generateCustomRangeQuestion(rng, topicID, difficultyLevelID, operandRanges)
	}

	// Topic IDs now represent skill levels, not addends
	// Format: "add-level-N" where N is the level
	var level int
	if len(topicID) > 9 && topicID[:9] == "add-level" {
		_, err := fmt.Sscanf(topicID, "add-level-%d", &level)
		if err != nil {
			return nil, fmt.Errorf("invalid topic ID: %s", topicID)
		}
	} else {
		return nil, fmt.Errorf("invalid topic ID format: %s (expected add-level-N)", topicID)
	}

	var addend1, addend2 int
	var hintText string

	// Generate numbers based on skill level
	switch level {
	case 1:
		// Level 1: Single digit, sums up to 10 (e.g., 3+2, 4+5)
		addend2 = rng.Intn(5) + 1 // 1-5
		addend1 = rng.Intn(10-addend2-1) + 1
		hintText = g.generateSingleDigitHint(addend1, addend2)

	case 2:
		// Level 2: Single digit, sums 10-18 (e.g., 7+6, 8+5)
		addend2 = rng.Intn(9) + 1  // 1-9
		minAddend1 := (10 - addend2)
		if minAddend1 < 1 { minAddend1 = 1 }
		maxAddend1 := (18 - addend2)
		if maxAddend1 > 9 { maxAddend1 = 9 }
		addend1 = rng.Intn(maxAddend1-minAddend1+1) + minAddend1
		hintText = g.generateSingleDigitHint(addend1, addend2)

	case 3:
		// Level 3: Double digit, no carry (e.g., 23+15, 41+32)
		// Both addends 10-99, sum of ones place < 10, sum of tens place < 100
		ones1 := rng.Intn(8) + 1  // 1-8
		ones2 := rng.Intn(9 - ones1)  // so ones1+ones2 < 10
		tens1 := rng.Intn(8) + 1  // 1-8
		tens2 := rng.Intn(9 - tens1)  // so tens1+tens2 < 9
		addend1 = tens1*10 + ones1
		addend2 = tens2*10 + ones2
		hintText = fmt.Sprintf("Think: Add the ones: %d+%d=%d, then the tens: %d+%d=%d. Answer: %d",
			ones1, ones2, ones1+ones2, tens1, tens2, tens1+tens2, addend1+addend2)

	case 4:
		// Level 4: Double digit with carry (e.g., 47+25, 68+19)
		addend1 = rng.Intn(89) + 10  // 10-99
		addend2 = rng.Intn(89) + 10  // 10-99
		hintText = fmt.Sprintf("Think: Add from right to left. Remember to carry!")

	case 5:
		// Level 5: Triple digit, no carry in hundreds (e.g., 123+456)
		ones1 := rng.Intn(8) + 1
		ones2 := rng.Intn(9 - ones1)
		tens1 := rng.Intn(8) + 1
		tens2 := rng.Intn(9 - tens1)
		hundreds1 := rng.Intn(8) + 1
		hundreds2 := rng.Intn(9 - hundreds1)
		addend1 = hundreds1*100 + tens1*10 + ones1
		addend2 = hundreds2*100 + tens2*10 + ones2
		hintText = fmt.Sprintf("Think: Add ones, then tens, then hundreds place by place.")

	case 6:
		// Level 6: Triple digit with carries (e.g., 567+789)
		addend1 = rng.Intn(899) + 100  // 100-999
		addend2 = rng.Intn(899) + 100  // 100-999
		hintText = fmt.Sprintf("Think: Add carefully, carrying over when each place reaches 10 or more.")

	case 7:
		// Level 7: Any 3-digit numbers, including larger sums
		addend1 = rng.Intn(899) + 100  // 100-999
		addend2 = rng.Intn(899) + 100  // 100-999
		hintText = fmt.Sprintf("Think: Take your time adding each place value, carrying as needed.")

	default:
		// Default to level 1
		addend2 = rng.Intn(5) + 1
		addend1 = rng.Intn(10-addend2-1) + 1
		hintText = g.generateSingleDigitHint(addend1, addend2)
	}

	sum := addend1 + addend2

	question := &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d + %d?", addend1, addend2),
		QuestionData: map[string]interface{}{
			"addend1": addend1,
			"addend2": addend2,
			"sum":     sum,
			"level":   level,
		},
		VisualHintData: map[string]interface{}{
			"hint_type": "counting",
			"addend1":   addend1,
			"addend2":   addend2,
			"sum":       sum,
			"color":     "blue",
			"level":     level,
		},
		VerbalHint:    hintText,
		CorrectAnswer: fmt.Sprintf("%d", sum),
		GeneratedAt:   time.Now(),
	}

	return question, nil
}

// generateCustomRangeQuestion generates a question within custom operand ranges
// For addition: operand1 = addend1 range, operand2 = addend2 range
func (g *AdditionGame) generateCustomRangeQuestion(rng *rand.Rand, topicID, difficultyLevelID string, ranges *models.OperandRanges) (*models.Question, error) {
	// Ensure we have valid ranges
	addend1Min, addend1Max := ranges.Operand1Min, ranges.Operand1Max
	addend2Min, addend2Max := ranges.Operand2Min, ranges.Operand2Max

	if addend1Max <= addend1Min {
		addend1Max = addend1Min + 1
	}
	if addend2Max <= addend2Min {
		addend2Max = addend2Min + 1
	}

	// Generate addends within their respective custom ranges
	addend1 := rng.Intn(addend1Max-addend1Min+1) + addend1Min
	addend2 := rng.Intn(addend2Max-addend2Min+1) + addend2Min

	sum := addend1 + addend2
	level := 1 // Default level for custom ranges

	hintText := g.generateSingleDigitHint(addend1, addend2)

	return &models.Question{
		QuestionID:        uuid.New().String(),
		TopicID:           topicID,
		DifficultyLevelID: difficultyLevelID,
		QuestionText:      fmt.Sprintf("What is %d + %d?", addend1, addend2),
		QuestionData: map[string]interface{}{
			"addend1": addend1,
			"addend2": addend2,
			"sum":     sum,
			"level":   level,
		},
		VisualHintData: map[string]interface{}{
			"hint_type": "counting",
			"addend1":   addend1,
			"addend2":   addend2,
			"sum":       sum,
			"color":     "blue",
			"level":     level,
		},
		VerbalHint:    hintText,
		CorrectAnswer: fmt.Sprintf("%d", sum),
		GeneratedAt:   time.Now(),
	}, nil
}

// generateSingleDigitHint provides hints for single-digit addition
func (g *AdditionGame) generateSingleDigitHint(addend1, addend2 int) string {
	sum := addend1 + addend2

	if addend2 == 0 {
		return fmt.Sprintf("Think: Adding zero doesn't change the number. %d + 0 = %d!", addend1, addend1)
	}
	if addend2 == 1 {
		return fmt.Sprintf("Think: Just count up one from %d. What comes next?", addend1)
	}
	if addend2 == 2 {
		return fmt.Sprintf("Think: Count up two from %d. The next two numbers are %d and %d!",
			addend1, addend1+1, addend1+2)
	}
	if sum == 10 {
		return fmt.Sprintf("Think: You're making a friendly 10! %d + %d = 10", addend1, addend2)
	}
	if addend1 == addend2 {
		return fmt.Sprintf("Think: Double %d! Two %ds make %d.", addend1, addend1, sum)
	}

	return fmt.Sprintf("Think: Start at %d and count forward %d times.", addend1, addend2)
}

// ValidateAnswer checks if the submitted answer is correct
func (g *AdditionGame) ValidateAnswer(correctAnswer, submittedAnswer string) (*models.AnswerValidationResult, error) {
	isCorrect := submittedAnswer == correctAnswer

	result := &models.AnswerValidationResult{
		IsCorrect:       isCorrect,
		CorrectAnswer:   correctAnswer,
		FeedbackMessage: g.getFeedbackMessage(isCorrect),
	}

	return result, nil
}

func (g *AdditionGame) getFeedbackMessage(isCorrect bool) string {
	if isCorrect {
		messages := []string{"Amazing!", "Excellent!", "Great job!", "You got it!", "Fantastic!", "Brilliant!"}
		return messages[rand.Intn(len(messages))]
	}
	messages := []string{"Not quite, but close!", "Let's try another one!", "Good try!", "Almost there!"}
	return messages[rand.Intn(len(messages))]
}

// GetGameDefinition returns the game metadata
func (g *AdditionGame) GetGameDefinition() *models.Game {
	return &models.Game{
		GameID:      "addition-facts",
		DisplayName: "Addition Adventure",
		Description: "Master addition and become a math explorer!",
		IconName:    "addition-icon.svg",
		Version:     "1.0.0",
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
}
