package services

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/db"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/structs"
	"github.com/ishantSikdar/mindo-server/pkg/utils"
)

func CreateNewAppUser(newUserData structs.NewAppUserParams) (structs.AppUserDataDTO, error) {
	userCreationContext := context.Background()

	tx, err := db.DB.BeginTx(userCreationContext, nil)
	if err != nil {
		logger.Log.Errorf("failed to init a transaction, %s", err)
		return structs.AppUserDataDTO{}, err
	}

	qtx := db.Queries.WithTx(tx)
	newUserID := uuid.New()

	user, userErr := qtx.CreateNewUser(userCreationContext, models.CreateNewUserParams{
		ID: newUserID,
		UserType: models.NullUserType{
			UserType: models.UserTypeAppUser,
			Valid:    true,
		},
		UpdatedBy: uuid.NullUUID{
			UUID:  newUserID,
			Valid: true,
		},
	})

	if userErr != nil {
		tx.Rollback()
		logger.Log.Errorf("failed to create new user of user_id %s, due to %s", newUserID, userErr)
		return structs.AppUserDataDTO{}, userErr
	}

	appUser, appUserErr := qtx.CreateNewAppUser(userCreationContext, models.CreateNewAppUserParams{
		UserID:        newUserID,
		OauthClientID: utils.GetSQLNullString(newUserData.OauthClientID),
		Name:          utils.GetSQLNullString(newUserData.Name),
		Username:      utils.GetSQLNullString(newUserData.Username),
		Email:         utils.GetSQLNullString(newUserData.Email),
		Mobile:        utils.GetSQLNullString(newUserData.Mobile),
		PasswordHash:  sql.NullString{String: "", Valid: false},
		
		UpdatedBy: uuid.NullUUID{
			UUID:  newUserID,
			Valid: true,
		},
	})

	if appUserErr != nil {
		tx.Rollback()
		logger.Log.Errorf("failed to create new app_user and user of user_id %s, due to %s", newUserID, appUserErr)
		return structs.AppUserDataDTO{}, appUserErr
	}

	txErr := tx.Commit()
	if txErr != nil {
		logger.Log.Errorf("failed to create new app_user and user of user_id %s, due to %s", newUserID, txErr)
		return structs.AppUserDataDTO{}, userErr
	}

	return structs.AppUserDataDTO{
		UserID:            appUser.UserID,
		Username:          appUser.Username,
		ProfilePictureUrl: appUser.ProfilePictureUrl,
		OauthClientID:     appUser.OauthClientID,
		Bio:               appUser.Bio,
		Name:              appUser.Name,
		Mobile:            appUser.Mobile,
		Email:             appUser.Email,
		LastLoginAt:       appUser.LastLoginAt,
		UpdatedAt:         appUser.UpdatedAt,
		CreatedAt:         appUser.CreatedAt,
		UpdatedBy:         appUser.UpdatedBy,
		UserType:          user.UserType.UserType,
	}, nil
}

func GetAppUserByUserID(id uuid.UUID) (models.GetAppUserByUserIDRow, error) {
	user, userErr := db.Queries.GetAppUserByUserID(context.Background(), id)

	if userErr != nil {
		if userErr == sql.ErrNoRows {
			logger.Log.Errorf("user of ID %s not found, %s", id, userErr)
			return models.GetAppUserByUserIDRow{}, userErr
		}

		logger.Log.Errorf("failed to get user of ID %s, %s", id, userErr)
		return models.GetAppUserByUserIDRow{}, userErr
	}

	return user, nil
}
