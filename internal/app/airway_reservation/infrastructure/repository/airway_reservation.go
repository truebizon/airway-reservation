package repository

import (
	"airway-reservation/internal/app/airway_reservation/domain/model"
	"airway-reservation/internal/app/airway_reservation/domain/repositoryIF"
	"airway-reservation/internal/pkg/myerror"
	"airway-reservation/internal/pkg/value"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type airwayReservationRepository struct {
	db *gorm.DB
}

func NewAirwayReservationRepository(db *gorm.DB) repositoryIF.AirwayReservationRepositoryIF {
	return &airwayReservationRepository{db}
}
func (airwayReservationRepo *airwayReservationRepository) FetchAll(airwayReservations *[]model.AirwayReservation) error {
	if err := airwayReservationRepo.db.Find(&airwayReservations).Error; err != nil {
		return err
	}
	return nil
}

func (airwayReservationRepo *airwayReservationRepository) FetchALLWithPagination(
	page int64,
) (*[]model.AirwayReservation, *model.Page, error) {
	const perPage = 20
	offset := (page - 1) * perPage

	var total int64
	var airwayReservations *[]model.AirwayReservation

	// 総件数の取得
	if err := airwayReservationRepo.db.Model(&model.AirwayReservation{}).
		Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// lastPage の計算 (件数が perPage で割り切れる場合を考慮)
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))

	// ページングデータの取得
	if err := airwayReservationRepo.db.
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(perPage).
		Find(&airwayReservations).Error; err != nil {
		return nil, nil, err
	}

	// Page 構造体を作成
	pageInfo := &model.Page{
		Total:       total,
		CurrentPage: page,
		LastPage:    int64(lastPage),
		PerPage:     perPage,
	}

	return airwayReservations, pageInfo, nil
}

func (airwayReservationRepo *airwayReservationRepository) FetchByOperatorWithPagination(
	operatorID string, page int64,
) (*[]model.AirwayReservation, *model.Page, error) {
	const perPage = 20
	offset := (page - 1) * perPage

	var total int64
	var airwayReservations *[]model.AirwayReservation

	// 総件数の取得
	if err := airwayReservationRepo.db.Model(&model.AirwayReservation{}).
		Where("reserved_by = ?", operatorID).
		Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// lastPage の計算 (件数が perPage で割り切れる場合を考慮)
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))

	// ページングデータの取得
	if err := airwayReservationRepo.db.
		Where("reserved_by = ?", operatorID).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(perPage).
		Find(&airwayReservations).Error; err != nil {
		return nil, nil, err
	}

	// Page 構造体を作成
	pageInfo := &model.Page{
		Total:       total,
		CurrentPage: page,
		LastPage:    int64(lastPage),
		PerPage:     perPage,
	}

	return airwayReservations, pageInfo, nil
}

func (airwayReservationRepo *airwayReservationRepository) FindByID(airwayReservationID string, airwayReservation *model.AirwayReservation) error {
	if err := airwayReservationRepo.db.First(airwayReservation, "id = ?", airwayReservationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.Wrap(myerror.NotFound, fmt.Errorf(""), "record not found")
		}
		return err
	}
	return nil
}

func (airwayReservationRepo *airwayReservationRepository) InsertOne(airwayReservation *model.AirwayReservation) (*model.AirwayReservation, error) {
	// データ挿入前に詳細情報をログ出力
	fmt.Printf("Inserting AirwayReservation: %+v\n", airwayReservation)
	fmt.Printf("ExAirwaySections (Raw): %+v\n", airwayReservation.ExAirwaySections)
	fmt.Printf("ExAirwaySections (String): %s\n", string(airwayReservation.ExAirwaySections))

	// トランザクションを使用して、より詳細なエラー情報を取得
	tx := airwayReservationRepo.db.Begin()
	if err := tx.Create(airwayReservation).Error; err != nil {
		// エラーロールバック
		tx.Rollback()

		// エラー詳細を詳細に出力
		fmt.Printf("Insert Error: %v\n", err)

		// SQLエラーの場合、追加情報を出力
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Printf("PostgreSQL Error Code: %s\n", pqErr.Code)
			fmt.Printf("PostgreSQL Error Detail: %s\n", pqErr.Detail)
			fmt.Printf("PostgreSQL Error Hint: %s\n", pqErr.Hint)
		}

		return airwayReservation, err
	}

	// トランザクションをコミット
	tx.Commit()

	return airwayReservation, nil
}
func (airwayReservationRepo *airwayReservationRepository) UpdateOne(airwayReservation *model.AirwayReservation) (*model.AirwayReservation, error) {
	updateData := make(map[string]interface{})
	if airwayReservation.Status != "" {
		updateData["status"] = airwayReservation.Status
	}

	var updatedAirwayReservation model.AirwayReservation
	result := airwayReservationRepo.db.Model(&model.AirwayReservation{}).
		Select("*").
		Clauses(clause.Returning{}).
		Where("id = ?", airwayReservation.ID).
		Updates(updateData).
		Scan(&updatedAirwayReservation)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, myerror.Wrap(myerror.NotFound, fmt.Errorf(""), "record not found")

	}
	return &updatedAirwayReservation, nil
}
func (airwayReservationRepo *airwayReservationRepository) DeleteOne(airwayReservationID value.ModelID) (value.ModelID, error) {
	result := airwayReservationRepo.db.Where("id=?", airwayReservationID).Delete(&model.AirwayReservation{})
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected < 1 {
		return "", fmt.Errorf("object does not exist")
	}
	return airwayReservationID, nil
}
