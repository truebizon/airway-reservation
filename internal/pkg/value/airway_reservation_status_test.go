package value

import (
	"encoding/json"
	"testing"
)

func TestAirwayReservationStatusUnmarshal(t *testing.T) {
	data := []byte(`"RESERVED"`)
	var status AirwayReservationStatus
	if err := json.Unmarshal(data, &status); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != AIRWAY_RESERVATON_STATUS_RESERVED {
		t.Errorf("expected %s got %s", AIRWAY_RESERVATON_STATUS_RESERVED, status)
	}
}

func TestAirwayReservationStatusUnmarshalInvalid(t *testing.T) {
	data := []byte(`"INVALID"`)
	var status AirwayReservationStatus
	if err := json.Unmarshal(data, &status); err == nil {
		t.Fatal("expected error")
	}
}
