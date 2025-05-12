package aiservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
)

func GenerateRoadmaps(params []dto.GeneratePlaylistParams) ([]dto.GeneratedPlaylist, error) {
	url := "https://arbazkhan-cs-mindo-apis.hf.space/MindoSyllabusGenerator"

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	var lastErr error
	client := &http.Client{}

	for attempt := 1; attempt <= 5; attempt++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed on attempt %d: %w", attempt, err)
		} else {
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}

			logger.Log.Infof("roadmap ai service response statusCode, %d, %s", res.StatusCode, string(body[:min(100, len(body))]))

			if res.StatusCode == http.StatusOK {
				var responseJson []dto.GeneratedPlaylist
				err = json.Unmarshal(body, &responseJson)
				if err != nil {
					return nil, fmt.Errorf("failed to parse roadmap ai response: %w, body: %s", err, string(body[:min(200, len(body))]))
				}
				return responseJson, nil
			}

			lastErr = fmt.Errorf("roadmap ai service returned status code %d: %s", res.StatusCode, string(body[:min(100, len(body))]))
		}

		if attempt < 5 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			logger.Log.Warnf("attempt %d failed: %v, retrying in %v", attempt, lastErr, backoff)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("all retry attempts failed: %w", lastErr)
}
