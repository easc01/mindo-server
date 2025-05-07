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

	tx, err := db.DB.BeginTx(c, nil)
	if err != nil {
		return dto.CommunityDTO{}, http.StatusInternalServerError, err
	}
	qtx := db.Queries.WithTx(tx)

	// Step 1: Create the community
	community, err := qtx.CreateNewCommunity(c, models.CreateNewCommunityParams{
		Title:        util.GetSQLNullString(req.Title),
		About:        util.GetSQLNullString(req.About),
		ThumbnailUrl: util.GetSQLNullString(req.ThumbnailUrl),
		LogoUrl:      util.GetSQLNullString(req.LogoUrl),
		UpdatedBy:    util.GetNullUUID(userID),
	})
	if err != nil {
		tx.Rollback()
		logger.Log.Errorf("failed to create community %s, because, %s", req.Title, err.Error())
		return dto.CommunityDTO{}, http.StatusInternalServerError, err
	}

	// Step 2: Join the community
	err = qtx.CreateNewUserJoinedCommunityById(
		c,
		models.CreateNewUserJoinedCommunityByIdParams{
			UserID:      userID,
			CommunityID: community.ID,
			UpdatedBy:   util.GetNullUUID(userID),
		},
	)
	if err != nil {
		tx.Rollback()
		logger.Log.Errorf("failed to join community %s, because %s", community.ID, err.Error())
		return dto.CommunityDTO{}, http.StatusInternalServerError, err
	}

	if err := tx.Commit(); err != nil {
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

func JoinExistingCommunity(c *gin.Context, communityId uuid.UUID) (int, error) {
	user, ok := middleware.GetUser(c)

	if user.AppUser == nil || !ok {
		return http.StatusUnauthorized, fmt.Errorf(message.NullAppUserContext)
	}

	userId := user.AppUser.UserID

	var exists = false
	for _, community := range user.AppUser.JoinedCommunities {
		if community.ID == communityId {
			exists = true
		}
	}

	if exists {
		return http.StatusConflict, fmt.Errorf("community already joined")
	}

	err := db.Queries.CreateNewUserJoinedCommunityById(
		c,
		models.CreateNewUserJoinedCommunityByIdParams{
			UserID:      userId,
			CommunityID: communityId,
			UpdatedBy:   util.GetNullUUID(userId),
		},
	)

	if err != nil {
		logger.Log.Errorf(
			"failed to create user joined community by id %s, because %s",
			communityId,
			err.Error(),
		)
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}
