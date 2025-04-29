package util

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/easc01/mindo-server/pkg/utils/constant"
	"github.com/google/uuid"
)

var adjectives = []string{
	"Swift", "Clever", "Brave", "Witty", "Calm", "Mighty", "Happy", "Silent", "Lucky", "Bold",
}

var nouns = []string{
	"Fox", "Eagle", "Tiger", "Panda", "Wolf", "Hawk", "Bear", "Lion", "Otter", "Shark",
}

// Returns sql.NullString for a string
func GetSQLNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != constant.Blank}
}

// Returns sql.NullString for a string
func GetNullUUID(s uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: s, Valid: true}
}

// ConvertStringToUUID converts a string to a uuid.UUID
func ConvertStringToUUID(id string) uuid.UUID {
	parsedId, _ := uuid.Parse(id)
	return parsedId
}

// Generate a random username with UUID suffix
func GenerateUsername() string {
	uid := uuid.New()
	adjective := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]

	// Combine random word with UUID suffix (first 8 characters)
	username := fmt.Sprintf("%s%s-%s", adjective, noun, uid.String()[:8])

	return username
}

func GenerateHexCode(num int) string {
	// Format the integer as a hexadecimal string with padding to ensure it's 6 characters long
	hexCode := fmt.Sprintf("%06X", num)
	return hexCode
}

func GetUUIDFromString(s string) uuid.NullUUID {
	if s != constant.Blank {
		parsedUUID, err := uuid.Parse(s)

		if err != nil {
			return uuid.NullUUID{
				Valid: false,
			}
		}

		return uuid.NullUUID{
			UUID:  parsedUUID,
			Valid: true,
		}
	}

	return uuid.NullUUID{
		Valid: false,
	}
}
