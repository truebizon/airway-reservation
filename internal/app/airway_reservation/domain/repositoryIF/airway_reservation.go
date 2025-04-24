package repositoryIF

import (
	"airway-reservation/internal/app/airway_reservation/domain/model"
	"airway-reservation/internal/pkg/value"
)

type AirwayReservationRepositoryIF interface {
	FetchAll(airwayReservations *[]model.AirwayReservation) error
	FetchALLWithPagination(page int64) (*[]model.AirwayReservation, *model.Page, error)
	FetchByOperatorWithPagination(operatorID string, page int64) (*[]model.AirwayReservation, *model.Page, error)
	FindByID(airwayReservationID string, airwayReservation *model.AirwayReservation) error
	InsertOne(airwayReservation *model.AirwayReservation) (*model.AirwayReservation, error)
	UpdateOne(airwayReservation *model.AirwayReservation) (*model.AirwayReservation, error)
	DeleteOne(airwayReservationID value.ModelID) (value.ModelID, error)
}
