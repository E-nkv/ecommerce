package postgres

import (
	"gorm.io/gorm"
)

type PostgresRepo struct {
	DB *gorm.DB
}
