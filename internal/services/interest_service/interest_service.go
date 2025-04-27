package interestservice

import (
	"net/http"

	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/gin-gonic/gin"
)

func GetMasterInterestList(c *gin.Context) ([]models.Interest, int, error) {
	interests, intErr := db.Queries.GetAllInterest(c)

	if intErr != nil {
		logger.Log.Errorf("failed to get master interests, %s", intErr.Error())
		return []models.Interest{}, http.StatusInternalServerError, intErr
	}

	if interests == nil {
		interests = []models.Interest{}
	}

	return interests, http.StatusAccepted, nil
}
