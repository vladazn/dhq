package service

import (
	"context"
	"github.com/vladazn/dhq/service/app/domain"
	"github.com/vladazn/dhq/service/app/pkg/redis"
	"github.com/vladazn/dhq/service/app/repository"
)

type Storage interface {
	Create(ctx context.Context, userId int, answer *domain.Answer) (bool,
		error)
	Delete(ctx context.Context, userId int, key string) (bool, error)
	Update(ctx context.Context, userId int, answer *domain.Answer) (bool,
		error)
	Get(ctx context.Context, userId int, key string) (*domain.Answer, error)
	History(ctx context.Context, userId int, key string) ([]domain.Action,
		error)
}

type Services struct {
	Storage Storage
}

func InitServices(repo *repository.Repositories, redis redis.Redis) *Services {

	return &Services{
		Storage: newStorageService(repo, redis),
	}
}
