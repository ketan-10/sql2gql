package internal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// https://github.com/99designs/gqlgen/issues/1428

func UnmarshalDatetime(v interface{}) (time.Time, error) {
	if str, ok := v.(string); ok {
		layout := "2006-01-02 15:04:05"
		result, err := time.ParseInLocation(layout, str, time.Local)
		if err == nil {
			return result, err
		}

		layout = "2006-01-02"
		result, err = time.ParseInLocation(layout, str, time.Local)
		if err == nil {
			return result, err
		}

		layout = "15:04:05"
		return time.ParseInLocation(layout, str, time.Local)
	}
	return time.Time{}, errors.New("ErrCannotParseDate")
}

func MarshalDatetime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, `"`+t.Format("2006-01-02 15:04:05")+`"`)
	})
}

func UnmarshalIntBool(v interface{}) (int8, error) {
	if value, ok := v.(bool); ok {
		if value {
			return 1, nil
		}
		return 0, nil
	}
	return 0, errors.New("ErrCannotParseBool")
}

func MarshalIntBool(v int8) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v == 1 {
			io.WriteString(w, "true")
		} else {
			io.WriteString(w, "false")
		}
	})
}

func UnmarshalNullInt64(v interface{}) (sql.NullInt64, error) {
	nullInt64 := sql.NullInt64{}
	if v == nil {
		return nullInt64, nil
	}
	value, err := graphql.UnmarshalInt(v)
	nullInt64.Int64 = int64(value)
	nullInt64.Valid = err == nil

	return nullInt64, err
}

func MarshalNullInt64(v sql.NullInt64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.Valid {
			io.WriteString(w, fmt.Sprint(v.Int64))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullFloat64(v interface{}) (sql.NullFloat64, error) {
	nullFloat64 := sql.NullFloat64{}
	if v == nil {
		return nullFloat64, nil
	}
	value, err := graphql.UnmarshalFloat(v)
	nullFloat64.Float64 = value
	nullFloat64.Valid = err == nil

	return nullFloat64, err
}

func MarshalNullFloat64(v sql.NullFloat64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.Valid {
			io.WriteString(w, fmt.Sprint(v.Float64))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullString(v interface{}) (sql.NullString, error) {
	nullString := sql.NullString{}
	if v == nil {
		return nullString, nil
	}
	value, err := graphql.UnmarshalString(v)
	nullString.String = value
	nullString.Valid = err == nil

	return nullString, err
}

func MarshalNullString(v sql.NullString) graphql.Marshaler {
	if v.Valid {
		return graphql.MarshalString(v.String)
	} else {
		return graphql.WriterFunc(func(w io.Writer) {
			io.WriteString(w, "null")
		})
	}
}

func UnmarshalNullBool(v interface{}) (sql.NullBool, error) {
	nullBool := sql.NullBool{}
	if v == nil {
		return nullBool, nil
	}
	value, err := graphql.UnmarshalBoolean(v)
	nullBool.Bool = value
	nullBool.Valid = err == nil

	return nullBool, err
}

func MarshalNullBool(v sql.NullBool) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if v.Valid {
			io.WriteString(w, fmt.Sprint(v.Bool))
		} else {
			io.WriteString(w, "null")
		}
	})
}

func MarshalNullTime(t sql.NullTime) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if t.Valid {
			io.WriteString(w, `"`+t.Time.Format("2006-01-02 15:04:05")+`"`)
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalNullTime(v interface{}) (sql.NullTime, error) {
	nt := sql.NullTime{}
	if str, ok := v.(string); ok {
		t, err := UnmarshalDatetime(str)
		if err == nil {
			nt.Time = t
			nt.Valid = true
		}
		return nt, err
	}
	return nt, errors.New("ErrCannotParseDate")
}

func MarshalMap(t map[string]interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if bytes, err := json.Marshal(t); err == nil {
			w.Write(bytes)
		} else {
			io.WriteString(w, "null")
		}
	})
}

func UnmarshalMap(v interface{}) (map[string]interface{}, error) {
	var nt map[string]interface{}
	if str, ok := v.(string); ok {
		err := json.Unmarshal([]byte(str), &nt)
		return nt, err
	}
	return nt, errors.New("ErrCannotParseMap")
}
