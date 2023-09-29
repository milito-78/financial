package bootstrap

import (
	"financial/application"
	"financial/config"
	"financial/infrastructure/cache"
	database "financial/infrastructure/db"
	"github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitDatabase() {
	err := container.Singleton(func() *gorm.DB {
		t := database.Factory(config.Default.(*config.App).Database)
		return t
	})

	if err != nil {
		log.Fatalf("Error during generate singleton : %s", err)
	}
}

func InitCache() {
	err := container.Singleton(func() cache.ICache {
		c := cache.NewCache(config.Default.(*config.App).Cache)
		return c
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
	_ = container.Singleton(func(db *gorm.DB) *database.UserRepository {
		return database.NewUserRepository(db)
	})

	_ = container.Singleton(func(repo *database.UserRepository) database.UserReader {
		return repo
	})

	_ = container.Singleton(func(repo *database.UserRepository) database.UserWriter {
		return repo
	})

	_ = container.Singleton(func(db *gorm.DB) *database.GroupRepository {
		return database.NewGroupRepository(db)
	})

	_ = container.Singleton(func(repo *database.GroupRepository) database.GroupReader {
		return repo
	})

	_ = container.Singleton(func(repo *database.GroupRepository) database.GroupWriter {
		return repo
	})
}

func initUseCases() {
	_ = container.Singleton(func(reader database.UserReader, writer database.UserWriter) application.IUserService {
		return application.NewUserService(reader, writer)
	})

	_ = container.Singleton(func(reader database.GroupReader, writer database.GroupWriter) application.IGroupService {
		return application.NewGroupService(reader, writer)
	})
}
