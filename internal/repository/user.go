package repository

import (
	"context"
	"database/sql"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate  = dao.ErrUserDuplicate
	ErrRecordNotFound = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	// Update 更新数据，只有非 0 值才会更新
	Update(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	// FindByWechat 暂时可以认为按照 openId来查询
	// 将来可能需要按照 unionId 来查询
	//FindByWechat(ctx context.Context, openId string) (domain.User, error)
}

type CacheUserRepository struct {
	dao       dao.UserDAO
	userCache cache.UserCache
}

func NewCacheUserRepository(dao dao.UserDAO, userCache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:       dao,
		userCache: userCache,
	}
}

func (repo *CacheUserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Password: user.Password,
	})
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func (repo *CacheUserRepository) Update(ctx context.Context, user domain.User) error {
	err := repo.dao.UpdateById(ctx, repo.toEntity(user))
	return err
}

func (repo *CacheUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	// 两种方式，要不就不查数据库，redis gg = 系统业务 gg
	// 可能是真没有
	// 但也有可能是redis挂了，网络链接不好
	// 可以通过定义错误来屏蔽掉第二种情况，让它也去查数据库
	// 什么都不做也是去查数据库了
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
	// 保守做法：处理错误，如果返回，会导致业务没存上
	// 激进做法，数据不一致，但业务保住了
	return repo.toDomain(user), nil
}

func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func (repo *CacheUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
	}
}

func (repo *CacheUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Email,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
