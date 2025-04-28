package youtubeservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/easc01/mindo-server/internal/config"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
)

func SearchVideosByTopic(query string, maxResults int) ([]dto.VideoMiniDTO, error) {
	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?key=%s&q=%s&safeSearch=strict&type=video&videoEmbeddable=true&part=snippet&videoDuration=any&maxResults=%s",
		config.GetConfig().YoutubeAPIKey,
		url.QueryEscape(query),
		strconv.Itoa(maxResults),
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return []dto.VideoMiniDTO{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return []dto.VideoMiniDTO{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []dto.VideoMiniDTO{}, err
	}

	// Debug the actual response
	logger.Log.Debugf("youTube api response status: %d", res.StatusCode)

	// Only try to parse if we got a 200 OK
	if res.StatusCode != http.StatusOK {
		return []dto.VideoMiniDTO{}, fmt.Errorf(
			"youtube api returned status code %d: %s",
			res.StatusCode,
			string(body[:100]),
		)
	}

	var responseJson dto.YouTubeSearchResponse
	err = json.Unmarshal(body, &responseJson)
	if err != nil {
		return []dto.VideoMiniDTO{}, fmt.Errorf(
			"failed to parse YouTube response: %w, body: %s",
			err,
			string(body[:200]),
		)
	}

	return serializeYoutubeResponse(responseJson), nil
}

func serializeYoutubeResponse(response dto.YouTubeSearchResponse) []dto.VideoMiniDTO {
	var serializedVideos []dto.VideoMiniDTO

	for _, video := range response.Items {
		serializedVideos = append(serializedVideos, dto.VideoMiniDTO{
			VideoID:      video.ID.VideoID,
			Title:        video.Snippet.Title,
			VideoDate:    video.Snippet.PublishedAt,
			ChannelTitle: video.Snippet.ChannelTitle,
			ThumbnailURL: video.Snippet.Thumbnails.Default.URL,
		})
	}

	return serializedVideos
}
