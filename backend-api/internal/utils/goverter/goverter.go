// Some utils for https://github.com/jmattheis/goverter codegen tool

package goverter

import (
	"database/sql"
	"time"
)

func SQLStringToPString(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

func PStringToSQLString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: *value,
		Valid:  true,
	}
}

func StringToSQLString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func SQLTimeToPTime(value sql.NullTime) *time.Time {
	if value.Valid {
		return &value.Time
	}
	return nil
}

func PTimeToSQLTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  *value,
		Valid: true,
	}
}

func TimeToTime(t time.Time) time.Time {
	return t
}
