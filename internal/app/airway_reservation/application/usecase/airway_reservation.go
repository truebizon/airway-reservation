package usecase

import (
	"airway-reservation/internal/app/airway_reservation/application/converter"
	"airway-reservation/internal/app/airway_reservation/application/validator"
	"airway-reservation/internal/app/airway_reservation/domain/model"
	"airway-reservation/internal/app/airway_reservation/domain/repositoryIF"
	pkgIF "airway-reservation/internal/pkg/database/interfaces"
	"airway-reservation/internal/pkg/mqtt"
	"airway-reservation/internal/pkg/value"
	"context"
	"encoding/json"
	"fmt"
)

type AirwayReservationUsecase struct {
	ctx                     context.Context
	transactionDBIF         pkgIF.TransactionDBIF
	airwayReservationRepoIF repositoryIF.AirwayReservationRepositoryIF
}

func NewAirwayReservationUsecase(
	ctx context.Context,
	transactionDBIF pkgIF.TransactionDBIF,
	repoIF repositoryIF.AirwayReservationRepositoryIF,
) *AirwayReservationUsecase {
	return &AirwayReservationUsecase{
		ctx:                     ctx,
		transactionDBIF:         transactionDBIF,
		airwayReservationRepoIF: repoIF,
	}
}

func (uc *AirwayReservationUsecase) ListAirwayReservations(page int64) (*converter.AirwayReservationResponses, error) {
	fmt.Println("ListAirwayReservations")
	domainConv := converter.NewAirwayReservationDomain()
	airwayReservations, pageInfo, err := uc.airwayReservationRepoIF.FetchALLWithPagination(page)
	if err != nil {
		return nil, err
	}
	return domainConv.ToAirwayReservationResponses(airwayReservations, pageInfo)
}

func (uc *AirwayReservationUsecase) FetchAirwayReservationsByOperatorID(operatorID string, page int64) (*converter.AirwayReservationResponses, error) {
	fmt.Println("FetchAirwayReservationsByOperatorID")
	domainConv := converter.NewAirwayReservationDomain()
	valid := validator.NewAirwayReservation()
	if err := valid.FetchAirwayReservationsByOperatorID(operatorID); err != nil {
		return nil, err
	}
	airwayReservations, pageInfo, err := uc.airwayReservationRepoIF.FetchByOperatorWithPagination(operatorID, page)
	if err != nil {
		return nil, err
	}

	return domainConv.ToAirwayReservationResponses(airwayReservations, pageInfo)
}

func (uc *AirwayReservationUsecase) GetAirwayReservation(airwayReservationID string) (res *converter.AirwayReservationResponse, err error) {
	fmt.Println("GetAirwayReservation")
	domainConv := converter.NewAirwayReservationDomain()
	valid := validator.NewAirwayReservation()
	if err := valid.GetRequest(airwayReservationID); err != nil {
		return nil, err
	}
	airwayReservation := &model.AirwayReservation{}
	if err := uc.airwayReservationRepoIF.FindByID(airwayReservationID, airwayReservation); err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	return domainConv.ToAirwayReservationResponse(airwayReservation)
}

func (uc *AirwayReservationUsecase) CreateAirwayReservation(req *converter.CreateAirwayReservationRequest) (res *converter.AirwayReservationResponse, err error) {
	fmt.Println("CreateAirwayReservation")
	domainConv := converter.NewAirwayReservationDomain()
	valid := validator.NewAirwayReservation()
	if err := valid.RegisterRequest(req); err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	airwayReservation, err := domainConv.ToCreateAirwayReservationModel(req)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	airwayReservationDomain, err := uc.airwayReservationRepoIF.InsertOne(airwayReservation)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	res, err = domainConv.ToAirwayReservationResponse(airwayReservationDomain)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	msg := domainConv.ToAirwayReservationMessage(res)
	mqttClient, err := mqtt.NewMQTTClient()
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	topic := fmt.Sprintf("airway/operator/%s/airwayReservation/%s", msg.OperatorID, msg.AirwayReservationID)
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	mqttClient.Publish(topic, 1, false, messageBytes)
	return res, nil
}

func (uc *AirwayReservationUsecase) UpdateAirwayReservation(airwayReservationID, status string) (res *converter.AirwayReservationResponse, err error) {
	fmt.Println("UpdateAirwayReservation")
	domainConv := converter.NewAirwayReservationDomain()

	req := domainConv.ToUpdateAirwayReservationRequest(airwayReservationID, status)
	valid := validator.NewAirwayReservation()
	if err := valid.UpdateRequest(req); err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	airwayReservation, err := domainConv.ToUpdateAirwayReservationModel(req)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	airwayReservationDomain, err := uc.airwayReservationRepoIF.UpdateOne(airwayReservation)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	res, err = domainConv.ToAirwayReservationResponse(airwayReservationDomain)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	msg := domainConv.ToAirwayReservationMessage(res)
	mqttClient, err := mqtt.NewMQTTClient()
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	topic := fmt.Sprintf("airway/operator/%s/airwayReservation/%s", msg.OperatorID, msg.AirwayReservationID)
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return domainConv.NewAirwayReservationResponse(), err
	}
	mqttClient.Publish(topic, 1, false, messageBytes)
	return res, nil
}

func (uc *AirwayReservationUsecase) DeleteAirwayReservation(airwayReservationID string) (string, error) {
	fmt.Println("DeleteAirwayReservation")
	valid := validator.NewAirwayReservation()
	if err := valid.DeleteRequest(airwayReservationID); err != nil {
		return airwayReservationID, err
	}
	deleteID := value.ModelID(airwayReservationID)
	deleteAirwayReservationID, err := uc.airwayReservationRepoIF.DeleteOne(deleteID)
	if err != nil {
		return airwayReservationID, err
	}
	return deleteAirwayReservationID.ToString(), nil
}
