package value

import (
	"encoding/json"
	"fmt"
)

type AirwayReservationStatus string

const (
	AIRWAY_RESERVATON_STATUS_RESERVED  AirwayReservationStatus = "RESERVED"
	AIRWAY_RESERVATON_STATUS_CANCELED  AirwayReservationStatus = "CANCELED"
	AIRWAY_RESERVATON_STATUS_RESCINDED AirwayReservationStatus = "RESCINDED"
)

func (status AirwayReservationStatus) ToString() string {
	return string(status)
}

func (status *AirwayReservationStatus) UnmarshalJSON(data []byte) error {
	var statusString string
	if err := json.Unmarshal(data, &statusString); err != nil {
		return err
	}

	switch statusString {
	case "RESERVED", "CANCELED", "RESCINDED":
		*status = AirwayReservationStatus(statusString)
	default:
		return fmt.Errorf("invalid airway_reservation status: %s", statusString)
	}

	return nil
}
