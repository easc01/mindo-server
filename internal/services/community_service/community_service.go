package communityservice

import (
	"fmt"
	"net/http"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNewCommunity(
	c *gin.Context,
	req *dto.CreateCommunityDTO,
) (dto.CommunityDTO, int, error) {
	user, ok := middleware.GetUser(c)

	if !ok {
		return dto.CommunityDTO{}, http.StatusUnauthorized, fmt.Errorf(message.NullUserContext)
	}

	var userID uuid.UUID
	if user.AppUser != nil {
		userID = user.AppUser.UserID
	} else {
		userID = user.AdminUser.UserID
	}

	community, err := db.Queries.CreateNewCommunity(c, models.CreateNewCommunityParams{
		Title:        util.GetSQLNullString(req.Title),
		About:        util.GetSQLNullString(req.About),
		ThumbnailUrl: util.GetSQLNullString(req.ThumbnailUrl),
		LogoUrl:      util.GetSQLNullString(req.LogoUrl),
		UpdatedBy:    util.GetNullUUID(userID),
	})

	if err != nil {
		logger.Log.Errorf("failed to create community %s, because, %s", req.Title, err.Error())
		return dto.CommunityDTO{}, http.StatusInternalServerError, err
	}

	return dto.CommunityDTO{
		ID:           community.ID,
		Title:        community.Title.String,
		About:        community.About.String,
		ThumbnailUrl: community.ThumbnailUrl.String,
		LogoUrl:      community.LogoUrl.String,
		CreatedAt:    community.CreatedAt.Time,
		UpdatedAt:    community.UpdatedAt.Time,
		UpdatedBy:    community.UpdatedBy.UUID.String(),
	}, http.StatusCreated, nil
}

func JoinExistingCommunity(c *gin.Context, communityId string) {

}
