package handler

import (
	"airway-reservation/internal/app/airway_reservation/application/converter"
	"airway-reservation/internal/app/airway_reservation/application/usecase"
	"airway-reservation/internal/pkg/myerror"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

/*
******************************
** AirwayReservation
******************************
 */
type AirwayReservationHandler interface {
	List(c echo.Context) error
	FetchByOperatorID(c echo.Context) error
	Get(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type airwayReservationHandler struct {
	airwayReservationUC usecase.AirwayReservationUsecase
}

func NewAirwayReservationHandler(uc usecase.AirwayReservationUsecase) AirwayReservationHandler {
	return &airwayReservationHandler{
		airwayReservationUC: uc,
	}
}

func (th *airwayReservationHandler) List(c echo.Context) error {
	var page int64 = 1

	if queryPage := c.QueryParam("page"); queryPage != "" {
		parsedPage, err := strconv.Atoi(queryPage)
		if err == nil && parsedPage > 0 {
			page = int64(parsedPage)
		}
	}
	airwayReservations, err := th.airwayReservationUC.ListAirwayReservations(page)
	if err != nil {
		statusCode, errMsg := myerror.ConvertToHTTPError(err)
		return c.JSON(statusCode, map[string]string{"error": errMsg})
	}
	return c.JSON(http.StatusOK, airwayReservations)
}

func (th *airwayReservationHandler) FetchByOperatorID(c echo.Context) error {
	operatorID := c.Param("operatorID")
	var page int64 = 1

	if queryPage := c.QueryParam("page"); queryPage != "" {
		parsedPage, err := strconv.Atoi(queryPage)
		if err == nil && parsedPage > 0 {
			page = int64(parsedPage)
		}
	}

	airwayReservations, err := th.airwayReservationUC.FetchAirwayReservationsByOperatorID(operatorID, page)
	if err != nil {
		statusCode, errMsg := myerror.ConvertToHTTPError(err)
		return c.JSON(statusCode, map[string]string{"error": errMsg})
	}
	return c.JSON(http.StatusOK, airwayReservations)
}

func (th *airwayReservationHandler) Get(c echo.Context) error {
	airwayReservationID := c.Param("airwayReservationID")
	airwayReservation, err := th.airwayReservationUC.GetAirwayReservation(airwayReservationID)
	if err != nil {
		statusCode, errMsg := myerror.ConvertToHTTPError(err)
		return c.JSON(statusCode, map[string]string{"error": errMsg})
	}
	return c.JSON(http.StatusOK, airwayReservation)
}

func (th *airwayReservationHandler) Create(c echo.Context) error {
	airwayReservation := converter.NewCreateAirwayReservationRequest()
	// airwayReservation := model.AirwayReservation{}
	fmt.Println("airwayReservation_http", airwayReservation)
	if err := c.Bind(&airwayReservation); err != nil {
		statusCode, errMsg := myerror.ConvertToHTTPError(err)
		return c.JSON(statusCode, map[string]string{"error": errMsg})
	}

	res, err := th.airwayReservationUC.CreateAirwayReservation(airwayReservation)
	if err != nil {
		// return err
		// return c.JSON(http.StatusBadRequest, err.Error())
		statusCode, errMsg := myerror.ConvertToHTTPError(err)
		return c.JSON(statusCode, map[string]string{"error": errMsg})

	}
	fmt.Println("http.Status", http.StatusOK)
	return c.JSON(http.StatusOK, res)
}

func (th *airwayReservationHandler) Update(c echo.Context) error {
	airwayReservationID := c.Param("airwayReservationID")

	// パスを取得
	path := c.Path()
	var status string

	// パスの末尾でstatusを設定
	switch {
	case strings.HasSuffix(path, "/cancel"):
		status = "CANCELED"
	case strings.HasSuffix(path, "/rescind"):
		status = "RESCINDED"
	default:
		return c.JSON(http.StatusBadRequest, "Invalid operation")
	}

	updatedAirwayReservation, err := th.airwayReservationUC.UpdateAirwayReservation(airwayReservationID, status)
	if err != nil {
		statusCode, errMsg := myerror.ConvertToHTTPError(err)
		return c.JSON(statusCode, map[string]string{"error": errMsg})
	}
	return c.JSON(http.StatusOK, updatedAirwayReservation)
}

func (th *airwayReservationHandler) Delete(c echo.Context) error {
	airwayReservationID := c.Param("airwayReservationID")

	deletedAirwayReservationID, err := th.airwayReservationUC.DeleteAirwayReservation(airwayReservationID)
	if err != nil {
		statusCode, errMsg := myerror.ConvertToHTTPError(err)
		return c.JSON(statusCode, map[string]string{"error": errMsg})
	}
	return c.JSON(http.StatusOK, deletedAirwayReservationID)
}
