package topicrepository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/google/uuid"
)

type GetTopicByIDWithVideosRow struct {
	ID         uuid.UUID
	Name       sql.NullString
	Number     sql.NullInt32
	PlaylistID uuid.UUID
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
	UpdatedBy  uuid.NullUUID
	Videos     []dto.VideoDataDTO
}

func GetTopicByIDWithVideos(ctx context.Context, id uuid.UUID) (GetTopicByIDWithVideosRow, error) {

	const query = `-- name: GetTopicByIDWithVideos :one
		SELECT 
			t.id,
			t.name,
			t.number,
			t.playlist_id,
			t.created_at,
			t.updated_at,
			t.updated_by,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', yv.id,
						'videoId', yv.video_id,
						'title', yv.title,
						'topicId', yv.topic_id,
						'videoPublishedAt', TO_CHAR(yv.video_date, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
						'channelTitle', yv.channel_title,
						'thumbnailUrl', yv.thumbnail_url,
						'expiryAt', TO_CHAR(yv.expiry_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
						'updatedAt', TO_CHAR(yv.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
						'createdAt', TO_CHAR(yv.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
						'updatedBy', yv.updated_by
					)
				) FILTER (WHERE yv.id IS NOT NULL),
				'[]'::json
			) AS videos
		FROM
			topic t
		LEFT JOIN 
			youtube_video yv ON t.id = yv.topic_id
		WHERE t.id = $1
		GROUP BY t.id
	`

	row := db.DB.QueryRowContext(ctx, query, id)
	var i GetTopicByIDWithVideosRow
	var videosJson string
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Number,
		&i.PlaylistID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&videosJson,
	)

	if err != nil {
		return i, err
	}

	// Unmarshal the videos JSON into the Videos field
	if err := json.Unmarshal([]byte(videosJson), &i.Videos); err != nil {
		return i, err
	}

	return i, nil
}
