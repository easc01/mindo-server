package youtubevideorepository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/google/uuid"
)

func BatchInsertYoutubeVideos(
	videos []dto.VideoMiniDTO,
	topicId uuid.UUID,
	userId uuid.UUID,
) ([]dto.VideoDataDTO, error) {
	var placeholders []string
	var values []interface{}
	expiry := time.Now().Add(time.Hour * 24)

	// Construct placeholders and values
	for i, video := range videos {
		placeholders = append(
			placeholders,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*8+1,
				i*8+2,
				i*8+3,
				i*8+4,
				i*8+5,
				i*8+6,
				i*8+7,
				i*8+8,
			),
		)
		values = append(
			values,
			topicId,
			video.Title,
			video.ChannelTitle,
			video.ThumbnailURL,
			video.VideoID,
			video.VideoDate,
			expiry,
			userId,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO youtube_video (topic_id, title, channel_title, thumbnail_url, video_id, video_date, expiry_at, updated_by)
		VALUES %s
		RETURNING id, topic_id, title, channel_title, thumbnail_url, video_id, video_date, expiry_at, created_at, updated_at, updated_by
	`, strings.Join(placeholders, ", "))

	// Execute query in transaction
	rows, err := db.DB.QueryContext(context.Background(), query, values...)
	if err != nil {
		logger.Log.Errorf("failed to insert topics, %s", err.Error())
		return nil, err
	}

	logger.Log.Infof(
		"inserted %d youtube videos during batch insert to topic id %s",
		len(videos),
		topicId,
	)
	defer rows.Close()

	var insertedVideos []dto.VideoDataDTO
	for rows.Next() {
		var video dto.VideoDataDTO
		if err := rows.Scan(
			&video.ID,
			&video.TopicID,
			&video.Title,
			&video.ChannelTitle,
			&video.ThumbnailURL,
			&video.VideoID,
			&video.VideoDate,
			&video.ExpiryAt,
			&video.CreatedAt,
			&video.UpdatedAt,
			&video.UpdatedBy,
		); err != nil {
			logger.Log.Errorf("failed to scan inserted video, %s", err.Error())
			return nil, err
		}
		insertedVideos = append(insertedVideos, video)
	}

	// Handle any row iteration errors
	if err := rows.Err(); err != nil {
		logger.Log.Errorf("failed during row iteration, %s", err.Error())
		return nil, err
	}

	return insertedVideos, nil
}
