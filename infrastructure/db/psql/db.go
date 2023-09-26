package psql

import (
	"financial/config"
	"financial/infrastructure/db/psql/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetConnection(cfg config.Database) *gorm.DB {
	if db != nil {
		return db
	}
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=" + cfg.Host + " user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.Name + " port=" + cfg.Port + " sslmode=disable TimeZone=UTC",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db = conn

	db.AutoMigrate(models.UserEntity{})

	return db
}
