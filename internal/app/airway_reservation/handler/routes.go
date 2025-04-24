package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHealthCheckRouter(e *echo.Echo) *echo.Echo {
	// ヘルスチェック用の API エンドポイント
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "tm service / path ok")
	})

	return e
}

// airwayReservation
func NewAirwayReservationRouter(airwayReservationHH AirwayReservationHandler, e *echo.Echo) *echo.Echo {
	admin := e.Group("/v1/admin/airwayReservations")
	airwayReservation := e.Group("/v1/airwayReservations")
	operator := e.Group("/v1/operator")
	admin.GET("", airwayReservationHH.List)
	operator.GET("/:operatorID/airwayReservations", airwayReservationHH.FetchByOperatorID)
	airwayReservation.GET("/:airwayReservationID", airwayReservationHH.Get)
	airwayReservation.POST("", airwayReservationHH.Create)
	airwayReservation.PUT("/:airwayReservationID/cancel", airwayReservationHH.Update)
	admin.PUT("/:airwayReservationID/rescind", airwayReservationHH.Update)

	return e
}
