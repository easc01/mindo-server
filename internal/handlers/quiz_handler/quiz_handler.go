package quizhandler

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	quizservice "github.com/easc01/mindo-server/internal/services/quiz_service"
	"github.com/easc01/mindo-server/pkg/dto"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
	"github.com/easc01/mindo-server/pkg/utils/route"
	"github.com/gin-gonic/gin"
)

func RegisterQuiz(rg *gin.RouterGroup) {
	quizRg := rg.Group(route.Quizzes)

	{
		quizRg.POST(
			"/gen-ai",
			middleware.RequireRole(models.UserTypeAppUser),
			generateQuizHandler,
		)

		quizRg.POST(
			"/verify",
			middleware.RequireRole(models.UserTypeAppUser),
			verifyQuizAnswersHandler,
		)
	}
}

func generateQuizHandler(c *gin.Context) {
	topicName := c.Query("topicName")

	if topicName == "" {
		networkutil.NewErrorResponse(
			http.StatusBadRequest,
			"topicName query param is required",
			nil,
		).Send(c)
		return
	}

	quizData, err := quizservice.GenerateAndSaveQuiz(c, dto.GenerateQuizParams{
		TopicName:     topicName,
		QuestionCount: 10,
	})

	if err != nil {
		networkutil.NewErrorResponse(
			http.StatusInternalServerError,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	networkutil.NewResponse(
		http.StatusAccepted,
		quizData,
	).Send(c)
}

func verifyQuizAnswersHandler(c *gin.Context) {
	req, ok := networkutil.GetRequestBody[dto.VerifyQuizParams](c)
	if !ok {
		return
	}

	quizData, err := quizservice.VerifyQuizResults(c, req)

	if err != nil {
		networkutil.NewErrorResponse(
			http.StatusInternalServerError,
			err.Error(),
			nil,
		).Send(c)
		return
	}

	networkutil.NewResponse(
		http.StatusAccepted,
		quizData,
	).Send(c)
}
