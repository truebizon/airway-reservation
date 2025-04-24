package interfaces

import (
	"context"

	"gorm.io/gorm"
)

type DB struct {
	Gorm *gorm.DB
}

type TransactionDBIF interface {
	TransactionScope(context.Context, func(context.Context, *DB) error) (err error)
	GetDB() *DB
}
