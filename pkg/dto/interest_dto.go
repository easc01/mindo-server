package dto

type UpsertInterestDTO struct {
	Interests []string `json:"interests" binding:"required,dive"`
}