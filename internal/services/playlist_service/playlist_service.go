package playlistservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/easc01/mindo-server/internal/models"
	interestservice "github.com/easc01/mindo-server/internal/services/interest_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func serializeTopics(topics *[]models.Topic) *[]string {
	var serializedTopics []string

	for _, topic := range *topics {
		serializedTopics = append(serializedTopics, topic.Name.String)
	}

	return &serializedTopics
}

func ProcessPlaylistCreationByAdmin(
	c *gin.Context,
	req dto.CreatePlaylistRequest,
	adminId uuid.UUID,
) (dto.PlaylistDetailsDTO, int, error) {

	// Begin a new transaction
	tx, err := db.DB.BeginTx(c, nil)
	if err != nil {
		logger.Log.Errorf("failed to begin transaction, %s", err.Error())
		return dto.PlaylistDetailsDTO{}, http.StatusInternalServerError, err
	}

	// Track the current sequence value before performing any operations
	var originalSequenceValue int32
	err = tx.QueryRowContext(c, "SELECT last_value FROM playlist_count_seq").
		Scan(&originalSequenceValue)
	if err != nil {
		logger.Log.Errorf("failed to get current sequence value: %s", err.Error())
		_ = tx.Rollback()
		return dto.PlaylistDetailsDTO{}, http.StatusInternalServerError, err
	}

	// Ensure the transaction is committed or rolled back
	defer func() {
		if err != nil {
			// Reset the sequence back to the original value before rolling back
			_, resetErr := tx.ExecContext(c, `
				SELECT setval('playlist_count_seq', $1, false)
			`, originalSequenceValue)
			if resetErr != nil {
				logger.Log.Errorf("failed to reset sequence: %s", resetErr.Error())
			}

			// Rollback the transaction
			_ = tx.Rollback()
		} else {
			// Commit the transaction if everything went smoothly
			err = tx.Commit()
		}
	}()

	// Create the playlist
	playlist, statusCode, err := CreatePlaylist(c, req, adminId, tx)
	if err != nil {
		return dto.PlaylistDetailsDTO{}, statusCode, err
	}

	// Batch insert topics
	topics, statusCode, err := BatchInsertPlaylistTopic(c, req.Topics, adminId, playlist.ID, tx)
	if err != nil {
		return dto.PlaylistDetailsDTO{}, statusCode, err
	}

	// Serialize topics for the response
	serializedTopics := serializeTopics(&topics)

	return dto.PlaylistDetailsDTO{
		ID:           playlist.ID.String(),
		Name:         playlist.Name.String,
		Description:  playlist.Description.String,
		InterestID:   playlist.InterestID.UUID.String(),
		ThumbnailURL: playlist.ThumbnailUrl.String,
		Views:        int(playlist.Views.Int32),
		Code:         playlist.Code,
		CreatedAt:    playlist.CreatedAt.Time,
		UpdatedAt:    playlist.UpdatedAt.Time,
		UpdatedBy:    playlist.UpdatedBy.UUID.String(),
		Topics:       *serializedTopics,
	}, http.StatusCreated, nil
}

func CreatePlaylist(
	c *gin.Context,
	req dto.CreatePlaylistRequest,
	userId uuid.UUID,
	tx *sql.Tx,
) (models.Playlist, int, error) {
	var playlistCount int
	err := tx.QueryRowContext(c, "SELECT nextval('playlist_count_seq')").Scan(&playlistCount)
	if err != nil {
		logger.Log.Errorf("failed to get playlist count sequence, %s", err.Error())
		return models.Playlist{}, http.StatusInternalServerError, err
	}

	interest, intStatus, intErr := interestservice.GetInterestByName(c, req.DomainName)
	if intErr != nil {
		err = intErr
		return models.Playlist{}, intStatus, intErr
	}

	playlistParams := models.CreatePlaylistParams{
		Name:         util.GetSQLNullString(req.Name),
		Description:  util.GetSQLNullString(req.Description),
		ThumbnailUrl: util.GetSQLNullString(req.ThumbnailURL),
		Code:         util.GenerateHexCode(playlistCount),
		UpdatedBy:    util.GetNullUUID(userId),
		InterestID:   util.GetNullUUID(interest.ID),
	}

	playlist, err := db.Queries.WithTx(tx).CreatePlaylist(c, playlistParams)
	if err != nil {
		logger.Log.Errorf("failed to create playlist, %s", err.Error())
		return models.Playlist{}, http.StatusInternalServerError, err
	}

	return playlist, http.StatusCreated, nil
}

func BatchInsertPlaylistTopic(
	c *gin.Context,
	topics []string,
	userId uuid.UUID,
	playlistId uuid.UUID,
	tx *sql.Tx,
) ([]models.Topic, int, error) {
	var placeholders []string
	var values []interface{}

	// Construct placeholders and values
	for i, topic := range topics {
		placeholders = append(
			placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4),
		)
		values = append(values, topic, i+1, playlistId, userId)
	}

	query := fmt.Sprintf(`
		INSERT INTO topic (name, number, playlist_id, updated_by)
		VALUES %s
		RETURNING id, name, number, playlist_id, created_at, updated_at, updated_by
	`, strings.Join(placeholders, ", "))

	// Execute query in transaction
	rows, err := tx.QueryContext(c, query, values...)
	if err != nil {
		logger.Log.Errorf("Failed to insert topics, %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	var insertedTopics []models.Topic
	for rows.Next() {
		var topic models.Topic
		if err := rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Number,
			&topic.PlaylistID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
			&topic.UpdatedBy,
		); err != nil {
			logger.Log.Errorf("Failed to scan inserted topic, %s", err.Error())
			return nil, http.StatusInternalServerError, err
		}
		insertedTopics = append(insertedTopics, topic)
	}

	// Handle any row iteration errors
	if err := rows.Err(); err != nil {
		logger.Log.Errorf("Failed during row iteration, %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return insertedTopics, http.StatusCreated, nil
}

func GetPlaylistWithTopics(
	c *gin.Context,
	playlistID uuid.UUID,
) (dto.PlaylistDetailsDTO, int, error) {
	playlist, err := GetPlaylistWithTopicsQuery(c, playlistID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Errorf("playlist of id %s not found", playlistID)
			return dto.PlaylistDetailsDTO{}, http.StatusNotFound, fmt.Errorf(
				"playlist of id %s not found",
				playlistID,
			)
		}
		logger.Log.Errorf("failed to get playlist of id %s, %s", playlistID, err.Error())
		return dto.PlaylistDetailsDTO{}, http.StatusInternalServerError, err
	}

	// Return the playlist with topics
	return dto.PlaylistDetailsDTO{
		ID:           playlist.ID.String(),
		Name:         playlist.Name.String,
		Description:  playlist.Description.String,
		Code:         playlist.Code,
		ThumbnailURL: playlist.ThumbnailUrl.String,
		Views:        int(playlist.Views.Int32),
		CreatedAt:    playlist.CreatedAt.Time,
		UpdatedAt:    playlist.UpdatedAt.Time,
		UpdatedBy:    playlist.UpdatedBy.UUID.String(),
		Topics:       playlist.Topics,
	}, http.StatusAccepted, nil
}

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
	const getPlaylistWithTopics = `
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

	row := db.DB.QueryRowContext(ctx, getPlaylistWithTopics, id)
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
