package repositories

import (
	"context"
	"sort"

	"github.com/samber/lo"
)

type LeaderboardItem struct {
	ParticipantID int
	Score         int
}

type LeaderboardRepository interface {
	UpdateScore(ctx context.Context, sessionID string, participantID int, score int) error
	GetScore(ctx context.Context, sessionID string, participantID int) (int, int, error)
	LeaderBoardBySessionID(ctx context.Context, sessionID string) ([]*LeaderboardItem, error)
}

type inMemLeaderboardRepository struct {
	scores map[string]map[int]int
}

func NewInMemLeaderboardRepository() LeaderboardRepository {
	return &inMemLeaderboardRepository{
		scores: map[string]map[int]int{},
	}
}

func (i *inMemLeaderboardRepository) UpdateScore(ctx context.Context, sessionID string, participantID int, score int) error {
	if _, ok := i.scores[sessionID]; !ok {
		i.scores[sessionID] = map[int]int{}
	}
	i.scores[sessionID][participantID] += score
	return nil
}

func (i *inMemLeaderboardRepository) GetScore(ctx context.Context, sessionID string, participantID int) (int, int, error) {
	leaderBoard, err := i.LeaderBoardBySessionID(ctx, sessionID)
	if err != nil {
		return 0, 0, err
	}
	for idx, item := range leaderBoard {
		if item.ParticipantID == participantID {
			return item.Score, idx + 1, nil
		}
	}
	return 0, 0, nil
}

func (i *inMemLeaderboardRepository) LeaderBoardBySessionID(ctx context.Context, sessionID string) ([]*LeaderboardItem, error) {
	if _, ok := i.scores[sessionID]; ok {
		leaderBoard := lo.MapToSlice(i.scores[sessionID], func(key int, val int) *LeaderboardItem {
			return &LeaderboardItem{
				ParticipantID: key,
				Score:         val,
			}
		})
		sort.Slice(leaderBoard, func(i, j int) bool {
			return leaderBoard[i].Score > leaderBoard[j].Score
		})
		return leaderBoard, nil
	}
	return nil, nil
}
