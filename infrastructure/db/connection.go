package db

import (
	"financial/config"
	"financial/infrastructure/db/psql"
	"gorm.io/gorm"
)

const (
	postgres = "postgres"
	mysql    = "mysql"
)

func Factory(cfg config.Database) *gorm.DB {
	var connect *gorm.DB
	switch cfg.Driver {
	case mysql:
		break
	case postgres:
		connect = psql.GetConnection(cfg)
		break
	default:
		connect = psql.GetConnection(cfg)
	}

	return connect
}
