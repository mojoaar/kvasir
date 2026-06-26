package storage

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestSQLTimeScanNil(t *testing.T) {
	var st SQLTime
	if err := st.Scan(nil); err != nil {
		t.Fatalf("Scan nil: %v", err)
	}
	if !st.Time.IsZero() {
		t.Error("expected zero time after scanning nil")
	}
}

func TestSQLTimeScanTimeType(t *testing.T) {
	var st SQLTime
	now := time.Now().Truncate(time.Second)
	if err := st.Scan(now); err != nil {
		t.Fatalf("Scan time.Time: %v", err)
	}
	if !st.Time.Equal(now) {
		t.Errorf("expected %v, got %v", now, st.Time)
	}
}

func TestSQLTimeScanString(t *testing.T) {
	var st SQLTime
	if err := st.Scan("2024-01-15 10:30:00"); err != nil {
		t.Fatalf("Scan string: %v", err)
	}
	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !st.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, st.Time)
	}
}

func TestSQLTimeScanByteSlice(t *testing.T) {
	var st SQLTime
	if err := st.Scan([]byte("2024-06-01 12:00:00")); err != nil {
		t.Fatalf("Scan []byte: %v", err)
	}
	expected := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	if !st.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, st.Time)
	}
}

func TestSQLTimeScanInvalidType(t *testing.T) {
	var st SQLTime
	if err := st.Scan(42); err == nil {
		t.Fatal("expected error for invalid scan type")
	}
}

func TestSQLTimeScanRFC3339(t *testing.T) {
	var st SQLTime
	if err := st.Scan("2024-01-15T10:30:00Z"); err != nil {
		t.Fatalf("Scan RFC3339: %v", err)
	}
	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !st.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, st.Time)
	}
}

func TestSQLTimeScanRFC3339Nano(t *testing.T) {
	var st SQLTime
	if err := st.Scan("2024-01-15T10:30:00.123456789Z"); err != nil {
		t.Fatalf("Scan RFC3339Nano: %v", err)
	}
	expected := time.Date(2024, 1, 15, 10, 30, 0, 123456789, time.UTC)
	if !st.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, st.Time)
	}
}

func TestSQLTimeScanMillisecond(t *testing.T) {
	var st SQLTime
	if err := st.Scan("2024-01-15T10:30:00.123Z"); err != nil {
		t.Fatalf("Scan millisecond: %v", err)
	}
	expected := time.Date(2024, 1, 15, 10, 30, 0, 123000000, time.UTC)
	if !st.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, st.Time)
	}
}

func TestSQLTimeValueZero(t *testing.T) {
	st := SQLTime{}
	val, err := st.Value()
	if err != nil {
		t.Fatalf("Value zero: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil for zero time, got %v", val)
	}
}

func TestSQLTimeValueNonZero(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	st := SQLTime{Time: now}
	val, err := st.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	expected := "2024-01-15 10:30:00"
	if val != expected {
		t.Errorf("expected %q, got %q", expected, val)
	}
}

func TestSQLTimeInterfaceCompliance(t *testing.T) {
	var _ driver.Valuer = SQLTime{}
	var _ interface{ Scan(any) error } = &SQLTime{}
}

func TestParseSQLiteTimeInvalid(t *testing.T) {
	if _, err := parseSQLiteTime("not-a-date"); err == nil {
		t.Fatal("expected error for invalid date string")
	}
}

func TestSearchEmptyQueryResultsNilSlice(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('hello', 'world')`)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	results, err := store.Search("\"nonexistent phrase that won't match\"", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if results == nil {
		t.Fatal("expected non-nil slice")
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchDeletedAtIsNil(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('test', 'hello')`)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	results, err := store.Search("hello", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].DeletedAt != nil {
		t.Error("expected nil DeletedAt")
	}
}
