package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"github.com/jasonzhao47/cuddle/internal/web/cache"
	"time"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrRecordNotFound = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao       *dao.UserDAO
	userCache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, userCache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:       dao,
		userCache: userCache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func (repo *UserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	err := repo.dao.UpdateById(ctx, repo.toEntity(user))
	return err
}

func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	// 两种方式，要不就不查数据库，redis gg = 系统业务 gg
	// 可能是真没有
	// 但也有可能是redis挂了，网络链接不好
	// 可以通过定义错误来屏蔽掉第二种情况，让它也去查数据库

	// 一定要查
	cu, err := repo.userCache.Get(ctx, uid)
	if err == nil {
		return cu, nil
	}
	user, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	domainUser := repo.toDomain(user)
	_ = repo.userCache.Set(ctx, domainUser)
	// 如果缓存没更新，会造成数据直接的不一致
	// 保守做法：处理错误
	// 激进做法，数据不一致
	return repo.toDomain(user), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
	}
}

func (repo *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
