import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

console.log('API Client initialized with base URL:', API_BASE_URL);

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor for debugging
apiClient.interceptors.request.use(
  (config) => {
    console.log('API Request:', config.method.toUpperCase(), config.url, config.params || config.data);
    return config;
  },
  (error) => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for debugging
apiClient.interceptors.response.use(
  (response) => {
    console.log('API Response:', response.status, response.config.url);
    return response;
  },
  (error) => {
    console.error('API Response Error:', error.message);
    if (error.response) {
      console.error('Error response:', error.response.status, error.response.data);
    }
    return Promise.reject(error);
  }
);

// API methods
// TODO: Implement all API calls matching backend routes

export const api = {
  // Player Management
  createPlayer: (playerName) => {
    // TODO: POST /api/players
    return apiClient.post('/api/players', { player_name: playerName });
  },

  listPlayers: () => {
    // TODO: GET /api/players
    return apiClient.get('/api/players');
  },

  getPlayer: (playerId) => {
    // TODO: GET /api/players/:id
    return apiClient.get(`/api/players/${playerId}`);
  },

  deletePlayer: (playerId) => {
    return apiClient.delete(`/api/players/${playerId}`);
  },

  // Session Management
  startSession: (playerId, gameId) => {
    // TODO: POST /api/sessions/start
    return apiClient.post('/api/sessions/start', {
      player_id: playerId,
      game_id: gameId,
    });
  },

  endSession: (playerId, gameId, sessionId) => {
    // TODO: POST /api/sessions/end
    return apiClient.post('/api/sessions/end', {
      player_id: playerId,
      game_id: gameId,
      session_id: sessionId,
    });
  },

  // Game Play
  getNextQuestion: (playerId, gameId) => {
    // TODO: GET /api/questions/next?player_id=X&game_id=Y
    return apiClient.get('/api/questions/next', {
      params: { player_id: playerId, game_id: gameId },
    });
  },

  submitAnswer: (playerId, gameId, questionId, topicId, difficultyLevelId, submittedAnswer, correctAnswer, hintUsed, timeToAnswerMs) => {
    return apiClient.post('/api/answers/submit', {
      player_id: playerId,
      game_id: gameId,
      question_id: questionId,
      topic_id: topicId,
      difficulty_level_id: difficultyLevelId,
      submitted_answer: submittedAnswer,
      correct_answer: correctAnswer,
      hint_used: hintUsed,
      time_to_answer_ms: timeToAnswerMs,
    });
  },

  // Progress & Analytics
  getPlayerProgress: (playerId, gameId) => {
    // TODO: GET /api/progress/:player_id/:game_id
    return apiClient.get(`/api/progress/${playerId}/${gameId}`);
  },

  getMasteryHeatmap: (playerId, gameId) => {
    // TODO: GET /api/mastery/:player_id/:game_id
    return apiClient.get(`/api/mastery/${playerId}/${gameId}`);
  },

  // Rewards
  getRewardBalance: (playerId) => {
    // TODO: GET /api/rewards/:player_id
    return apiClient.get(`/api/rewards/${playerId}`);
  },

  getRewardTransactions: (playerId) => {
    // TODO: GET /api/rewards/:player_id/transactions
    return apiClient.get(`/api/rewards/${playerId}/transactions`);
  },

  // Admin
  resetRewards: (playerId, starsToConsolidate, realWorldReward) => {
    // TODO: POST /api/admin/rewards/reset
    return apiClient.post('/api/admin/rewards/reset', {
      player_id: playerId,
      stars_to_consolidate: starsToConsolidate,
      real_world_reward: realWorldReward,
    });
  },

  getPlayerAnalytics: (playerId) => {
    // TODO: GET /api/admin/analytics/:player_id
    return apiClient.get(`/api/admin/analytics/${playerId}`);
  },

  // Player Number Ranges (Admin: grade-based difficulty customization)
  getPlayerNumberRanges: (playerId) => {
    return apiClient.get(`/api/players/${playerId}/number-ranges`);
  },

  setPlayerNumberRanges: (playerId, numberRanges) => {
    return apiClient.put(`/api/players/${playerId}/number-ranges`, {
      number_ranges: numberRanges,
    });
  },

  deletePlayerNumberRange: (playerId, gameId) => {
    return apiClient.delete(`/api/players/${playerId}/number-ranges/${gameId}`);
  },
};

export default apiClient;
