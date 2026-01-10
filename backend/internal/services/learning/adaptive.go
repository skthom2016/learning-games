package learning

import (
	"math/rand"
	"time"

	"github.com/learning-game/backend/internal/models"
)

// AdaptiveLearningEngine handles mastery tracking and question selection
type AdaptiveLearningEngine struct{}

// SelectNextTopic uses mastery records to select the best next topic
func (e *AdaptiveLearningEngine) SelectNextTopic(masteryRecords []*models.MasteryRecord, recentTopics []string, inRecoveryMode bool) (string, error) {
	if len(masteryRecords) == 0 {
		return "times-2", nil // Default to first topic
	}

	// Check if all topics are in "unknown" state
	allUnknown := true
	for _, record := range masteryRecords {
		if record.MasteryState != models.MasteryStateUnknown {
			allUnknown = false
			break
		}
	}

	// If all topics are unknown, return the first one (lowest level)
	// This ensures new players start at Level 1 and progress sequentially
	if allUnknown {
		return masteryRecords[0].TopicID, nil
	}

	// Calculate weights for each topic
	weights := make(map[string]float64)

	for _, record := range masteryRecords {
		// Base weight by mastery state
		var baseWeight float64
		switch record.MasteryState {
		case models.MasteryStateUnknown:
			baseWeight = 100
		case models.MasteryStateWeak:
			baseWeight = 200
		case models.MasteryStateLearning:
			baseWeight = 80
		case models.MasteryStateStrong:
			baseWeight = 30
		case models.MasteryStateMastered:
			baseWeight = 5
		default:
			baseWeight = 100
		}

		// Apply modifiers
		// Check if practiced recently
		practicedRecently := false
		for _, recentTopic := range recentTopics {
			if recentTopic == record.TopicID {
				practicedRecently = true
				break
			}
		}

		if practicedRecently {
			baseWeight *= 0.5 // Reduce weight if practiced recently
		}

		// Confidence recovery mode modifier
		if inRecoveryMode {
			if record.MasteryState == models.MasteryStateWeak {
				baseWeight *= 3.0
			} else {
				baseWeight *= 0.5
			}
		}

		weights[record.TopicID] = baseWeight
	}

	// Normalize to probabilities and select
	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}

	if totalWeight == 0 {
		return masteryRecords[0].TopicID, nil
	}

	// Weighted random selection
	r := rand.Float64() * totalWeight
	cumulative := 0.0
	for topicID, weight := range weights {
		cumulative += weight
		if r <= cumulative {
			return topicID, nil
		}
	}

	return masteryRecords[0].TopicID, nil
}

// SelectDifficultyForTopic determines appropriate difficulty based on mastery
func (e *AdaptiveLearningEngine) SelectDifficultyForTopic(masteryRecord *models.MasteryRecord, recentAccuracy float64) string {
	switch masteryRecord.MasteryState {
	case models.MasteryStateUnknown, models.MasteryStateWeak:
		return "easy"
	case models.MasteryStateLearning:
		if recentAccuracy < 0.5 {
			return "easy"
		}
		return "medium"
	case models.MasteryStateStrong:
		// 70% medium, 30% hard
		if rand.Float64() < 0.7 {
			return "medium"
		}
		return "hard"
	case models.MasteryStateMastered:
		return "hard"
	default:
		return "easy"
	}
}

// AssessMastery recalculates mastery state from recent attempts
func (e *AdaptiveLearningEngine) AssessMastery(attempts []*models.AttemptRecord) models.MasteryState {
	if len(attempts) < 3 {
		return models.MasteryStateUnknown
	}

	// Calculate accuracy over last 10 attempts
	last10 := attempts
	if len(attempts) > 10 {
		last10 = attempts[len(attempts)-10:]
	}

	correct := 0
	for _, att := range last10 {
		if att.WasCorrect {
			correct++
		}
	}
	accuracy := float64(correct) / float64(len(last10))

	// Check for mastered state (needs more criteria)
	if accuracy >= 0.9 && len(attempts) >= 20 {
		// Check last 10 for no incorrect
		last10Correct := 0
		for _, att := range last10 {
			if att.WasCorrect {
				last10Correct++
			}
		}
		if last10Correct == 10 {
			// Check for 10 consecutive correct
			consecutiveCorrect := 0
			for i := len(attempts) - 1; i >= 0; i-- {
				if attempts[i].WasCorrect {
					consecutiveCorrect++
					if consecutiveCorrect >= 10 {
						return models.MasteryStateMastered
					}
				} else {
					break
				}
			}
		}
	}

	// Standard states
	if accuracy < 0.4 {
		return models.MasteryStateWeak
	} else if accuracy < 0.7 {
		return models.MasteryStateLearning
	} else if accuracy < 0.9 {
		return models.MasteryStateStrong
	}

	return models.MasteryStateStrong
}

// CheckConfidenceRecoveryTrigger determines if recovery mode should activate
func (e *AdaptiveLearningEngine) CheckConfidenceRecoveryTrigger(recentAttempts []*models.AttemptRecord, alreadyInRecovery bool) bool {
	if alreadyInRecovery || len(recentAttempts) < 7 {
		return false
	}

	// Require at least 10 total attempts before triggering confidence recovery
	// This prevents false positives during early gameplay when child is learning
	if len(recentAttempts) < 10 {
		return false
	}

	// Get last 7 attempts
	last7 := recentAttempts
	if len(recentAttempts) > 7 {
		last7 = recentAttempts[len(recentAttempts)-7:]
	}

	// Count incorrect in last 7
	incorrect := 0
	topicsMap := make(map[string]bool)

	for _, attempt := range last7 {
		if !attempt.WasCorrect {
			incorrect++
		}
		topicsMap[attempt.TopicID] = true
	}

	return incorrect >= 5 && len(topicsMap) >= 2
}

// CheckConfidenceRecoveryExit determines if recovery mode should end
func (e *AdaptiveLearningEngine) CheckConfidenceRecoveryExit(recentAttempts []*models.AttemptRecord, enteredAt *time.Time) bool {
	if enteredAt == nil {
		return false
	}

	// Check if 15 minutes have passed
	if time.Since(*enteredAt) > 15*time.Minute {
		return true
	}

	if len(recentAttempts) < 3 {
		return false
	}

	// Check for 3 consecutive correct
	if len(recentAttempts) >= 3 {
		last3 := recentAttempts[len(recentAttempts)-3:]
		allCorrect := true
		for _, attempt := range last3 {
			if !attempt.WasCorrect {
				allCorrect = false
				break
			}
		}
		if allCorrect {
			return true
		}
	}

	// Check for 5 out of last 7 correct
	if len(recentAttempts) >= 7 {
		last7 := recentAttempts[len(recentAttempts)-7:]
		correct := 0
		for _, attempt := range last7 {
			if attempt.WasCorrect {
				correct++
			}
		}
		if correct >= 5 {
			return true
		}
	}

	return false
}
