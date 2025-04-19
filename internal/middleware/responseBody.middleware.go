package middleware

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ishantSikdar/mindo-server/pkg/utils"
)

// Custom response writer to capture the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	// Only write to the buffer, not to the original response writer
	return r.body.Write(b)
}

func (r *responseBodyWriter) WriteString(s string) (int, error) {
	// Only write to the buffer, not to the original response writer
	return r.body.WriteString(s)
}

// ResponseFormatter is a middleware that formats SQL responses
func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store the original writer
		originalWriter := c.Writer

		// Create a custom writer to capture the response
		writer := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// Process the request
		c.Next()

		// After the request is processed, check the response
		if writer.body.Len() > 0 {
			var data interface{}
			if err := json.Unmarshal(writer.body.Bytes(), &data); err == nil {
				// Process the data based on its type
				switch v := data.(type) {
				case map[string]interface{}:
					// Look for SQL Null* objects and transform them
					processedData := processSQLResponse(v)
					// Write the processed data to the original writer
					originalWriter.Header().Set("Content-Type", "application/json")
					originalWriter.WriteHeader(writer.Status())
					jsonData, _ := json.Marshal(processedData)
					originalWriter.Write(jsonData)
					return
				case []interface{}:
					// Handle array of objects
					processedArray := make([]interface{}, 0, len(v))
					for _, item := range v {
						if mapItem, ok := item.(map[string]interface{}); ok {
							processedArray = append(processedArray, processSQLResponse(mapItem))
						} else {
							processedArray = append(processedArray, item)
						}
					}
					// Write the processed array to the original writer
					originalWriter.Header().Set("Content-Type", "application/json")
					originalWriter.WriteHeader(writer.Status())
					jsonData, _ := json.Marshal(processedArray)
					originalWriter.Write(jsonData)
					return
				}
			}

			// If we couldn't process it, write the original response
			originalWriter.WriteHeader(writer.Status())
			originalWriter.Write(writer.body.Bytes())
		}
	}
}

// processSQLResponse processes SQL response objects and transforms them
func processSQLResponse(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		// Convert key to camelCase
		camelKey := utils.ToCamelCase(key)

		switch v := value.(type) {
		case map[string]interface{}:
			// Handle SQL Null* types
			if valid, hasValid := v["Valid"].(bool); hasValid {
				if valid {
					switch {
					case v["String"] != nil:
						result[camelKey] = v["String"]
					case v["Time"] != nil:
						result[camelKey] = v["Time"]
					case v["Int64"] != nil:
						result[camelKey] = v["Int64"]
					case v["Float64"] != nil:
						result[camelKey] = v["Float64"]
					case v["Bool"] != nil:
						result[camelKey] = v["Bool"]
					default:
						result[camelKey] = nil
					}
				} else {
					// Provide empty string for string fields, nil otherwise
					if _, ok := v["String"]; ok {
						result[camelKey] = ""
					} else {
						result[camelKey] = nil
					}
				}
			}

		case []interface{}:
			// Process array items
			processedArray := make([]interface{}, 0, len(v))
			for _, item := range v {
				if mapItem, ok := item.(map[string]interface{}); ok {
					processedArray = append(processedArray, processSQLResponse(mapItem))
				} else {
					processedArray = append(processedArray, item)
				}
			}
			result[camelKey] = processedArray

		default:
			// Regular value
			result[camelKey] = v
		}
	}

	return result
}
