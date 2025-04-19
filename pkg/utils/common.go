package utils

import "database/sql"

func GetSQLNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != Blank}
}
