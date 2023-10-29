package service

import (
	"context"
	"errors"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// 定义业务处理的异常
var (
	ErrDuplicateEmail        = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

// UserService 用户相关服务
type UserService interface {
	FindById(ctx context.Context, uid int64) (domain.User, error)
	Signup(context.Context, domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
	UpdateNonPII(ctx context.Context, user domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrRecordNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, err
	}
	return user, err
}

func (svc *userService) UpdateNonPII(ctx context.Context, user domain.User) error {
	err := svc.repo.Update(ctx, user)
	return err
}

func (svc *userService) FindById(ctx context.Context, uid int64) (domain.User, error) {
	user, err := svc.repo.FindById(ctx, uid)
	return user, err
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	_, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrRecordNotFound {
		return domain.User{}, err
	}

	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}

	return svc.repo.FindByPhone(ctx, phone)
}
