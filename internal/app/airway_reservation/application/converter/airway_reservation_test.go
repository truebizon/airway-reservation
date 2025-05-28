package converter

import (
	"gorm.io/datatypes"
	"reflect"
	"testing"
)

func TestConvertExAirwaySectionsRoundTrip(t *testing.T) {
	req := &CreateAirwayReservationRequest{
		AirwaySections: []AirwaySection{
			{AirwaySectionID: "sec1", StartAt: "2024-01-01T00:00:00Z", EndAt: "2024-01-01T00:05:00Z"},
			{AirwaySectionID: "sec2", StartAt: "2024-01-01T00:05:00Z", EndAt: "2024-01-01T00:10:00Z"},
		},
	}
	data, err := convertExAirwaySections(req)
	if err != nil {
		t.Fatalf("convertExAirwaySections error: %v", err)
	}
	// ensure valid JSON
	if len(data) == 0 {
		t.Fatalf("expected non empty json")
	}
	sections, err := convertExAirwaySectionsFromJSON(datatypes.JSON(data))
	if err != nil {
		t.Fatalf("convertExAirwaySectionsFromJSON error: %v", err)
	}
	if !reflect.DeepEqual(req.AirwaySections, sections) {
		t.Errorf("sections mismatch: expected %#v got %#v", req.AirwaySections, sections)
	}
}
