package airwayValidator

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

type CreateAirwayReservationInput struct {
	ExAirwaySections []string `json:"exAirwaySections"`
	AcceptedAt       string   `json:"acceptedAt"`
	ReservedBy       string   `json:"reservedBy"`
	AirspaceID       string   `json:"airspaceId"`
	Status           string   `json:"status"`
}

// UpdateAirwayReservationの引数の構造体
type UpdateAirwayReservationInput struct {
	AirwayReservationID string `json:"airwayReservationId"`
	Status              string `json:"status"`
}

func New(modelType interface{}) (*validator.Validate, error) {
	valid := validator.New()

	valid.RegisterStructValidation(validateAirwayStruct, modelType)

	// JSON フィールド名の設定
	useJsonFieldName(valid)

	return valid, nil
}

func useJsonFieldName(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func CustomErrorMessage(err error) error {
	if err == nil {
		return nil
	}

	for _, fieldErr := range err.(validator.ValidationErrors) {
		switch fieldErr.Tag() {
		case "datetime-format":
			return fmt.Errorf("%s の形式が不正です", fieldErr.Field())
		case "datetime-after-now":
			return fmt.Errorf("start_at は現在時刻以降である必要があります")
		case "datetime-order":
			return fmt.Errorf("end_at は start_at より後である必要があります")
		case "required-field":
			return fmt.Errorf("マーカー情報を入力する場合 MarkerName と MarkerArea の両方を指定してください")
		case "required":
			return fmt.Errorf("%s フィールドが存在しません", fieldErr.Field())
		}
	}

	return err
}

func validateAirwayStruct(sl validator.StructLevel) {
	// 日付のバリデーション
	validAirwayDates(sl)
}

func validAirwayDates(sl validator.StructLevel) {
	param := sl.Current().Interface()

	// StartAt と EndAt と AcceptedAtを取得
	startAtStr := reflect.ValueOf(param).FieldByName("StartAt").String()
	endAtStr := reflect.ValueOf(param).FieldByName("EndAt").String()
	acceptedAtStr := reflect.ValueOf(param).FieldByName("AcceptedAt").String()

	// startAt の形式を検証
	startAt, err := time.Parse(time.RFC3339, startAtStr)
	if err != nil {
		sl.ReportError(param, "startAt", "StartAt", "datetime-format", "")
		return
	}

	// endAt の形式を検証
	endAt, err := time.Parse(time.RFC3339, endAtStr)
	if err != nil {
		sl.ReportError(endAt, "endAt", "EndAt", "datetime-format", "")
		return
	}

	// acceptedAt の形式を検証
	acceptedAt, err := time.Parse(time.RFC3339, acceptedAtStr)
	if err != nil {
		sl.ReportError(acceptedAt, "acceptedAt", "AcceptedAt", "datetime-format", "")
		return
	}

	// 現在時刻以降であるかを確認
	if startAt.Before(time.Now()) {
		sl.ReportError(startAt, "startAt", "StartAt", "datetime-after-now", "")
	}

	// startAt が endAt より前であるかを確認
	if endAt.Before(startAt) {
		sl.ReportError(endAt, "endAt", "EndAt", "datetime-order", "")
	}
}
