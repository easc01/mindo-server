package userservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/internal/models"
	authservice "github.com/easc01/mindo-server/internal/services/auth_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

func GoogleAuthService(
	c *gin.Context,
	googleReq *dto.TokenDTO,
) (dto.AppUserDataDTO, int, error) {

	payload, payloadErr := idtoken.Validate(
		c,
		googleReq.AccessToken,
		config.GetConfig().GoogleClientId,
	)
	if payloadErr != nil {
		logger.Log.Errorf(
			"invalid app user token, %s, for token %s",
			payloadErr,
			googleReq.AccessToken,
		)
		return dto.AppUserDataDTO{}, http.StatusBadRequest, payloadErr
	}

	name, _ := payload.Claims["name"].(string)
	email, _ := payload.Claims["email"].(string)

	appUserParams := dto.NewAppUserParams{
		Name:          name,
		Email:         email,
		OauthClientID: payload.Subject,
		Username:      util.GenerateUsername(),
		Mobile:        constant.Blank,
	}

	// Check if appUser exists by oauthclientId and update last login
	appUser, appUserErr := db.Queries.UpdateAppUserLastLoginAtByOAuthClientID(
		c,
		util.GetSQLNullString(appUserParams.OauthClientID),
	)

	if appUserErr != nil {
		if errors.Is(appUserErr, sql.ErrNoRows) {
			// Create new user
			newAppUser, newAppUserErr := CreateNewAppUser(appUserParams)
			if newAppUserErr != nil {
				logger.Log.Errorf(
					"failed to create app user %s for email %s oauth client id %s",
					newAppUserErr,
					appUserParams.Email,
					appUserParams.OauthClientID,
				)
				return dto.AppUserDataDTO{}, http.StatusInternalServerError, newAppUserErr
			}

			logger.Log.Infof("new app user created %s", newAppUser.UserID)
			return newAppUser, http.StatusCreated, nil
		}

		// Log unexpected DB errors
		logger.Log.Errorf(
			"failed to update last login %s for oauth client id %s",
			appUserErr,
			appUserParams.OauthClientID,
		)
		return dto.AppUserDataDTO{}, http.StatusInternalServerError, appUserErr
	}

	// create tokens
	accessToken, tokenErr := authservice.IssueAuthTokens(c, appUser.UserID, models.UserTypeAppUser)
	if tokenErr != nil {
		logger.Log.Errorf(
			"failed to issue auth tokens for user id: %s, %s",
			appUser.UserID.String(),
			tokenErr.Error(),
		)
		return dto.AppUserDataDTO{}, http.StatusInternalServerError, tokenErr
	}

	return dto.AppUserDataDTO{
		AccessToken:       accessToken,
		UserID:            appUser.UserID,
		Username:          appUser.Username.String,
		ProfilePictureUrl: appUser.ProfilePictureUrl.String,
		OauthClientID:     appUser.OauthClientID.String,
		Bio:               appUser.Bio.String,
		Name:              appUser.Name.String,
		Mobile:            appUser.Mobile.String,
		Email:             appUser.Email.String,
		LastLoginAt:       appUser.LastLoginAt.Time,
		UpdatedAt:         appUser.UpdatedAt.Time,
		CreatedAt:         appUser.CreatedAt.Time,
		UpdatedBy:         appUser.UpdatedBy.UUID,
		UserType:          models.UserTypeAppUser,
	}, http.StatusAccepted, nil
}

func CreateNewAppUser(newUserData dto.NewAppUserParams) (dto.AppUserDataDTO, error) {
	userCreationContext := context.Background()

	tx, err := db.DB.BeginTx(userCreationContext, nil)
	if err != nil {
		logger.Log.Errorf("failed to init a transaction, %s", err)
		return dto.AppUserDataDTO{}, err
	}

	qtx := db.Queries.WithTx(tx)
	newUserID := uuid.New()

	user, userErr := qtx.CreateNewUser(userCreationContext, models.CreateNewUserParams{
		ID:       newUserID,
		UserType: models.UserTypeAppUser,
		UpdatedBy: uuid.NullUUID{
			UUID:  newUserID,
			Valid: true,
		},
	})

	if userErr != nil {
		tx.Rollback()
		logger.Log.Errorf("failed to create new user of user_id %s, due to %s", newUserID, userErr)
		return dto.AppUserDataDTO{}, userErr
	}

	appUser, appUserErr := qtx.CreateNewAppUser(userCreationContext, models.CreateNewAppUserParams{
		UserID:        newUserID,
		OauthClientID: util.GetSQLNullString(newUserData.OauthClientID),
		Name:          util.GetSQLNullString(newUserData.Name),
		Username:      util.GetSQLNullString(newUserData.Username),
		Email:         util.GetSQLNullString(newUserData.Email),
		Mobile:        util.GetSQLNullString(newUserData.Mobile),
		PasswordHash:  sql.NullString{String: constant.Blank, Valid: false},

		UpdatedBy: uuid.NullUUID{
			UUID:  newUserID,
			Valid: true,
		},
	})

	if appUserErr != nil {
		tx.Rollback()
		logger.Log.Errorf(
			"failed to create new app_user and user of user_id %s, due to %s",
			newUserID,
			appUserErr,
		)
		return dto.AppUserDataDTO{}, appUserErr
	}

	txErr := tx.Commit()
	if txErr != nil {
		logger.Log.Errorf(
			"failed to create new app_user and user of user_id %s, due to %s",
			newUserID,
			txErr,
		)
		return dto.AppUserDataDTO{}, userErr
	}

	return dto.AppUserDataDTO{
		UserID:            appUser.UserID,
		Username:          appUser.Username.String,
		ProfilePictureUrl: appUser.ProfilePictureUrl.String,
		OauthClientID:     appUser.OauthClientID.String,
		Bio:               appUser.Bio.String,
		Name:              appUser.Name.String,
		Mobile:            appUser.Mobile.String,
		Email:             appUser.Email.String,
		LastLoginAt:       appUser.LastLoginAt.Time,
		UpdatedAt:         appUser.UpdatedAt.Time,
		CreatedAt:         appUser.CreatedAt.Time,
		UpdatedBy:         appUser.UpdatedBy.UUID,
		UserType:          user.UserType,
	}, nil
}

func GetAppUserByUserID(id uuid.UUID) (dto.AppUserDataDTO, int, error) {
	appUser, appUserErr := db.Queries.GetAppUserByUserID(context.Background(), id)

	if appUserErr != nil {
		if appUserErr == sql.ErrNoRows {
			logger.Log.Errorf("user of ID %s not found, %s", id, appUserErr)
			return dto.AppUserDataDTO{}, http.StatusNotFound, fmt.Errorf(message.UserNotFound)
		}

		logger.Log.Errorf("failed to get user of ID %s, %s", id, appUserErr)
		return dto.AppUserDataDTO{}, http.StatusInternalServerError, fmt.Errorf(
			message.SomethingWentWrong,
		)
	}

	return dto.AppUserDataDTO{
		UserID:            appUser.UserID,
		Username:          appUser.Username.String,
		ProfilePictureUrl: appUser.ProfilePictureUrl.String,
		OauthClientID:     appUser.OauthClientID.String,
		Bio:               appUser.Bio.String,
		Name:              appUser.Name.String,
		Mobile:            appUser.Mobile.String,
		Email:             appUser.Email.String,
		LastLoginAt:       appUser.LastLoginAt.Time,
		UpdatedAt:         appUser.UpdatedAt.Time,
		CreatedAt:         appUser.CreatedAt.Time,
		UpdatedBy:         appUser.UpdatedBy.UUID,
		UserType:          appUser.UserType,
	}, http.StatusAccepted, nil
}
