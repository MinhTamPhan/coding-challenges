package routers

import (
	"net/http"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

var (
	LeaderboardRouterProviders = wire.NewSet(
		NewHttpHandler,
		NewEcho,
		ProvideMiddleware,
		NewLeaderBoardRouter,
	)
)

func ProvideMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func NewEcho(echoMiddlewares []echo.MiddlewareFunc) *echo.Echo {
	echoRouter := echo.New()
	if len(echoMiddlewares) > 0 {
		echoRouter.Use(echoMiddlewares...)
	}
	return echoRouter
}

func NewHttpHandler(handler *echo.Echo, leaderBoardRouter LeaderBoardRouter) http.Handler {
	handler.GET("/sessions", leaderBoardRouter.Sessions)
	handler.GET("/sessions/:session_id/leader-board", leaderBoardRouter.LeaderBoard)
	handler.POST("/sessions/:session_id/participants/:participant_id", leaderBoardRouter.AnswerQuiz)
	return handler
}
