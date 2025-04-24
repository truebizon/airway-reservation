package converter

import (
	"airway-reservation/internal/app/airway_reservation/domain/model"
	"airway-reservation/internal/pkg/util"
	"airway-reservation/internal/pkg/value"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

type airwayReservationConverter struct{}

func NewAirwayReservationDomain() *airwayReservationConverter {
	return &airwayReservationConverter{}
}

type CreateAirwayReservationRequest struct {
	OperatorID     string          `json:"operatorId"`
	AirwaySections []AirwaySection `json:"airwaySections"`
}

type UpdateAirwayReservationRequest struct {
	AirwayReservationID string `json:"airwayReservationId"`
	Status              string `json:"status"`
}
type AirwaySection struct {
	AirwaySectionID string `json:"airwaySectionId"`
	StartAt         string `json:"startAt"`
	EndAt           string `json:"endAt"`
}

type AirwayReservationResponse struct {
	AirwayReservationID string                        `json:"airwayReservationId"`
	OperatorID          string                        `json:"operatorId"`
	AirwaySections      []AirwaySection               `json:"airwaySections"`
	Status              value.AirwayReservationStatus `json:"status"`
	ReservedAt          string                        `json:"reservedAt"`
	UpdatedAt           string                        `json:"updatedAt"`
}

type AirwayReservationResponses struct {
	Result      []AirwayReservationResponse `json:"result"`
	Total       int64                       `json:"total"`
	CurrentPage int64                       `json:"currentPage"`
	LastPage    int64                       `json:"lastPage"`
	PerPage     int64                       `json:"perPage"`
}

type AirwayReservationMessage struct {
	EventID             string                        `json:"eventId"`
	AirwayReservationID string                        `json:"airwayReservationId"`
	OperatorID          string                        `json:"operatorId"`
	AirwaySections      []AirwaySection               `json:"airwaySections"`
	Status              value.AirwayReservationStatus `json:"status"`
	ReservedAt          string                        `json:"reservedAt"`
	UpdatedAt           string                        `json:"updatedAt"`
}

func NewCreateAirwayReservationRequest() *CreateAirwayReservationRequest {
	return &CreateAirwayReservationRequest{
		OperatorID:     "",
		AirwaySections: []AirwaySection{},
	}
}

func NewUpdateAirwayReservationRequest() *UpdateAirwayReservationRequest {
	return &UpdateAirwayReservationRequest{
		AirwayReservationID: "",
		Status:              "",
	}
}

func (airwayReservationConverter) ToUpdateAirwayReservationRequest(AirwayReservationID, status string) *UpdateAirwayReservationRequest {
	return &UpdateAirwayReservationRequest{
		AirwayReservationID: AirwayReservationID,
		Status:              status,
	}
}

func convertExAirwaySections(req *CreateAirwayReservationRequest) (datatypes.JSON, error) {
	airwaySections := req.AirwaySections
	exAirwaySections := make([]map[string]interface{}, 0, len(airwaySections))
	for _, section := range airwaySections {
		exAirwaySections = append(exAirwaySections, map[string]interface{}{
			"airway_section_id": section.AirwaySectionID,
			"start_at":          section.StartAt,
			"end_at":            section.EndAt,
		})
	}
	jsonData, err := json.Marshal(exAirwaySections)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ExAirwaySections to JSON: %v", err)
	}
	return datatypes.JSON(jsonData), nil
}
func convertExAirwaySectionsFromJSON(jsonData datatypes.JSON) ([]AirwaySection, error) {
	var exAirwaySections []map[string]interface{}
	err := json.Unmarshal(jsonData, &exAirwaySections)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ExAirwaySections from JSON: %v", err)
	}
	airwaySections := make([]AirwaySection, 0, len(exAirwaySections))
	for _, section := range exAirwaySections {
		airwaySection := AirwaySection{
			AirwaySectionID: section["airway_section_id"].(string),
			StartAt:         section["start_at"].(string),
			EndAt:           section["end_at"].(string),
		}
		airwaySections = append(airwaySections, airwaySection)
	}
	return airwaySections, nil
}

func (airwayReservationConverter) ToCreateAirwayReservationModel(req *CreateAirwayReservationRequest) (*model.AirwayReservation, error) {
	t, err := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	exAirwaySections, err := convertExAirwaySections(req)
	if err != nil {
		return nil, err
	}
	return &model.AirwayReservation{
		ExAirwaySections: exAirwaySections,
		AcceptedAt:       t,
		AirspaceID:       "a73b0f60-574c-409c-a9ab-3bebeb60dcfa",
		ReservedBy:       value.ModelID(req.OperatorID),
		Status:           "RESERVED",
	}, nil
}

func (airwayReservationConverter) NewAirwayReservationResponse() *AirwayReservationResponse {
	return &AirwayReservationResponse{
		AirwayReservationID: "",
		OperatorID:          "",
		AirwaySections:      []AirwaySection{},
		Status:              value.AirwayReservationStatus(""),
		ReservedAt:          "",
		UpdatedAt:           "",
	}
}
func (airwayReservationConverter) ToAirwayReservationResponse(ar *model.AirwayReservation) (*AirwayReservationResponse, error) {
	exAirwaySections, err := convertExAirwaySectionsFromJSON(ar.ExAirwaySections)
	if err != nil {
		return nil, err
	}
	return &AirwayReservationResponse{
		AirwayReservationID: ar.ID.ToString(),
		OperatorID:          ar.ReservedBy.ToString(),
		AirwaySections:      exAirwaySections,
		Status:              ar.Status,
		ReservedAt:          ar.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           ar.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (airwayReservationConverter) ToAirwayReservationResponses(ars *[]model.AirwayReservation, pageInfo *model.Page) (*AirwayReservationResponses, error) {
	if ars == nil {
		return nil, nil
	}

	responses := AirwayReservationResponses{
		Result: []AirwayReservationResponse{},
	}
	for _, ar := range *ars {
		exAirwaySections, err := convertExAirwaySectionsFromJSON(ar.ExAirwaySections)
		if err != nil {
			return nil, err
		}
		response := AirwayReservationResponse{
			AirwayReservationID: ar.ID.ToString(),
			OperatorID:          ar.ReservedBy.ToString(),
			AirwaySections:      exAirwaySections,
			Status:              ar.Status,
			ReservedAt:          ar.CreatedAt.Format(time.RFC3339),
			UpdatedAt:           ar.UpdatedAt.Format(time.RFC3339),
		}
		responses.Result = append(responses.Result, response)
	}
	responses.Total = pageInfo.Total
	responses.CurrentPage = pageInfo.CurrentPage
	responses.LastPage = pageInfo.LastPage
	responses.PerPage = pageInfo.PerPage
	return &responses, nil
}

func (airwayReservationConverter) ToAirwayReservationMessage(response *AirwayReservationResponse) *AirwayReservationMessage {
	if response == nil {
		return nil
	}

	return &AirwayReservationMessage{
		EventID:             util.GetId(),
		AirwayReservationID: response.AirwayReservationID,
		OperatorID:          response.OperatorID,
		AirwaySections:      response.AirwaySections,
		Status:              response.Status,
		ReservedAt:          response.ReservedAt,
		UpdatedAt:           response.UpdatedAt,
	}
}

func (airwayReservationConverter) ToUpdateAirwayReservationModel(req *UpdateAirwayReservationRequest) (*model.AirwayReservation, error) {
	return &model.AirwayReservation{
		ID:     value.ModelID(req.AirwayReservationID),
		Status: value.AirwayReservationStatus(req.Status),
	}, nil
}
