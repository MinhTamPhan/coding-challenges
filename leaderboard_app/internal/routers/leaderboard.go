package routers

import (
	"errors"
	"net/http"

	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/dtos"
	"github.com/MinhTamPhan/coding-challenges/leaderboard_app/internal/handlers"
	"github.com/labstack/echo/v4"
)

type LeaderBoardRouter interface {
	Sessions(c echo.Context) error
	LeaderBoard(c echo.Context) error
	AnswerQuiz(c echo.Context) error
}

type leaderBoardRouter struct {
	leaderBoardHandler handlers.LeaderBoardHandler
}

func NewLeaderBoardRouter(leaderBoardHandler handlers.LeaderBoardHandler) LeaderBoardRouter {
	return &leaderBoardRouter{
		leaderBoardHandler: leaderBoardHandler,
	}
}

func (l *leaderBoardRouter) Sessions(c echo.Context) error {
	sessions, err := l.leaderBoardHandler.Sessions(c.Request().Context())
	if err != nil {
		return err
	}
	response := dtos.ServerResponseDTO[*dtos.SessionResponse]{
		Code:    http.StatusOK,
		Message: "Success",
		Result:  sessions,
	}
	return c.JSON(http.StatusOK, response)
}

func (l *leaderBoardRouter) LeaderBoard(c echo.Context) error {
	request := &dtos.LeaderboardRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	leaderBoard, err := l.leaderBoardHandler.LeaderBoard(c.Request().Context(), request.SessionID)
	if err != nil {
		return err
	}
	response := dtos.ServerResponseDTO[*dtos.LeaderboardResponse]{
		Code:    http.StatusOK,
		Message: "Success",
		Result:  leaderBoard,
	}
	return c.JSON(http.StatusOK, response)
}

func (l *leaderBoardRouter) AnswerQuiz(c echo.Context) error {
	request := &dtos.AnswerQuizRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	result, err := l.leaderBoardHandler.AnswerQuiz(c.Request().Context(),
		request.SessionID, request.ParticipantID, request.QuizID, request.Answer)

	if err != nil && errors.Is(err, handlers.ErrorParticipantAnswered) {
		response := dtos.ServerResponseDTO[*dtos.AnswerQuizResponse]{
			Code:    http.StatusConflict,
			Message: err.Error(),
		}
		return c.JSON(http.StatusConflict, response)
	}

	if err != nil {
		return err
	}
	response := dtos.ServerResponseDTO[*dtos.AnswerQuizResponse]{
		Code:    http.StatusOK,
		Message: "Success",
		Result:  result,
	}
	return c.JSON(http.StatusOK, response)
}
