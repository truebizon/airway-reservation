package model

import (
	"airway-reservation/internal/pkg/value"
	"time"

	"gorm.io/datatypes"
)

type AirwayReservation struct {
	ID               value.ModelID                 `json:"id" gorm:"primaryKey;size:255;default:uuid_generate_v4();"`
	ExAirwaySections datatypes.JSON                `json:"ex_airway_sections" gorm:"type:jsonb"`
	AcceptedAt       time.Time                     `json:"accepted_at"`
	AirspaceID       value.ModelID                 `json:"airspace_id"`
	ReservedBy       value.ModelID                 `json:"reserved_by"`
	Status           value.AirwayReservationStatus `json:"status"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
}

type Page struct {
	Total       int64 `json:"total"`
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`
	PerPage     int64 `json:"per_page"`
}

func (a *AirwayReservation) TableName() string {
	return "airway_reservation.airway_reservations"
}
