//go:build wireinject

package wire

import (
	"github.com/google/wire"
	repository2 "github.com/jasonzhao47/cuddle/playground/wire/repository"
	"github.com/jasonzhao47/cuddle/playground/wire/repository/dao"
)

func InitUserRepository() *repository2.UserRepository {

	wire.Build(repository2.NewUserRepository, dao.NewUserDAO, repository2.InitDB)
	return &repository2.UserRepository{}
}
