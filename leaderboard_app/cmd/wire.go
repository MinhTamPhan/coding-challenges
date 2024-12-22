//go:build wireinject
// +build wireinject

package main

import (
	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/applications"
	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/handlers"
	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/repositories"
	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/routers"
	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/services"
	"github.com/google/wire"
)

func wireLeaderBoardApiApplication() (applications.Application, error) {
	wire.Build(
		services.LeaderBoardServiceProviders,
		handlers.LeaderBoardHandlerProviders,
		routers.LeaderboardRouterProviders,
		applications.RestApplicationProviders,
		repositories.RepositoriesProviders,
	)
	return nil, nil
}
