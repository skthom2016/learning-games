import React, { useEffect, useRef, useCallback } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import confetti from 'canvas-confetti';
import { playCorrectSound, playIncorrectSound, playCelebrationSound } from '../../utils/soundEffects';
import './FeedbackScreen.css';

// Global timer registry that persists across component mounts
const timerRegistry = {
  timer: null,
  interval: null,
  scheduled: false,
};

/**
 * FeedbackScreen - Show feedback after answer submission
 * See ARCHITECTURE.md Section 6 - Screen 5
 *
 * TODO:
 * - Display celebration message if correct
 * - Display encouragement message if incorrect
 * - Animate star earning (star flies into counter)
 * - Show correct answer if was incorrect
 * - Auto-advance to next question after 2-3 seconds
 * - Add confetti animation for special achievements (Phase 2)
 */
function FeedbackScreen() {
  const navigate = useNavigate();
  const location = useLocation();
  const {
    playerId,
    gameId,
    questionText,
    isCorrect,
    starsEarned,
    correctAnswer,
    submittedAnswer,
  } = location.state || {};

  const hasPlayedSound = useRef(false);

  const handleNextQuestion = useCallback(() => {
    navigate('/play', { state: { playerId, gameId } });
  }, [navigate, playerId, gameId]);

  useEffect(() => {
    // Play sound effects and trigger confetti on first render
    if (!hasPlayedSound.current) {
      hasPlayedSound.current = true;

      if (isCorrect) {
        console.log('FeedbackScreen: Correct answer! Playing sounds and showing confetti...');

        // Play cheerful sound for correct answer
        playCorrectSound();
        playCelebrationSound();

        // ===== REDUCED CONFETTI CELEBRATION =====
        // Trigger confetti celebration - immediate burst
        const end = Date.now() + 1000; // Reduced from 2000ms to 1000ms

        // First immediate burst from center (smaller)
        console.log('Triggering confetti burst...');
        confetti({
          particleCount: 50, // Reduced from 100
          spread: 50, // Reduced from 70
          origin: { y: 0.6, x: 0.5 },
          colors: ['#FFD700', '#FFA500', '#FF6347', '#00CED1'],
          zIndex: 9999,
        });

        // Fewer side bursts, less frequent
        timerRegistry.interval = setInterval(function() {
          const timeLeft = end - Date.now();

          if (timeLeft <= 0) {
            clearInterval(timerRegistry.interval);
            return;
          }

          // Left side burst (smaller)
          confetti({
            particleCount: 15, // Reduced from 30
            angle: 60,
            spread: 45, // Reduced from 55
            origin: { x: 0.1, y: 0.5 },
            colors: ['#FFD700', '#FFA500'],
            zIndex: 9999,
          });

          // Right side burst (smaller)
          confetti({
            particleCount: 15, // Reduced from 30
            angle: 120,
            spread: 45, // Reduced from 55
            origin: { x: 0.9, y: 0.5 },
            colors: ['#00CED1', '#32CD32'],
            zIndex: 9999,
          });
        }, 400); // Increased interval from 300ms to 400ms (less frequent)

        /*
        ===== ORIGINAL FULL CONFETTI (COMMENTED OUT) =====
        // If you want to revert to the full celebration, uncomment this section
        // and comment out the reduced version above.

        const end = Date.now() + 2000;

        // First immediate burst from center
        console.log('Triggering confetti burst...');
        confetti({
          particleCount: 100,
          spread: 70,
          origin: { y: 0.6, x: 0.5 },
          colors: ['#FFD700', '#FFA500', '#FF6347', '#00CED1', '#32CD32', '#FF69B4'],
          zIndex: 9999,
        });

        // Continuous side bursts
        timerRegistry.interval = setInterval(function() {
          const timeLeft = end - Date.now();

          if (timeLeft <= 0) {
            clearInterval(timerRegistry.interval);
            return;
          }

          // Left side burst
          confetti({
            particleCount: 30,
            angle: 60,
            spread: 55,
            origin: { x: 0, y: 0.5 },
            colors: ['#FFD700', '#FFA500', '#FF6347'],
            zIndex: 9999,
          });

          // Right side burst
          confetti({
            particleCount: 30,
            angle: 120,
            spread: 55,
            origin: { x: 1, y: 0.5 },
            colors: ['#00CED1', '#32CD32', '#FF69B4'],
            zIndex: 9999,
          });
        }, 300);
        */

        // Only schedule navigation once using global registry
        if (!timerRegistry.scheduled) {
          timerRegistry.scheduled = true;

          // Auto-advance after 2 seconds (2000ms) for correct answers
          console.log('Setting auto-advance timer for 2000ms (2 seconds)...');
          timerRegistry.timer = setTimeout(() => {
            console.log('Auto-advancing to next question...');
            timerRegistry.scheduled = false;
            handleNextQuestion();
          }, 2000);
        }

        return () => {
          console.log('Component unmounting (not clearing timer - it persists in global registry)');
        };
      } else {
        console.log('FeedbackScreen: Incorrect answer. Playing incorrect sound.');
        // Play gentle "oh-oh" sound for incorrect answer
        playIncorrectSound();
      }
    }
  }, []);

  const celebrationMessages = [
    'Amazing!',
    'Excellent!',
    'Great job!',
    'You got it!',
    'Fantastic!',
    'Brilliant!',
  ];

  const encouragementMessages = [
    'Not quite, but close!',
    'Let\'s try another one!',
    'Good try!',
    'Almost there!',
  ];

  const getMessage = () => {
    if (isCorrect) {
      return celebrationMessages[Math.floor(Math.random() * celebrationMessages.length)];
    }
    return encouragementMessages[Math.floor(Math.random() * encouragementMessages.length)];
  };

  return (
    <div className={`feedback-screen ${isCorrect ? 'correct' : 'incorrect'}`}>
      <div className="feedback-content">
        {isCorrect ? (
          <>
            <h1 className="feedback-title">⭐ {getMessage()} ⭐</h1>
            <p className="feedback-answer">Your answer: {submittedAnswer} ✓</p>
            <div className="stars-earned">
              <p>You earned {starsEarned} stars!</p>
              {/* TODO: Add star animation (stars flying into counter) */}
              <div className="star-animation">
                {Array.from({ length: Math.min(starsEarned, 10) }).map((_, i) => (
                  <span key={i} className="star-icon">⭐</span>
                ))}
              </div>
            </div>
          </>
        ) : (
          <>
            <h1 className="feedback-title">{getMessage()}</h1>
            <div className="question-display">
              <p className="question-label">Question:</p>
              <p className="question-text">{questionText}</p>
            </div>
            <div className="answer-comparison">
              <p className="your-answer">You said: {submittedAnswer}</p>
              <p className="correct-answer">Correct answer: {correctAnswer}</p>
            </div>
            <p className="encouragement">Let's try another one!</p>
          </>
        )}

        <button className="next-button" onClick={handleNextQuestion}>
          Next Question →
        </button>
      </div>
    </div>
  );
}

export default FeedbackScreen;
