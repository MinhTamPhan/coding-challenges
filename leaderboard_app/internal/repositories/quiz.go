package repositories

import "context"

type QuizRepository interface {
	GetAnswer(ctx context.Context, sessionID string, quizID int) (int, error)
	IsParticipantAnswered(ctx context.Context, sessionID string, participantID int) (bool, error)
	MarkParticipantAnswered(ctx context.Context, sessionID string, participantID int) error
}

type inMemQuizRepository struct {
	answers      map[string]map[int]int
	participants map[string]map[int]bool
}

func NewInMemQuizRepository() QuizRepository {
	return &inMemQuizRepository{
		answers: map[string]map[int]int{
			"session1": {
				1: 1,
				2: 2,
				3: 3,
			},
			"session2": {
				1: 1,
				2: 2,
				3: 3,
			},
			"session3": {
				1: 1,
				2: 2,
				3: 3,
			},
		},
		participants: map[string]map[int]bool{},
	}
}

func (i *inMemQuizRepository) GetAnswer(ctx context.Context, sessionID string, quizID int) (int, error) {
	if answers, ok := i.answers[sessionID]; ok {
		if answer, ok := answers[quizID]; ok {
			return answer, nil
		}
	}
	return 0, nil
}

func (i *inMemQuizRepository) IsParticipantAnswered(ctx context.Context, sessionID string, participantID int) (bool, error) {
	if participants, ok := i.participants[sessionID]; ok {
		if _, ok := participants[participantID]; ok {
			return true, nil
		}
	}
	return false, nil
}

func (i *inMemQuizRepository) MarkParticipantAnswered(ctx context.Context, sessionID string, participantID int) error {
	if _, ok := i.participants[sessionID]; !ok {
		i.participants[sessionID] = map[int]bool{}
	}
	i.participants[sessionID][participantID] = true
	return nil
}
