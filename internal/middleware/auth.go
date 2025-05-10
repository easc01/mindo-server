package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/easc01/mindo-server/internal/models"
	userrepository "github.com/easc01/mindo-server/internal/repository/user_repository"
	authservice "github.com/easc01/mindo-server/internal/services/auth_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
)

type UserContextUnion struct {
	AppUser   *dto.AppUserDataDTO
	AdminUser *dto.AdminUserDataDTO
}

func containsUserType(slice []models.UserType, val models.UserType) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func AuthenticateAndFetchUser(
	r *http.Request,
	allowedRoles ...models.UserType,
) (UserContextUnion, error) {
	token := r.Header.Get(constant.Authorization)
	if token == "" {
		return UserContextUnion{}, errors.New("authorization header required")
	}

	claims, err := authservice.ValidateJWT(token)
	if err != nil {
		return UserContextUnion{}, fmt.Errorf("invalid auth token: %w", err)
	}

	if !containsUserType(allowedRoles, claims.Role) {
		return UserContextUnion{}, fmt.Errorf("access denied for role: %s", claims.Role)
	}

	userID := util.ConvertStringToUUID(claims.Subject)

	switch claims.Role {
	case models.UserTypeAppUser:
		appUser, err := userrepository.GetAppUserByUserID(r.Context(), userID)

		if err != nil {
			return UserContextUnion{}, err
		}
		return UserContextUnion{
			AppUser: &dto.AppUserDataDTO{
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
				Color:             appUser.Color,
				JoinedCommunities: appUser.JoinedCommunities,
				RecentPlaylists:   appUser.RecentPlaylists,
			},
			AdminUser: nil,
		}, nil

	case models.UserTypeAdminUser:
		adminUser, err := db.Queries.GetAdminUserByUserID(r.Context(), userID)
		if err != nil {
			return UserContextUnion{}, err
		}
		return UserContextUnion{
			AdminUser: &dto.AdminUserDataDTO{
				UserID:      adminUser.UserID,
				Name:        adminUser.Name.String,
				Email:       adminUser.Email.String,
				LastLoginAt: adminUser.LastLoginAt.Time,
				UpdatedAt:   adminUser.UpdatedAt.Time,
				CreatedAt:   adminUser.CreatedAt.Time,
				UpdatedBy:   adminUser.UpdatedBy.UUID,
				UserType:    models.UserTypeAdminUser,
			},
			AppUser: nil,
		}, nil

	default:
		return UserContextUnion{}, fmt.Errorf("unsupported role: %s", claims.Role)
	}
}

func RequireRole(allowedRoles ...models.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, err := AuthenticateAndFetchUser(c.Request, allowedRoles...)
		if err != nil {
			networkutil.NewErrorResponse(http.StatusUnauthorized, err.Error(), nil).Send(c)
			c.Abort()
			return
		}
		c.Set(string(constant.UserContextKey), userData)
		c.Next()
	}
}

func GetUser(ctx *gin.Context) (UserContextUnion, bool) {
	user, ok := ctx.Get(string(constant.UserContextKey))
	if !ok {
		return UserContextUnion{}, false
	}

	if userUnion, ok := user.(UserContextUnion); ok {
		return userUnion, true
	}

	return UserContextUnion{}, false
}
