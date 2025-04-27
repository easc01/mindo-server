package dto

type UpsertInterestDTO struct {
	Interests []string `json:"interests" binding:"required,dive"`
}

type GetInterestDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
