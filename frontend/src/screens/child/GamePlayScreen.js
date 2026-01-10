import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { api } from '../../api/client';
import './GamePlayScreen.css';

/**
 * GamePlayScreen - Main gameplay screen
 * See ARCHITECTURE.md Section 6 - Screen 4
 *
 * TODO:
 * - Fetch next question from API
 * - Display question text and answer input
 * - Implement "Show Hint" button (conditional on difficulty)
 * - Track time to answer (optional, for analytics)
 * - On submit, validate answer and navigate to /feedback
 * - Implement visual hint rendering (dot array for multiplication)
 * - Add quit button that ends session
 */
function GamePlayScreen() {
  const navigate = useNavigate();
  const location = useLocation();
  const { playerId, gameId } = location.state || {};

  const [question, setQuestion] = useState(null);
  const [answer, setAnswer] = useState('');
  const [showHint, setShowHint] = useState(false);
  const [startTime, setStartTime] = useState(null);
  const [loading, setLoading] = useState(true);
  const [starCount, setStarCount] = useState(0);

  useEffect(() => {
    console.log('GamePlayScreen mounted', { playerId, gameId });
    if (!playerId || !gameId) {
      console.log('Missing playerId or gameId, redirecting to home');
      navigate('/');
      return;
    }
    fetchNextQuestion();
    fetchStarCount();
  }, [playerId, gameId, navigate]);

  const fetchStarCount = async () => {
    try {
      console.log('Fetching star count for player:', playerId);
      const response = await api.getRewardBalance(playerId);
      console.log('Star count response:', response.data);
      setStarCount(response.data.total_stars || 0);
    } catch (error) {
      console.error('Failed to fetch star count:', error);
      console.error('Error details:', error.response?.data || error.message);
    }
  };

  const fetchNextQuestion = async () => {
    try {
      setLoading(true);
      console.log('Fetching next question for:', { playerId, gameId });
      const response = await api.getNextQuestion(playerId, gameId);
      console.log('Question response:', response.data);
      setQuestion(response.data);

      // Also refresh star count when getting new question
      await fetchStarCount();

      setAnswer('');
      setShowHint(false);
      setStartTime(Date.now());
      setLoading(false);
    } catch (error) {
      console.error('Failed to fetch question:', error);
      console.error('Error details:', error.response?.data || error.message);
      setLoading(false);
    }
  };

  const handleSubmit = async () => {
    console.log('handleSubmit called', { answer, question, playerId, gameId });

    if (!answer || !answer.trim()) {
      console.log('Answer is empty, returning');
      return;
    }

    if (!question || !question.question_id) {
      console.error('Question is not loaded properly');
      return;
    }

    const timeToAnswer = Date.now() - startTime;

    const submittedAnswer = String(answer).trim();
    // Calculate correct answer from question_data (backend doesn't send it for security)
    let correctAnswer;
    if (gameId === 'division-facts') {
      // Division: dividend ÷ divisor = quotient
      const quotient = question.question_data?.quotient || 0;
      correctAnswer = String(quotient);
    } else if (gameId === 'addition-facts') {
      // Addition: addend1 + addend2 = sum
      const sum = question.question_data?.sum || 0;
      correctAnswer = String(sum);
    } else if (gameId === 'subtraction-facts') {
      // Subtraction: minuend - subtrahend = difference
      const difference = question.question_data?.difference || 0;
      correctAnswer = String(difference);
    } else {
      // Multiplication: operand1 × operand2 = product
      const operand1 = question.question_data?.operand1 || 0;
      const operand2 = question.question_data?.operand2 || 0;
      correctAnswer = String(operand1 * operand2);
    }

    console.log('Submitting answer...', {
      questionId: question.question_id,
      topicId: question.topic_id,
      difficultyLevelId: question.difficulty_level_id,
      gameId,
      submittedAnswer,
      correctAnswer,
      questionData: question.question_data,
      hintUsed: showHint,
      timeToAnswer
    });

    try {
      const response = await api.submitAnswer(
        playerId,
        gameId,
        question.question_id,
        question.topic_id,
        question.difficulty_level_id,
        submittedAnswer,
        correctAnswer,
        showHint,
        timeToAnswer
      );

      console.log('Submit successful:', response.data);

      // Navigate to feedback screen with result from API
      navigate('/feedback', {
        state: {
          playerId,
          gameId,
          questionText: question.question_text,
          isCorrect: response.data.is_correct,
          starsEarned: response.data.stars_earned,
          correctAnswer: response.data.correct_answer,
          submittedAnswer: answer,
        },
      });
    } catch (error) {
      console.error('Failed to submit answer:', error);
      console.error('Error details:', error.response?.data || error.message);
    }
  };

  const handleQuit = async () => {
    try {
      await api.endSession(playerId, gameId, null);
      navigate('/games', { state: { playerId } });
    } catch (error) {
      console.error('Failed to end session:', error);
      navigate('/games', { state: { playerId } });
    }
  };

  const renderVisualHint = () => {
    if (!question?.visual_hint_data) return null;

    const { rows, cols, color } = question.visual_hint_data;

    // Only show visual dots for small numbers (up to 5x5 = 25 dots)
    const showDots = rows <= 5 && cols <= 5;

    return (
      <div className="visual-hint">
        {showDots && <p>Count the dots to find the answer!</p>}
        {showDots && (
          <div className="dot-array">
            {Array.from({ length: rows }).map((_, r) => (
              <div key={r} className="dot-row">
                {Array.from({ length: cols }).map((_, c) => (
                  <div key={c} className="dot" style={{ background: color }} />
                ))}
              </div>
            ))}
          </div>
        )}
        {question.verbal_hint && (
          <div className="verbal-hint">
            <p className="verbal-hint-text">{question.verbal_hint}</p>
          </div>
        )}
      </div>
    );
  };

  if (loading || !question) {
    return <div className="loading">Loading question...</div>;
  }

  return (
    <div className="game-play-screen">
      <div className="top-bar">
        <div className="star-count">⭐ {starCount}</div>
        <button className="quit-button" onClick={handleQuit}>❮ Quit</button>
      </div>

      <div className="question-area">
        <h1 className="question-text">{question.question_text}</h1>

        <div className="answer-section">
          <input
            type="number"
            className="answer-input"
            value={answer}
            onChange={(e) => setAnswer(e.target.value)}
            placeholder="Your answer"
            autoFocus
          />
        </div>

        {!showHint && (
          <button
            className="hint-button"
            onClick={() => setShowHint(true)}
          >
            Show Hint 💡
          </button>
        )}

        {showHint && renderVisualHint()}

        <button
          className="submit-button"
          onClick={(e) => {
            console.log('Submit button clicked!', { answer, disabled: !answer.trim() });
            handleSubmit(e);
          }}
          disabled={!answer.trim()}
        >
          Submit ✓
        </button>
      </div>
    </div>
  );
}

export default GamePlayScreen;
