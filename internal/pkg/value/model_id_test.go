package value

import (
	"airway-reservation/internal/pkg/util"
	"testing"
)

func TestNewModelIDFromUUIDString(t *testing.T) {
	uuid := "a73b0f60-574c-409c-a9ab-3bebeb60dcfa"
	id, err := NewModelIDFromUUIDString(uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.ToString() != uuid {
		t.Errorf("expected %s got %s", uuid, id.ToString())
	}
}

func TestNewModelIDFromEncodedString(t *testing.T) {
	uuid := "a73b0f60-574c-409c-a9ab-3bebeb60dcfa"
	enc := util.Encode(uuid)
	id, err := NewModelIDFromEncodedString(enc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.ToString() != uuid {
		t.Errorf("expected %s got %s", uuid, id.ToString())
	}
}

func TestNewModelIDFromUUIDStringInvalid(t *testing.T) {
	_, err := NewModelIDFromUUIDString("invalid")
	if err == nil {
		t.Fatalf("expected error")
	}
}
