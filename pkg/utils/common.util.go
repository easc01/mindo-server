package utils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"time"
	"unicode"

	"github.com/google/uuid"
)

var adjectives = []string{
	"Swift", "Clever", "Brave", "Witty", "Calm", "Mighty", "Happy", "Silent", "Lucky", "Bold",
}

var nouns = []string{
	"Fox", "Eagle", "Tiger", "Panda", "Wolf", "Hawk", "Bear", "Lion", "Otter", "Shark",
}

func GetSQLNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != Blank}
}

// Converts a map of SQL result to camelCase and removes Valid fields
func ParseSQLResponse(input interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	// If input is a pointer to struct, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := ToCamelCase(fieldType.Name)

		// Unwrap based on type
		switch field.Interface().(type) {
		case sql.NullString:
			ns := field.Interface().(sql.NullString)
			if ns.Valid {
				result[fieldName] = ns.String
			} else {
				result[fieldName] = ""
			}
		case sql.NullInt64:
			ni := field.Interface().(sql.NullInt64)
			if ni.Valid {
				result[fieldName] = ni.Int64
			} else {
				result[fieldName] = nil
			}
		case sql.NullTime:
			nt := field.Interface().(sql.NullTime)
			if nt.Valid {
				result[fieldName] = nt.Time.Format(time.RFC3339)
			} else {
				result[fieldName] = ""
			}
		default:
			// For regular types, assign directly
			result[fieldName] = field.Interface()
		}
	}

	return result, nil
}

// Converts a string to camelCase format
func ToCamelCase(str string) string {
	// Convert the first letter to lower case
	result := ""
	upperNext := false
	for i, runeValue := range str {
		if i == 0 {
			// First character to lowercase
			result += string(unicode.ToLower(runeValue))
		} else if runeValue == '_' || runeValue == '-' {
			// Skip underscores and hyphens
			upperNext = true
		} else {
			// Capitalize the next letter after an underscore
			if upperNext {
				result += string(unicode.ToUpper(runeValue))
				upperNext = false
			} else {
				result += string(runeValue)
			}
		}
	}
	return result
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
