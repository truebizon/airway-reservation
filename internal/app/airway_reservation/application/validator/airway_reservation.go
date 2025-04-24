package validator

import (
	"airway-reservation/internal/app/airway_reservation/application/converter"
	"airway-reservation/internal/pkg/logger"
	"airway-reservation/internal/pkg/myvalidator/baseValidator"
	"airway-reservation/internal/pkg/value"
)

type airwayReservation struct{}

func NewAirwayReservation() *airwayReservation {
	return &airwayReservation{}
}

func (airwayReservation) FetchAirwayReservationsByOperatorID(operatorID string) error {
	type valid struct {
		OperatorID string `json:"OperatorId" validate:"required,model-id"`
	}
	v := valid{
		OperatorID: operatorID,
	}
	validate, err := baseValidator.New()
	if err != nil {
		return err
	}
	err = validate.Struct(v)
	if err != nil {
		return baseValidator.CustomErrorMessage(err)
	}
	return nil
}

func (airwayReservation) GetRequest(airwayReservationID string) error {
	type valid struct {
		AirwayReservationID string `json:"airwayReservationId" validate:"required,model-id"`
	}
	v := valid{
		AirwayReservationID: airwayReservationID,
	}
	validate, err := baseValidator.New()
	if err != nil {
		return err
	}
	err = validate.Struct(v)
	if err != nil {
		return baseValidator.CustomErrorMessage(err)
	}
	return nil
}
func (airwayReservation) RegisterRequest(req *converter.CreateAirwayReservationRequest) error {
	type validAirwaySection struct {
		AirwaySectionID string `json:"airwaySectionId" validate:"required,model-id"`
		StartAt         string `json:"startAt" validate:"required,datetime"`
		EndAt           string `json:"endAt" validate:"required,datetime"`
	}
	type valid struct {
		OperatorID     string               `json:"operatorId" validate:"required,model-id"`
		AirwaySections []validAirwaySection `json:"airwaySections" validate:"gt=0,dive,len=1,dive,required"`
	}
	validSections := make([]validAirwaySection, len(req.AirwaySections))
	for i, section := range req.AirwaySections {
		validSections[i] = validAirwaySection{
			AirwaySectionID: section.AirwaySectionID,
			StartAt:         section.StartAt,
			EndAt:           section.EndAt,
		}
	}
	v := valid{
		OperatorID:     req.OperatorID,
		AirwaySections: validSections,
	}

	logger.StructToJson("validairwayReservation", v)

	validate, err := baseValidator.New()
	if err != nil {
		return err
	}
	err = validate.Struct(v)
	if err != nil {
		return baseValidator.CustomErrorMessage(err)
	}
	return nil
}
func (airwayReservation) UpdateRequest(req *converter.UpdateAirwayReservationRequest) error {
	type valid struct {
		AirwayReservationID string                        `json:"airwayReservationId" validate:"required,model-id"`
		Status              value.AirwayReservationStatus `json:"status" validate:"omitempty,oneof=CANCELED RESCINDED"`
	}
	v := valid{
		AirwayReservationID: req.AirwayReservationID,
		Status:              value.AirwayReservationStatus(req.Status),
	}
	validate, err := baseValidator.New()
	if err != nil {
		return err
	}
	err = validate.Struct(v)
	if err != nil {
		return baseValidator.CustomErrorMessage(err)
	}
	return nil
}

func (airwayReservation) DeleteRequest(airwayReservationID string) error {
	type valid struct {
		ID string `json:"id" validate:"required,model-id"`
	}
	v := valid{
		ID: airwayReservationID,
	}
	validate, err := baseValidator.New()
	if err != nil {
		return err
	}
	err = validate.Struct(v)
	if err != nil {
		return baseValidator.CustomErrorMessage(err)
	}
	return nil
}
