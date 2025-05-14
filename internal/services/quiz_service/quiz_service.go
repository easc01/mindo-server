package quizservice

import (
	"database/sql"
	"fmt"

	"github.com/easc01/mindo-server/internal/models"
	aiservice "github.com/easc01/mindo-server/internal/services/ai_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateAndSaveQuiz(
	c *gin.Context,
	params dto.GenerateQuizParams,
) (dto.GeneratedQuiz, error) {
	generatedQuiz, err := aiservice.GenerateQuiz(params)

	if err != nil {
		logger.Log.Errorf("failed to generate quiz - %s, because %s", params.TopicName, err.Error())
		return dto.GeneratedQuiz{}, err
	}

	savedQuiz, err := db.Queries.SaveQuiz(c, models.SaveQuizParams{
		Name:         util.GetSQLNullString(params.TopicName),
		ThumbnailUrl: util.GetSQLNullString(""),
		PlayCount: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
		UpdatedBy: uuid.NullUUID{
			Valid: false,
		},
	})

	for i, question := range generatedQuiz.Questions {
		options := make([]string, 0)
		for _, option := range question.Options {
			options = append(options, option.Option)
		}

		ques, _ := db.Queries.SaveQuizQuestion(c, models.SaveQuizQuestionParams{
			QuizID:        savedQuiz.ID,
			Question:      util.GetSQLNullString(question.Question),
			Options:       options,
			CorrectOption: sql.NullInt32{Int32: int32(question.CorrectOption), Valid: true},
			UpdatedBy:     uuid.NullUUID{Valid: false},
		})
		generatedQuiz.Questions[i].QuestionId = ques.ID.String()
	}

	if err != nil {
		return dto.GeneratedQuiz{}, err
	}

	generatedQuiz.QuizId = savedQuiz.ID.String()
	return generatedQuiz, nil
}

func VerifyQuizResults(
	c *gin.Context,
	params dto.VerifyQuizParams,
) (dto.VerifyQuizResults, error) {
	quizUuid, _ := uuid.Parse(params.QuizId)
	questions, err := db.Queries.GetQuestionsByQuizId(c, quizUuid)

	if err != nil || len(questions) == 0 {
		return dto.VerifyQuizResults{}, nil
	}

	correctResponses := 0

	for _, question := range params.Questions {
		var questionData models.QuizQuestion

		for _, ques := range questions {
			fmt.Println(question.QuestionId, ques.ID.String())

			if question.QuestionId == ques.ID.String() {
				questionData = ques
				break
			}
		}

		if question.AttemptedOption == int(questionData.CorrectOption.Int32) {
			correctResponses++
		}
	}

	return dto.VerifyQuizResults{
		Grade: util.GetGrade(10, correctResponses),
		Marks: float64(correctResponses) / 10.00 * 100,
	}, nil
}
