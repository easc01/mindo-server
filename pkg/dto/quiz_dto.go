package dto

type GenerateQuizParams struct {
	TopicName     string `json:"topicName"`
	QuestionCount int    `json:"questionCount"`
}

type GeneratedQuiz struct {
	QuizId    string                  `json:"quizId"`
	Questions []GeneratedQuizQuestion `json:"questions"`
	TopicName string                  `json:"topicName"`
}

type GeneratedQuizQuestion struct {
	QuestionId     string                 `json:"questionId"`
	CorrectOption  int                    `json:"correctOption"`
	Options        []GeneratedQuizOptions `json:"options"`
	Question       string                 `json:"question"`
	QuestionNumber int                    `json:"questionNumber"`
}

type GeneratedQuizOptions struct {
	Option       string `json:"option"`
	OptionNumber int    `json:"optionNumber"`
}

type VerifyQuizParams struct {
	QuizId    string               `json:"quizId"`
	Questions []VerifyQuizQuestion `json:"questions"`
}

type VerifyQuizQuestion struct {
	QuestionId      string `json:"questionId"`
	AttemptedOption int    `json:"attemptedOption"`
}

type VerifyQuizResults struct {
	Grade string  `json:"grade"`
	Marks float64 `json:"marks"`
}
