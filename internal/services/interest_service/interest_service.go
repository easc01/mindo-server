package interestservice

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/logger"
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

func GetMasterInterestList(c *gin.Context) ([]string, int, error) {
	interests, intErr := db.Queries.GetAllInterest(c)

	if intErr != nil {
		logger.Log.Errorf("failed to get master interests, %s", intErr.Error())
		return []string{}, http.StatusInternalServerError, intErr
	}

	if interests == nil {
		interests = []models.Interest{}
	}

	// serialize response
	serializedInterests := make([]string, len(interests))
	for i, interest := range interests {
		serializedInterests[i] = interest.Name.String
	}

	return serializedInterests, http.StatusAccepted, nil
}
