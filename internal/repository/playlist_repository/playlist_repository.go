package playlistrepository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/google/uuid"
)

type GetPlaylistWithTopicsRow struct {
	ID           uuid.UUID
	Name         sql.NullString
	Description  sql.NullString
	Code         string
	ThumbnailUrl sql.NullString
	Views        sql.NullInt32
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	UpdatedBy    uuid.NullUUID
	Topics       []dto.TopicsMiniDTO
}

func GetPlaylistWithTopicsQuery(
	ctx context.Context,
	id uuid.UUID,
) (GetPlaylistWithTopicsRow, error) {
	const query = `
		WITH ranked_videos AS (
			SELECT 
				yv.topic_id,
				yv.video_id,
				ROW_NUMBER() OVER (PARTITION BY yv.topic_id ORDER BY yv.created_at DESC) AS rn
			FROM youtube_video yv
		)
		SELECT 
				p.id, 
				p.name, 
				p.description, 
				p.code, 
				p.thumbnail_url, 
				p.views, 
				p.created_at, 
				p.updated_at, 
				p.updated_by,
				COALESCE(
					JSON_AGG(
						JSON_BUILD_OBJECT(
							'id', t.id,
							'name', t.name,
							'videoId', rv.video_id
						)
					) FILTER (WHERE t.id IS NOT NULL),
					'[]'::json
				) AS topics
		FROM playlist p
		LEFT JOIN topic t ON p.id = t.playlist_id
		LEFT JOIN ranked_videos rv ON rv.topic_id = t.id AND rv.rn = 1
		WHERE p.id = $1
		GROUP BY p.id
	`
	row := db.DB.QueryRowContext(ctx, query, id)
	var i GetPlaylistWithTopicsRow
	var topicsJSON []byte
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Code,
		&i.ThumbnailUrl,
		&i.Views,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&topicsJSON,
	)
	if err != nil {
		return i, err
	}
	// Unmarshal the JSON array into Topics
	if err := json.Unmarshal(topicsJSON, &i.Topics); err != nil {
		return i, err
	}
	return i, nil
}
