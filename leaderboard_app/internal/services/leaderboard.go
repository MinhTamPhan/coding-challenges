package services

import (
	"context"
	"errors"

	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/repositories"
	"github.com/google/wire"
	"github.com/samber/lo"
)

var (
	LeaderBoardServiceProviders = wire.NewSet(
		NewLeaderBoardService,
	)
)

var (
	ErrParticipantAnswered = errors.New("participant already answered")
)

type LeaderBoardItem struct {
	ParticipantID int
	Score         int
}

type Score struct {
	Point int
	Rank  int
}

type LeaderBoardService interface {
	Sessions(ctx context.Context) ([]string, error)
	LeaderBoard(ctx context.Context, sessionID string) ([]*LeaderBoardItem, error)
	AnswerQuiz(ctx context.Context, sessionID string, participantID, quizID, answer int) (bool, error)
	GetScore(ctx context.Context, sessionID string, participantID int) (*Score, error)
}

type leaderBoardService struct {
	quizRepository repositories.QuizRepository
	leaderBoard    repositories.LeaderboardRepository
}

func NewLeaderBoardService(quizRepository repositories.QuizRepository,
	leaderBoard repositories.LeaderboardRepository) LeaderBoardService {
	return &leaderBoardService{
		quizRepository: quizRepository,
		leaderBoard:    leaderBoard,
	}
}

func (l *leaderBoardService) Sessions(ctx context.Context) ([]string, error) {
	return []string{"session1", "session2", "session3"}, nil
}

func (l *leaderBoardService) LeaderBoard(ctx context.Context, sessionID string) ([]*LeaderBoardItem, error) {
	leaderboard, err := l.leaderBoard.LeaderBoardBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	result := lo.Map(leaderboard, func(item *repositories.LeaderboardItem, _ int) *LeaderBoardItem {
		return &LeaderBoardItem{
			ParticipantID: item.ParticipantID,
			Score:         item.Score,
		}
	})
	return result, nil
}

func (l *leaderBoardService) AnswerQuiz(ctx context.Context, sessionID string, participantID, quizID, answer int) (bool, error) {
	isAnswered, err := l.quizRepository.IsParticipantAnswered(ctx, sessionID, participantID)
	if err != nil {
		return false, err
	}
	if isAnswered {
		return false, ErrParticipantAnswered
	}
	correctAnswer, err := l.quizRepository.GetAnswer(ctx, sessionID, quizID)
	if err != nil {
		return false, err
	}
	if correctAnswer != answer {
		return false, nil
	}
	err = l.leaderBoard.UpdateScore(ctx, sessionID, participantID, 1)
	if err != nil {
		return false, err
	}
	err = l.quizRepository.MarkParticipantAnswered(ctx, sessionID, participantID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *leaderBoardService) GetScore(ctx context.Context, sessionID string, participantID int) (*Score, error) {
	score, rank, err := l.leaderBoard.GetScore(ctx, sessionID, participantID)
	if err != nil {
		return nil, err
	}
	return &Score{
		Point: score,
		Rank:  rank,
	}, nil
}
