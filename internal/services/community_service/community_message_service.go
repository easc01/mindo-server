package communityservice

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/db"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/utils/util"
	"github.com/google/uuid"
)

func SaveCommunityMessage(
	c context.Context,
	communityID uuid.UUID,
	userID uuid.UUID,
	msg string,
) (models.Message, error) {
	return db.Queries.CreateMessage(c, models.CreateMessageParams{
		CommunityID: communityID,
		UserID:      userID,
		Content:     util.GetSQLNullString(msg),
		UpdatedBy:   util.GetNullUUID(userID),
	})
}

func GetMessageHistoryPage(
	c context.Context,
	communityID uuid.UUID,
	lastDate time.Time,
) ([]dto.UserMessageDTO, int, error) {
	messages, err := db.Queries.GetMessagePageByCommunityID(
		c,
		models.GetMessagePageByCommunityIDParams{
			CommunityID: communityID,
			CreatedAt: sql.NullTime{
				Valid: true,
				Time:  lastDate,
			},
		},
	)

	if err != nil {
		return []dto.UserMessageDTO{}, http.StatusInternalServerError, err
	}

	return serializeUserMessages(messages), http.StatusAccepted, nil
}

func serializeUserMessages(messages []models.GetMessagePageByCommunityIDRow) []dto.UserMessageDTO {
	var userMessages = make([]dto.UserMessageDTO, 0)
	timeFrame := 2 * time.Minute

	for _, message := range messages {

		// append current message to last group message
		if len(userMessages) > 0 {
			lastGroup := &userMessages[len(userMessages)-1]

			if message.UserID == lastGroup.UserID &&
				lastGroup.Messages[len(lastGroup.Messages)-1].Timestamp.Sub(
					message.CreatedAt.Time,
				) <= timeFrame {

				lastGroup.Messages = append(lastGroup.Messages, dto.MessageDTO{
					ID:        message.ID,
					Content:   message.Content.String,
					Timestamp: message.CreatedAt.Time,
				})

				continue
			}
		}

		userMessages = append(userMessages, dto.UserMessageDTO{
			MessageGroupID: uuid.New(),
			UserID:         message.UserID,
			Username:       message.Username.String,
			UserProfileUrl: message.ProfilePictureUrl.String,
			Name:           message.Name.String,
			Messages: []dto.MessageDTO{
				{
					ID:        message.ID,
					Content:   message.Content.String,
					Timestamp: message.CreatedAt.Time,
				},
			},
		})
	}

	for i := range userMessages {
		util.ReverseSlice(userMessages[i].Messages)
	}

	return userMessages
}
