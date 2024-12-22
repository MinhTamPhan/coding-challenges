package repositories

import "github.com/google/wire"

var (
	RepositoriesProviders = wire.NewSet(
		NewInMemQuizRepository,
		NewInMemLeaderboardRepository,
	)
)
