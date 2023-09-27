package db

import (
	"financial/config"
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
		connect = GetConnection(cfg)
		break
	default:
		connect = GetConnection(cfg)
	}

	return connect
}
