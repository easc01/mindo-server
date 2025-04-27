package interestservice

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
)

func UpsertIntoMasterInterest(c *gin.Context, interests []string, adminId string) (int, error) {
	var validInterests []string

	// Filter out empty strings
	for _, interest := range interests {
		if interest != "" {
			validInterests = append(validInterests, interest)
		}
	}

	// Check if any valid interests are left
	if len(validInterests) == 0 {
		return http.StatusBadRequest, fmt.Errorf("no valid interests provided")
	}

	var placeholders []string
	var values []interface{}

	// construct placeholders and values
	for i, interest := range validInterests {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		values = append(values, interest, adminId)
	}

	// construct batch query
	query := fmt.Sprintf(`
		INSERT INTO interest (name, updated_by)
		VALUES %s
		ON CONFLICT (name) DO NOTHING
	`, strings.Join(placeholders, ", "))

	// execute query
	_, upsertErr := db.DB.ExecContext(c, query, values...)
	if upsertErr != nil {
		logger.Log.Errorf("failed to upsert interests, %s", upsertErr.Error())
		return http.StatusInternalServerError, upsertErr
	}

	return http.StatusAccepted, nil
}

func GetMasterInterestList(c *gin.Context) ([]dto.GetInterestDTO, int, error) {
	interests, intErr := db.Queries.GetAllInterest(c)

	if intErr != nil {
		logger.Log.Errorf("failed to get master interests, %s", intErr.Error())
		return []dto.GetInterestDTO{}, http.StatusInternalServerError, intErr
	}

	if interests == nil {
		interests = []models.Interest{}
	}

	// serialize response
	serializedInterests := make([]dto.GetInterestDTO, len(interests))
	for i, interest := range interests {
		serializedInterests[i] = dto.GetInterestDTO{
			ID:   interest.ID.String(),
			Name: interest.Name.String,
		}
	}

	return serializedInterests, http.StatusAccepted, nil
}

func GetInterestByName(c *gin.Context, interestName string) (models.Interest, int, error) {
	interest, intErr := db.Queries.GetInterestByName(c, util.GetSQLNullString(interestName))
	if intErr != nil {
		if errors.Is(intErr, sql.ErrNoRows) {
			logger.Log.Errorf("interest of name %s not found", interestName)
			return models.Interest{}, http.StatusNotFound, fmt.Errorf(
				"interest of name %s not found",
				interestName,
			)
		}

		logger.Log.Errorf("failed to get interest of name %s", interestName)
		return models.Interest{}, http.StatusInternalServerError, intErr
	}

	return interest, http.StatusAccepted, nil
}
