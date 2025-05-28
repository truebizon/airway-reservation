package value

import (
	"airway-reservation/internal/pkg/constant"
	"testing"
	"time"
)

func TestTimestampToString(t *testing.T) {
	tt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	ts := NewTimestampFromTime(tt)
	got := ts.ToString()
	expect := tt.Format(constant.TIME_MICRO_FORMAT)
	if got != expect {
		t.Errorf("expected %s got %s", expect, got)
	}
}

func TestNewTimestampFromString(t *testing.T) {
	str := "2024-01-01T10:00:00Z"
	ts, err := NewTimestampFromString(str)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.ToTime().UTC().Format(time.RFC3339) != str {
		t.Errorf("unexpected time: %s", ts.ToTime().Format(time.RFC3339))
	}
}

func TestNewTimestampFromStringInvalid(t *testing.T) {
	_, err := NewTimestampFromString("invalid")
	if err == nil {
		t.Fatalf("expected error")
	}
}
