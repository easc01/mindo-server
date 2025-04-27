package playlistrepository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/easc01/mindo-server/pkg/db"
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
	Topics       []string
}

func GetPlaylistWithTopicsQuery(
	ctx context.Context,
	id uuid.UUID,
) (GetPlaylistWithTopicsRow, error) {
	const query = `
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
						json_agg(t.name ORDER BY t.number ASC), 
						'[]'
				) AS topics
		FROM playlist p
		LEFT JOIN topic t ON p.id = t.playlist_id
		WHERE p.id = $1
		GROUP BY p.id
	`

	row := db.DB.QueryRowContext(ctx, query, id)
	var i GetPlaylistWithTopicsRow
	var topicsJSON string
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
	if err := json.Unmarshal([]byte(topicsJSON), &i.Topics); err != nil {
		return i, err
	}

	return i, nil
}
