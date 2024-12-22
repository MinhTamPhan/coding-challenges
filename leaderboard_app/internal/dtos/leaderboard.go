package dtos

type LeaderboardRequest struct {
	SessionID string `param:"session_id"`
}

type ServerResponseDTO[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  T      `json:"result,omitempty"`
}

type LeaderboardResponse struct {
	Items []*LeaderboardItem `json:"items"`
}

type LeaderboardItem struct {
	ParticipantID int `json:"participant_id"`
	Score         int `json:"score"`
	Rank          int `json:"rank"`
}

type SessionResponse struct {
	Items []*SessionItem `json:"items"`
}

type SessionItem struct {
	SessionID string `json:"session_id"`
}

type AnswerQuizRequest struct {
	SessionID     string `param:"session_id"`
	ParticipantID int    `param:"participant_id"`
	QuizID        int    `json:"quiz_id"`
	Answer        int    `json:"answer"`
}

type AnswerQuizResponse struct {
	Result       string `json:"result"`
	CurrentScore int    `json:"current_score"`
	CurrentRank  int    `json:"current_rank"`
}
