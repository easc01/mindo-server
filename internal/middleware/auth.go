package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/easc01/mindo-server/internal/models"
	authservice "github.com/easc01/mindo-server/internal/services/auth_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
)

type contextKey string

const UserContextKey contextKey = "user"

func containsUserType(slice []models.UserType, val models.UserType) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func RequireRole(allowedRoles ...models.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constant.Authorization)
		if token == "" {
			httputil.NewErrorResponse(http.StatusUnauthorized, message.AuthHeaderRequired, message.ProvideAuthHeader).
				Send(c)
			c.Abort()
			return
		}

		claims, err := authservice.ValidateJWT(token)
		if err != nil {
			logger.Log.Errorf("invalid auth token %s", err)
			httputil.NewErrorResponse(http.StatusUnauthorized, fmt.Sprintf("invalid auth token, %s", err.Error()), err.Error()).
				Send(c)
			c.Abort()
			return
		}

		if !containsUserType(allowedRoles, claims.Role) {
			logger.Log.Errorf("invalid resource access by user Id: %s", claims.Subject)
			httputil.NewErrorResponse(http.StatusForbidden, "invalid resource access", nil).Send(c)
			c.Abort()
			return
		}

		userID := util.ConvertStringToUUID(claims.Subject)

		if claims.Role == models.UserTypeAppUser {
			appUser, err := db.Queries.GetAppUserByUserID(c.Request.Context(), userID)
			if err != nil {
				handleUserFetchErr(c, claims.Subject, err)
				return
			}
			appUserContext := dto.AppUserDataDTO{
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
			}
			c.Set(string(UserContextKey), appUserContext)
			c.Next()
			return
		}

		if claims.Role == models.UserTypeAdminUser {
			adminUser, err := db.Queries.GetAdminUserByUserID(c.Request.Context(), userID)
			if err != nil {
				handleUserFetchErr(c, claims.Subject, err)
				return
			}
			adminUserContext := dto.AdminUserDataDTO{
				UserID:      adminUser.UserID,
				Name:        adminUser.Name.String,
				Email:       adminUser.Email.String,
				LastLoginAt: adminUser.LastLoginAt.Time,
				UpdatedAt:   adminUser.UpdatedAt.Time,
				CreatedAt:   adminUser.CreatedAt.Time,
				UpdatedBy:   adminUser.UpdatedBy.UUID,
				UserType:    models.UserTypeAdminUser,
			}
			c.Set(string(UserContextKey), adminUserContext)
			c.Next()
			return
		}
	}
}

func handleUserFetchErr(c *gin.Context, subject string, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		httputil.NewErrorResponse(http.StatusNotFound, message.UserNotFound, err.Error()).Send(c)
	} else {
		httputil.NewErrorResponse(http.StatusInternalServerError, message.SomethingWentWrong, err.Error()).Send(c)
	}
	logger.Log.Errorf("failed to get user by id: %s, %s", subject, err)
	c.Abort()
}

type UserContextUnion struct {
	AppUser   *dto.AppUserDataDTO
	AdminUser *dto.AdminUserDataDTO
}

func GetUser(ctx *gin.Context) (UserContextUnion, bool) {
	user, ok := ctx.Get(string(UserContextKey))
	if !ok {
		return UserContextUnion{}, false
	}

	switch u := user.(type) {
	case dto.AppUserDataDTO:
		return UserContextUnion{AppUser: &u}, true
	case dto.AdminUserDataDTO:
		return UserContextUnion{AdminUser: &u}, true
	default:
		return UserContextUnion{}, false
	}
}
