package userservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/easc01/mindo-server/pkg/utils/http"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/google/uuid"
)

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
		UserType:          user.UserType.UserType,
	}, nil
}

func CreateNewAdminUser(newUserData dto.NewAdminUserParams) (dto.AdminUserDataDTO, error) {
	userCreationContext := context.Background()

	tx, err := db.DB.BeginTx(userCreationContext, nil)
	if err != nil {
		logger.Log.Errorf("failed to init a transaction, %s", err)
		return dto.AdminUserDataDTO{}, err
	}

	qtx := db.Queries.WithTx(tx)
	newUserID := uuid.New()

	user, userErr := qtx.CreateNewUser(userCreationContext, models.CreateNewUserParams{
		ID: newUserID,
		UserType: models.NullUserType{
			UserType: models.UserTypeAdminUser,
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
		return dto.AdminUserDataDTO{}, userErr
	}

	adminUser, adminUserErr := qtx.CreateNewAdminUser(
		userCreationContext,
		models.CreateNewAdminUserParams{
			UserID:       newUserID,
			Name:         util.GetSQLNullString(newUserData.Name),
			Email:        util.GetSQLNullString(newUserData.Email),
			PasswordHash: sql.NullString{String: constant.Blank, Valid: false},
			UpdatedBy: uuid.NullUUID{
				UUID:  newUserID,
				Valid: true,
			},
		},
	)

	if adminUserErr != nil {
		tx.Rollback()
		logger.Log.Errorf(
			"failed to create new admin_user and user of user_id %s, due to %s",
			newUserID,
			adminUserErr,
		)
		return dto.AdminUserDataDTO{}, adminUserErr
	}

	txErr := tx.Commit()
	if txErr != nil {
		logger.Log.Errorf(
			"failed to create new admin_user and user of user_id %s, due to %s",
			newUserID,
			txErr,
		)
		return dto.AdminUserDataDTO{}, userErr
	}

	return dto.AdminUserDataDTO{
		UserID:      adminUser.UserID,
		Name:        adminUser.Name.String,
		Email:       adminUser.Email.String,
		LastLoginAt: adminUser.LastLoginAt.Time,
		UpdatedAt:   adminUser.UpdatedAt.Time,
		CreatedAt:   adminUser.CreatedAt.Time,
		UpdatedBy:   adminUser.UpdatedBy.UUID,
		UserType:    user.UserType.UserType,
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
		UserType:          models.UserTypeAppUser,
	}, http.StatusFound, nil
}
