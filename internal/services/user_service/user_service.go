package userservice

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ishantSikdar/mindo-server/internal/models"
	"github.com/ishantSikdar/mindo-server/pkg/db"
	"github.com/ishantSikdar/mindo-server/pkg/logger"
	"github.com/ishantSikdar/mindo-server/pkg/structs"
	"github.com/ishantSikdar/mindo-server/pkg/utils/util"
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
		OauthClientID: util.GetSQLNullString(newUserData.OauthClientID),
		Name:          util.GetSQLNullString(newUserData.Name),
		Username:      util.GetSQLNullString(newUserData.Username),
		Email:         util.GetSQLNullString(newUserData.Email),
		Mobile:        util.GetSQLNullString(newUserData.Mobile),
		PasswordHash:  sql.NullString{String: "", Valid: false},

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
		return structs.AppUserDataDTO{}, appUserErr
	}

	txErr := tx.Commit()
	if txErr != nil {
		logger.Log.Errorf(
			"failed to create new app_user and user of user_id %s, due to %s",
			newUserID,
			txErr,
		)
		return structs.AppUserDataDTO{}, userErr
	}

	return structs.AppUserDataDTO{
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
		UserType:          user.UserType.UserType,
	}, nil
}

func GetAppUserByUserID(id uuid.UUID) (structs.AppUserDataDTO, error) {
	appUser, appUserErr := db.Queries.GetAppUserByUserID(context.Background(), id)

	if appUserErr != nil {
		if appUserErr == sql.ErrNoRows {
			logger.Log.Errorf("user of ID %s not found, %s", id, appUserErr)
			return structs.AppUserDataDTO{}, appUserErr
		}

		logger.Log.Errorf("failed to get user of ID %s, %s", id, appUserErr)
		return structs.AppUserDataDTO{}, appUserErr
	}

	return structs.AppUserDataDTO{
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
	}, nil
}
