package userservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/easc01/mindo-server/internal/models"
	authservice "github.com/easc01/mindo-server/internal/services/auth_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/encrypt"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateNewAdminUser(newUserData *dto.NewAdminUserParams) (dto.AdminUserDataDTO, error) {
	userCreationContext := context.Background()

	hashPwd, hashErr := encrypt.HashPassword(newUserData.Password)
	if hashErr != nil {
		logger.Log.Errorf("failed to hash password, %s", hashErr)
		return dto.AdminUserDataDTO{}, hashErr
	}

	tx, err := db.DB.BeginTx(userCreationContext, nil)
	if err != nil {
		logger.Log.Errorf("failed to init a transaction, %s", err)
		return dto.AdminUserDataDTO{}, err
	}

	qtx := db.Queries.WithTx(tx)
	newUserID := uuid.New()

	user, userErr := qtx.CreateNewUser(userCreationContext, models.CreateNewUserParams{
		ID:       newUserID,
		UserType: models.UserTypeAdminUser,
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
			PasswordHash: util.GetSQLNullString(hashPwd),
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
		UserType:    user.UserType,
	}, nil
}

func AdminSignIn(
	c *gin.Context,
	adminData *dto.AdminSignInParams,
) (dto.AdminUserDataDTO, int, error) {
	adminUser, adminUserErr := db.Queries.GetAdminUserByEmail(
		c,
		util.GetSQLNullString(adminData.Email),
	)

	// check if admin exists
	if adminUserErr != nil {
		if errors.Is(adminUserErr, sql.ErrNoRows) {
			logger.Log.Errorf(
				"admin user of email: %s not found %s",
				adminData.Email,
				adminUserErr,
			)
			return dto.AdminUserDataDTO{}, http.StatusNotFound, fmt.Errorf(
				"admin of email %s not found",
				adminData.Email,
			)
		}

		logger.Log.Infof("failed to fetch admin user %s", adminUserErr)
		return dto.AdminUserDataDTO{}, http.StatusInternalServerError, adminUserErr
	}

	// check the password
	dbHashPassword := adminUser.PasswordHash.String
	isPasswordValid := encrypt.CheckPasswordHash(adminData.Password, dbHashPassword)

	if !isPasswordValid {
		logger.Log.Errorf(
			"incorrect password attempt for admin user ID: %s",
			adminUser.UserID,
		)
		return dto.AdminUserDataDTO{}, http.StatusForbidden, fmt.Errorf(
			message.IncorrectPassword,
		)
	}

	// issue tokens
	accessToken, tokenErr := authservice.IssueAuthTokens(
		c,
		adminUser.UserID,
		models.UserTypeAdminUser,
	)
	if tokenErr != nil {
		logger.Log.Errorf(
			"failed to issue auth tokens for user id: %s, %s",
			adminUser.UserID.String(),
			tokenErr.Error(),
		)
		return dto.AdminUserDataDTO{}, http.StatusInternalServerError, tokenErr
	}

	go db.Queries.UpdateAdminUserLastLoginByUserId(c, adminUser.UserID)

	return dto.AdminUserDataDTO{
		AccessToken: accessToken,
		UserID:      adminUser.UserID,
		UserType:    adminUser.UserType,
		Name:        adminUser.Name.String,
		Email:       adminUser.Email.String,
		LastLoginAt: adminUser.LastLoginAt.Time,
		UpdatedAt:   adminUser.UpdatedAt.Time,
		CreatedAt:   adminUser.CreatedAt.Time,
		UpdatedBy:   adminUser.UpdatedBy.UUID,
	}, http.StatusAccepted, nil
}

func GetAdminUserByUserID(id uuid.UUID) (dto.AdminUserDataDTO, int, error) {
	adminUser, adminUserErr := db.Queries.GetAdminUserByUserID(context.Background(), id)

	if adminUserErr != nil {
		if adminUserErr == sql.ErrNoRows {
			logger.Log.Errorf("admin of id %s not found, %s", id, adminUserErr)
			return dto.AdminUserDataDTO{}, http.StatusNotFound, fmt.Errorf(message.UserNotFound)
		}

		logger.Log.Errorf("failed to get admin of id %s, %s", id, adminUserErr)
		return dto.AdminUserDataDTO{}, http.StatusInternalServerError, fmt.Errorf(
			message.SomethingWentWrong,
		)
	}

	return dto.AdminUserDataDTO{
		UserID:      adminUser.UserID,
		Name:        adminUser.Name.String,
		Email:       adminUser.Email.String,
		LastLoginAt: adminUser.LastLoginAt.Time,
		UpdatedAt:   adminUser.UpdatedAt.Time,
		CreatedAt:   adminUser.CreatedAt.Time,
		UpdatedBy:   adminUser.UpdatedBy.UUID,
		UserType:    adminUser.UserType,
	}, http.StatusAccepted, nil
}
