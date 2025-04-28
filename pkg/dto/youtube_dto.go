package dto

import "time"

type YouTubeSearchResponse struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []YouTubeSearchItem `json:"items"`
}

type YouTubeSearchItem struct {
	Kind    string         `json:"kind"`
	Etag    string         `json:"etag"`
	ID      YouTubeVideoID `json:"id"`
	Snippet YouTubeSnippet `json:"snippet"`
}

type YouTubeVideoID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type YouTubeSnippet struct {
	PublishedAt          time.Time         `json:"publishedAt"`
	ChannelID            string            `json:"channelId"`
	Title                string            `json:"title"`
	Description          string            `json:"description"`
	Thumbnails           YouTubeThumbnails `json:"thumbnails"`
	ChannelTitle         string            `json:"channelTitle"`
	LiveBroadcastContent string            `json:"liveBroadcastContent"`
	PublishTime          time.Time         `json:"publishTime"`
}

type YouTubeThumbnails struct {
	Default YouTubeThumbnail `json:"default"`
	Medium  YouTubeThumbnail `json:"medium"`
	High    YouTubeThumbnail `json:"high"`
}

type YouTubeThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
