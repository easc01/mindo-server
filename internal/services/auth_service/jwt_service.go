package authservice

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/google/uuid"
)

type Claims struct {
	Role models.UserType `json:"role"`
	jwt.StandardClaims
}

var secretKey = []byte(string(config.GetConfig().JwtSecret))

func CreateAccessToken(id string, userId string, role models.UserType) (string, error) {
	token, err := createJWT(
		id,
		userId,
		role,
		constant.AppName,
		time.Now().Unix(),
		time.Now().Add(time.Hour*24).Unix(),
	)

	if err != nil {
		logger.Log.Errorf(
			"failed to create access token %s for userId %s, %s",
			id,
			userId,
			err.Error(),
		)
		return constant.Blank, err
	}

	return token, nil
}

func CreateRefreshTokenByUserId(userId uuid.UUID) (models.UserToken, error) {
	refreshToken := uuid.New().String() + uuid.New().String()

	userTokenParams := models.UpsertUserTokenParams{
		UserID:       userId,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(constant.Month),
		UpdatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
	}

	userToken, err := db.Queries.UpsertUserToken(context.Background(), userTokenParams)
	if err != nil {
		logger.Log.Errorf("failed to upsert user token of user_id: %s, %s", userId, err.Error())
		return models.UserToken{}, err
	}

	return userToken, nil
}

func createJWT(
	id string,
	subject string,
	role models.UserType,
	issuer string,
	issuesAt int64,
	expiresAt int64,
) (string, error) {
	claims := Claims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			Id:        id,
			Issuer:    issuer,
			Subject:   subject,
			IssuedAt:  issuesAt,
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return constant.Blank, err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if !claims.VerifyIssuer("Mindo2.0", true) {
		return nil, fmt.Errorf("invalid issuer")
	}

	return claims, nil
}
