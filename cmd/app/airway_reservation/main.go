package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"airway-reservation/internal/app/airway_reservation/application/usecase"
	"airway-reservation/internal/app/airway_reservation/handler"
	"airway-reservation/internal/app/airway_reservation/infrastructure/repository"
	"airway-reservation/internal/pkg/config"
	database "airway-reservation/internal/pkg/database/postgres"
	"airway-reservation/internal/pkg/logger"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	handler.NewHealthCheckRouter(e)
	e.GET("/airway_reservation", func(c echo.Context) error {
		return c.String(http.StatusOK, "/airway_reservation path ok")
	})
	db := database.NewDB()

	logger.InitLogger()
	config.Load()
	logger.SetLogLevelByConfig()
	ctx := context.Background()
	transactionDBIF := database.NewTransactionDBIF(db)

	// 環境変数で実施可否を制御する。内部で使用するときは別途migrateコマンドを実施するのでfalseとする
	if os.Getenv("USE_EXECUTE_SQL") == "true" {
		if !shouldSkipSQLExecution(db) {
			executeSQL(db)
		}
	}

	airwayReservationRepository := repository.NewAirwayReservationRepository(db)
	airwayReservationUsecase := usecase.NewAirwayReservationUsecase(ctx, transactionDBIF, airwayReservationRepository)
	airwayReservationHandler := handler.NewAirwayReservationHandler(*airwayReservationUsecase)
	handler.NewAirwayReservationRouter(airwayReservationHandler, e)
	e.Logger.Fatal(e.Start(":8088"))
}

// airway_reservationのスキーマとairway_reservations、airspacesテーブルの存在確認
func shouldSkipSQLExecution(db *gorm.DB) bool {
	var count int64

	// airway_reservationスキーマの存在確認
	err := db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.schemata 
		WHERE schema_name = 'airway_reservation';
	`).Scan(&count).Error
	if err != nil {
		log.Fatalf("スキーマ確認エラー: %v", err)
	}
	if count == 0 {
		// スキーマが存在しない場合
		return false
	}

	// airspacesテーブルの存在確認
	err = db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'airway_reservation' 
		AND table_name = 'airspaces';
	`).Scan(&count).Error
	if err != nil {
		log.Fatalf("airspacesテーブル確認エラー: %v", err)
	}
	if count == 0 {
		// テーブルが存在しない場合
		return false
	}

	// airway_reservationsテーブルの存在確認
	err = db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'airway_reservation' 
		AND table_name = 'airway_reservations';
	`).Scan(&count).Error
	if err != nil {
		log.Fatalf("airway_reservationsテーブル確認エラー: %v", err)
	}
	if count == 0 {
		// テーブルが存在しない場合
		return false
	}

	// すべての条件が満たされた場合
	return true
}

// DBの初期化処理
func executeSQL(db *gorm.DB) {
	sqlStatements := `
CREATE SCHEMA "airway_reservation";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
DROP TYPE IF EXISTS airspace_type CASCADE;
CREATE TYPE airspace_type AS ENUM ('AIRWAY', 'NFZ', 'UVR');
DROP TYPE IF EXISTS airspace_class CASCADE;
CREATE TYPE airspace_class AS ENUM ('CTA', 'Xa', 'Xu', 'Xx', 'Y', 'Z');
DROP TYPE IF EXISTS op_shape CASCADE;
CREATE TYPE op_shape AS ENUM ('POLYGON', 'POINT', 'LINESTRING');
DROP TYPE IF EXISTS airway_reservation_status CASCADE;
CREATE TYPE airway_reservation_status AS ENUM ('RESERVED', 'CANCELED', 'RESCINDED');
CREATE TABLE airway_reservation.airspaces (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "organization_id" uuid NOT NULL,
    "project_id" uuid NOT NULL,
    "type" airspace_type NOT NULL,
    "class" airspace_class NOT NULL,
    "volume" geometry NOT NULL,
    "shape" op_shape NOT NULL,
    "max_height" decimal(8, 2),
    "min_height" decimal(8, 2),
    "radius" decimal(8, 2),
    "file" varchar(255),
    "properties" jsonb,
    "valid_from" timestamp(0) with time zone,
    "valid_to" timestamp(0) with time zone,
    "created_at" timestamp(0) with time zone NOT NULL DEFAULT now(),
    "updated_at" timestamp(0) with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);
CREATE TABLE airway_reservation.airway_reservations (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "ex_airway_sections" jsonb NOT NULL,
    "accepted_at" timestamp(0) with time zone,
    "reserved_by" uuid,
    "ex_reserved_by" uuid,
    "ex_airway_id" uuid,
    "airspace_id" uuid NOT NULL,
    "plan_id" uuid,
    "operation_id" uuid,
    "status" airway_reservation_status NOT NULL,
    "created_at" timestamp(0) with time zone NOT NULL DEFAULT now(),
    "updated_at" timestamp(0) with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY ("id"),
    CONSTRAINT fk_airway_reservations_airspace FOREIGN KEY(airspace_id) REFERENCES airway_reservation.airspaces(id)
);
INSERT INTO airway_reservation.airspaces(
    id,
    organization_id,
    project_id,
    "type",
    "class",
    volume,
    shape,
    max_height,
    min_height,
    radius,
    file,
    properties,
    valid_from,
    valid_to,
    created_at,
    updated_at
)
VALUES(
    'a73b0f60-574c-409c-a9ab-3bebeb60dcfa',
    'a73b0f60-574c-409c-a9ab-3bebeb60dcfa',
    '123e4567-e89b-12d3-a456-426614174000',
    'AIRWAY',
    'CTA',
    'POLYGON ((35.8176 139.6915, 35.8099 139.6915, 35.8099 139.7063, 35.8176 139.7063, 35.8176 139.6915))',
    'POLYGON',
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    '2024-11-22 00:21:51.000',
    '2025-11-22 00:21:52.000',
    '2024-11-21 15:21:57.000',
    '2024-11-21 15:21:57.000'
);
INSERT INTO airway_reservation.airway_reservations(
    id,
    ex_airway_sections,
    accepted_at,
    reserved_by,
    ex_reserved_by,
    ex_airway_id,
    airspace_id,
    plan_id,
    operation_id,
    status,
    created_at,
    updated_at
)
VALUES(
    '5a50b6b3-f780-40d4-8c9b-7d6d369640ad',
    '[{"end_at": "2025-02-02T00:00:00Z", "start_at": "2025-02-01T23:59:59Z", "airway_section_id": "123e4567-e89b-12d3-a456-426614174000"}, {"end_at": "2025-02-02T00:02:00Z", "start_at": "2025-02-02T00:01:00Z", "airway_section_id": "123e4567-e89b-12d3-a456-426614174001"}]',
    '2025-01-21 19:09:15.000',
    '60c895e5-321a-fe8a-af39-f005f3206efb',
    NULL,
    NULL,
    'a73b0f60-574c-409c-a9ab-3bebeb60dcfa',
    NULL,
    NULL,
    'RESCINDED',
    '2025-01-21 19:09:15.000',
    '2025-01-22 12:38:26.000'
);
`
	if err := db.Exec(sqlStatements).Error; err != nil {
		fmt.Printf("SQL実行中にエラーが発生しました: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("SQLスクリプトの実行が完了しました")
}
