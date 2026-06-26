package storage

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type SQLTime struct {
	time.Time
}

func (st *SQLTime) Scan(value any) error {
	if value == nil {
		st.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		st.Time = v
	case string:
		t, err := parseSQLiteTime(v)
		if err != nil {
			return err
		}
		st.Time = t
	case []byte:
		t, err := parseSQLiteTime(string(v))
		if err != nil {
			return err
		}
		st.Time = t
	default:
		return fmt.Errorf("storage: cannot scan %T into SQLTime", value)
	}
	return nil
}

func (st SQLTime) Value() (driver.Value, error) {
	if st.Time.IsZero() {
		return nil, nil
	}
	return st.Time.Format("2006-01-02 15:04:05"), nil
}

func parseSQLiteTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		time.RFC3339,
		time.RFC3339Nano,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("storage: cannot parse time %q", s)
}

var (
	_ sql.Scanner   = (*SQLTime)(nil)
	_ driver.Valuer = SQLTime{}
)
