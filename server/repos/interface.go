package repos

import (
	"ecom/server/repos/mock"
	"ecom/server/repos/postgres"

	"gorm.io/gorm"
)

type IRepo interface {
}

func NewPostgresDB(db *gorm.DB) *postgres.PostgresRepo {
	return &postgres.PostgresRepo{DB: db}
}

func NewMockDB() *mock.MockRepo {
	return &mock.MockRepo{}
}
