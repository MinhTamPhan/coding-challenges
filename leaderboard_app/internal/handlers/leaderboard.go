package handlers

import (
	"context"
	"errors"

	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/dtos"
	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/services"
	"github.com/google/wire"
	"github.com/samber/lo"
)

const (
	CorrectAnswer = "correct"
	WrongAnswer   = "wrong"
)

var (
	ErrorParticipantAnswered = errors.New("participant already answered")
)

var (
	LeaderBoardHandlerProviders = wire.NewSet(NewLeaderBoardHandler)
)

type LeaderBoardHandler interface {
	Sessions(ctx context.Context) (*dtos.SessionResponse, error)
	LeaderBoard(ctx context.Context, sessionID string) (*dtos.LeaderboardResponse, error)
	AnswerQuiz(ctx context.Context, sessionID string, participantID, quizID, answer int) (*dtos.AnswerQuizResponse, error)
}

type leaderBoardHandler struct {
	leaderBoardService services.LeaderBoardService
}

func NewLeaderBoardHandler(leaderBoardService services.LeaderBoardService) LeaderBoardHandler {
	return &leaderBoardHandler{
		leaderBoardService: leaderBoardService,
	}
}

func (l *leaderBoardHandler) Sessions(ctx context.Context) (*dtos.SessionResponse, error) {
	sessions, err := l.leaderBoardService.Sessions(ctx)
	if err != nil {
		return nil, err
	}
	items := lo.Map(sessions, func(item string, _ int) *dtos.SessionItem {
		return &dtos.SessionItem{
			SessionID: item,
		}
	})
	return &dtos.SessionResponse{
		Items: items,
	}, nil
}

func (l *leaderBoardHandler) LeaderBoard(ctx context.Context, sessionID string) (*dtos.LeaderboardResponse, error) {
	leaderBoard, err := l.leaderBoardService.LeaderBoard(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	items := lo.Map(leaderBoard, func(item *services.LeaderBoardItem, idx int) *dtos.LeaderboardItem {
		return &dtos.LeaderboardItem{
			ParticipantID: item.ParticipantID,
			Score:         item.Score,
			Rank:          idx + 1,
		}
	})
	return &dtos.LeaderboardResponse{
		Items: items,
	}, nil
}

func (l *leaderBoardHandler) AnswerQuiz(ctx context.Context, sessionID string,
	participantID, quizID, answer int) (*dtos.AnswerQuizResponse, error) {
	result, err := l.leaderBoardService.AnswerQuiz(ctx, sessionID, participantID, quizID, answer)
	if err != nil && errors.Is(err, services.ErrParticipantAnswered) {
		return nil, ErrorParticipantAnswered
	}
	if err != nil {
		return nil, err
	}
	resp := CorrectAnswer
	if !result {
		resp = WrongAnswer
	}
	score, err := l.leaderBoardService.GetScore(ctx, sessionID, participantID)
	if err != nil {
		return nil, err
	}
	return &dtos.AnswerQuizResponse{
		Result:       resp,
		CurrentScore: score.Point,
		CurrentRank:  score.Rank,
	}, nil
}
