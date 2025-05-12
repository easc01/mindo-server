package dto

type GenerateQuizParams struct {
	TopicName     string `json:"topicName"`
	QuestionCount int    `json:"questionCount"`
}

type GeneratedQuiz struct {
	Questions []GeneratedQuizQuestion `json:"questions"`
	TopicName string                  `json:"topicName"`
}

type GeneratedQuizQuestion struct {
	CorrectOption  int                    `json:"correctOption"`
	Options        []GeneratedQuizOptions `json:"options"`
	Question       string                 `json:"question"`
	QuestionNumber int                    `json:"questionNumber"`
}

type GeneratedQuizOptions struct {
	Option       string `json:"option"`
	OptionNumber int    `json:"optionNumber"`
}
