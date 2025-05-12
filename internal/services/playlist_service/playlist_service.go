package playlistservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	playlistrepository "github.com/easc01/mindo-server/internal/repository/playlist_repository"
	topicrepository "github.com/easc01/mindo-server/internal/repository/topic_repository"
	youtubevideorepository "github.com/easc01/mindo-server/internal/repository/youtube_video_repository"
	aiservice "github.com/easc01/mindo-server/internal/services/ai_service"
	interestservice "github.com/easc01/mindo-server/internal/services/interest_service"
	youtubeservice "github.com/easc01/mindo-server/internal/services/youtube_service"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/message"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func serializeTopics(topics *[]models.Topic) *[]dto.TopicsMiniDTO {
	var serializedTopics []dto.TopicsMiniDTO

	for _, topic := range *topics {
		serializedTopics = append(serializedTopics, dto.TopicsMiniDTO{
			Id:          topic.ID.String(),
			Name:        topic.Name.String,
			VideoID:     "",
			TopicNumber: int(topic.Number.Int32),
		})
	}

	return &serializedTopics
}

func ProcessPlaylistCreation(
	c *gin.Context,
	req dto.CreatePlaylistRequest,
	userId uuid.UUID,
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
	playlist, statusCode, err := CreatePlaylist(c, req, userId, tx)
	if err != nil {
		return dto.PlaylistDetailsDTO{}, statusCode, err
	}

	// Batch insert topics
	topics, statusCode, err := BatchInsertPlaylistTopic(c, req.Topics, userId, playlist.ID, tx)
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
		IsAIGen:      playlist.IsAiGen,
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
		IsAiGen:      req.IsAIGen,
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
	playlist, err := playlistrepository.GetPlaylistWithTopicsQuery(c, playlistID)

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

	// Clone necessary data (user)
	user, ok := middleware.GetUser(c)
	if ok && user.AppUser != nil {
		appUser := user.AppUser
		go func(appUserID uuid.UUID, playlistID uuid.UUID) {
			ctx := context.Background()

			// Update views
			if err := db.Queries.UpdatePlaylistViewCountById(ctx, models.UpdatePlaylistViewCountByIdParams{
				ID: playlistID,
				Views: sql.NullInt32{
					Int32: 1,
					Valid: true,
				},
			}); err != nil {
				logger.Log.Errorf("failed to update playlist views: %v", err)
			}

			// Create user_playlist
			if _, err := db.Queries.CreateUserPlaylist(ctx, models.CreateUserPlaylistParams{
				UserID:     appUserID,
				PlaylistID: playlistID,
				UpdatedBy:  util.GetNullUUID(appUserID),
			}); err != nil {
				logger.Log.Errorf("failed to create user_playlist: %v", err)
			}
		}(appUser.UserID, playlistID)
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
		IsAIGen:      playlist.IsAIGen,
		Topics:       playlist.Topics,
	}, http.StatusAccepted, nil
}

func GetAllPlaylistPreviews(
	c *gin.Context,
	searchTag string,
) ([]dto.PlaylistPreviewDTO, int, error) {
	playlists, err := db.Queries.GetAllPlaylistsPreviews(c, searchTag)

	if err != nil {
		logger.Log.Error("failed to get playlist previews")
		return []dto.PlaylistPreviewDTO{}, http.StatusInternalServerError, err
	}

	if playlists == nil {
		return []dto.PlaylistPreviewDTO{}, http.StatusAccepted, nil
	}

	var serializedPlaylist []dto.PlaylistPreviewDTO
	for _, playlist := range playlists {
		serializedPlaylist = append(serializedPlaylist, dto.PlaylistPreviewDTO{
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
			IsAIGen:      playlist.IsAiGen,
			TopicsCount:  int(playlist.TopicsCount.(int64)),
		})
	}

	return serializedPlaylist, http.StatusAccepted, nil
}

func GetVideosByTopicId(
	c *gin.Context,
	topicId uuid.UUID,
	videoId string,
) (dto.GroupedVideoDataResponse, int, error) {
	topic, err := topicrepository.GetTopicByIDWithVideos(c, topicId)
	videos := topic.Videos

	if err != nil {
		logger.Log.Errorf("failed to get yt videos by topic id %s, %s", topicId, err.Error())
		return dto.GroupedVideoDataResponse{}, http.StatusInternalServerError, err
	}

	// no videos found in db, search and save new ones
	if len(videos) == 0 {
		newVideos, err := FetchAndSaveNewVideos(c, topic, topic.PlaylistName.String)
		if err != nil {
			return dto.GroupedVideoDataResponse{}, http.StatusInternalServerError, err
		}
		return GroupVideos(newVideos, videoId), http.StatusCreated, nil
	}

	return GroupVideos(videos, videoId), http.StatusAccepted, nil
}

func FetchAndSaveNewVideos(
	c *gin.Context,
	topic topicrepository.GetTopicByIDWithVideosRow,
	playlistName string,
) ([]dto.VideoDataDTO, error) {

	topicId := topic.ID
	topicName := topic.Name.String
	user, ok := middleware.GetUser(c)

	if !ok {
		return []dto.VideoDataDTO{}, fmt.Errorf(message.NullUserContext)
	}

	var userID uuid.UUID
	if user.AppUser != nil {
		userID = user.AppUser.UserID
	} else {
		userID = user.AdminUser.UserID
	}

	videos, err := youtubeservice.SearchVideosByTopic(
		fmt.Sprintf("%s in %s", topicName, playlistName),
		10,
	)

	if err != nil {
		logger.Log.Errorf("failed to search %s on youtube, %s", topicName, err.Error())
		return []dto.VideoDataDTO{}, err
	}

	savedVideos, err := youtubevideorepository.BatchInsertYoutubeVideos(videos, topicId, userID)

	if err != nil {
		logger.Log.Errorf("failed to batch insert youtube videos of %s, %s", topicName, err.Error())
		return []dto.VideoDataDTO{}, err
	}

	return savedVideos, nil
}

func GroupVideos(videos []dto.VideoDataDTO, videoID string) dto.GroupedVideoDataResponse {
	var result dto.GroupedVideoDataResponse

	if len(videos) == 0 {
		return result
	}

	var firstVideo *dto.VideoDataDTO

	for _, video := range videos {
		if video.VideoID == videoID {
			firstVideo = &video
			break
		}
	}

	if firstVideo == nil {
		firstVideo = &videos[0]
	}

	result.Video = *firstVideo

	for _, video := range videos {
		if video.VideoID != firstVideo.VideoID {
			result.MoreVideos = append(result.MoreVideos, video)
		}
	}

	return result
}

func GenerateAndSavePlaylist(
	c *gin.Context,
	playlistTitle string,
) (dto.PlaylistDetailsDTO, error) {
	generatedPlaylists, err := aiservice.GenerateRoadmaps([]dto.GeneratePlaylistParams{
		{
			Title: playlistTitle,
		},
	})

	if err != nil && len(generatedPlaylists) == 0 {
		logger.Log.Errorf(
			"failed to generate playlist %s, because %s, generated playlists: %v",
			playlistTitle,
			err.Error(),
			generatedPlaylists,
		)
		return dto.PlaylistDetailsDTO{}, err
	}

	user, _ := middleware.GetUser(c)
	playlistData := generatedPlaylists[0]

	savedPlaylistData, _, err := ProcessPlaylistCreation(c, dto.CreatePlaylistRequest{
		Name:         playlistData.Title,
		Description:  playlistData.Description,
		DomainName:   "None",
		ThumbnailURL: "/example.com",
		Topics:       playlistData.Topics,
		IsAIGen:      true,
	}, user.AppUser.UserID)

	if err != nil {
		logger.Log.Errorf("failed to save generated playlist, %s", err.Error())
		return dto.PlaylistDetailsDTO{}, err
	}

	return savedPlaylistData, nil
}
