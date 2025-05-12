package quizservice

import (
	aiservice "github.com/easc01/mindo-server/internal/services/ai_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/gin-gonic/gin"
)

func GenerateStatelessQuiz(
	c *gin.Context,
	params dto.GenerateQuizParams,
) (dto.GeneratedQuiz, error) {
	generatedQuiz, err := aiservice.GenerateQuiz(params)

	if err != nil {
		logger.Log.Errorf("failed to generate quiz - %s, because %s", params.TopicName, err.Error())
		return dto.GeneratedQuiz{}, err
	}

	return generatedQuiz, nil
}
