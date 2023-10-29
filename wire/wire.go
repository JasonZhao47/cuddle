//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jasonzhao47/cuddle/wire/repository"
	"github.com/jasonzhao47/cuddle/wire/repository/dao"
)

func InitUserRepository() *repository.UserRepository {

	wire.Build(repository.NewUserRepository, dao.NewUserDAO, repository.InitDB)
	return &repository.UserRepository{}
}
