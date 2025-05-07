package userrepository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/google/uuid"
)

type GetAppUserByUserIDRow struct {
	UserID            uuid.UUID
	UserType          models.UserType
	Username          sql.NullString
	ProfilePictureUrl sql.NullString
	OauthClientID     sql.NullString
	Bio               sql.NullString
	Name              sql.NullString
	Mobile            sql.NullString
	Email             sql.NullString
	LastLoginAt       sql.NullTime
	CreatedAt         sql.NullTime
	UpdatedAt         sql.NullTime
	UpdatedBy         uuid.NullUUID
	JoinedCommunities []dto.CommunityDTO
}

func GetAppUserByUserID(
	ctx context.Context,
	userId uuid.UUID,
) (GetAppUserByUserIDRow, error) {
	const query = `-- name: GetAppUserByUserID :one
		SELECT
				u.id AS user_id,
				u.user_type,
				au.username,
				au.profile_picture_url,
				au.oauth_client_id,
				au.bio,
				au.name,
				au.mobile,
				au.email,
				au.last_login_at,
				au.created_at,
				au.updated_at,
				au.updated_by,
				COALESCE(
						JSON_AGG(
								JSON_BUILD_OBJECT(
										'id', c.id,
										'title', c.title,
										'about', c.about,
										'thumbnailUrl', c.thumbnail_url,
										'logoUrl', c.logo_url,
										'updatedAt', c.updated_at,
										'createdAt', c.created_at,
										'updatedBy', c.updated_by
								)
						) FILTER (WHERE c.id IS NOT NULL),
						'[]'::json
				) AS joined_communities
		FROM app_user au
				JOIN "user" u ON u.id = au.user_id
				LEFT JOIN user_joined_community ujc ON ujc.user_id = au.user_id
				LEFT JOIN community c ON c.id = ujc.community_id
		WHERE
				au.user_id = $1
		GROUP BY
				u.id, u.user_type,
				au.username, au.profile_picture_url, au.oauth_client_id, au.bio, au.name,
				au.mobile, au.email, au.last_login_at, au.created_at, au.updated_at, au.updated_by
		`

	row := db.DB.QueryRowContext(ctx, query, userId)
	var i GetAppUserByUserIDRow
	var communitiesJSON []byte
	err := row.Scan(
		&i.UserID,
		&i.UserType,
		&i.Username,
		&i.ProfilePictureUrl,
		&i.OauthClientID,
		&i.Bio,
		&i.Name,
		&i.Mobile,
		&i.Email,
		&i.LastLoginAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&communitiesJSON,
	)
	if err != nil {
		return i, err
	}
	// Unmarshal the JSON array into communities
	if err := json.Unmarshal(communitiesJSON, &i.JoinedCommunities); err != nil {
		return i, err
	}
	return i, err
}
