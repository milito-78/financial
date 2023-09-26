package bootstrap

import (
	"financial/application"
	"financial/config"
	"financial/infrastructure/db"
	"financial/infrastructure/db/psql"
	"github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitDatabase() {
	err := container.Singleton(func() *gorm.DB {
		t := db.Factory(config.Default.(*config.App).Database)
		return t
	})

	if err != nil {
		log.Fatalf("Error during generate singleton : %s", err)
	}
}

func InitDependencies() {
	initRepositories()
	initUseCases()
}

func initRepositories() {
	_ = container.Singleton(func(db *gorm.DB) *psql.UserRepository {
		return psql.NewUserRepository(db)
	})

	_ = container.Singleton(func(repo *psql.UserRepository) db.UserReader {
		return repo
	})

	_ = container.Singleton(func(repo *psql.UserRepository) db.UserWriter {
		return repo
	})
}

func initUseCases() {
	_ = container.Singleton(func(reader db.UserReader, writer db.UserWriter) application.IUserService {
		return application.NewUserService(reader, writer)
	})
}
