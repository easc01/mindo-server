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
	Color             models.Color
	JoinedCommunities []dto.CommunityDTO
	RecentPlaylists   []dto.PlaylistPreviewDTO
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
				au.color,
				au.last_login_at,
				au.created_at,
				au.updated_at,
				au.updated_by,
				(
						SELECT COALESCE(
								JSON_AGG(
										JSON_BUILD_OBJECT(
												'id', c.id,
												'title', c.title,
												'about', c.about,
												'thumbnailUrl', c.thumbnail_url,
												'logoUrl', c.logo_url,
												'updatedAt', TO_CHAR(c.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
												'createdAt', TO_CHAR(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
												'updatedBy', c.updated_by
										)
								),
								'[]'::json
						)
						FROM (
								SELECT c.*
								FROM user_joined_community ujc
								LEFT JOIN community c ON c.id = ujc.community_id
								WHERE ujc.user_id = au.user_id
								ORDER BY ujc.updated_at DESC
						) c
				) AS joined_communities,
				(
						SELECT COALESCE(
								JSON_AGG(
										JSON_BUILD_OBJECT(
												'id', p.id,
												'name', p.name,
												'description', p.description,
												'interestId', p.interest_id,
												'thumbnailUrl', p.thumbnail_url,
												'views', p.views,
												'code', p.code,
												'updatedAt', TO_CHAR(p.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
												'createdAt', TO_CHAR(p.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
												'updatedBy', p.updated_by
										)
								),
								'[]'::json
						)
						FROM (
								SELECT p.*
								FROM user_playlist up
								LEFT JOIN playlist p ON p.id = up.playlist_id
								WHERE up.user_id = au.user_id
								ORDER BY up.updated_at DESC
						) p
				) AS recent_playlists
		FROM app_user au
		JOIN "user" u ON u.id = au.user_id
		WHERE au.user_id = $1;
		`

	row := db.DB.QueryRowContext(ctx, query, userId)
	var i GetAppUserByUserIDRow

	var communitiesJSON []byte
	var playlistsJSON []byte

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
		&i.Color,
		&i.LastLoginAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&communitiesJSON,
		&playlistsJSON,
	)
	if err != nil {
		return i, err
	}
	// Unmarshal the JSON array into communities
	if err := json.Unmarshal(communitiesJSON, &i.JoinedCommunities); err != nil {
		return i, err
	}
	// Unmarshal the JSON array into recent playlists
	if err := json.Unmarshal(playlistsJSON, &i.RecentPlaylists); err != nil {
		return i, err
	}
	return i, err
}
