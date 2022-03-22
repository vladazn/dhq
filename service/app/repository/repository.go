package repository

import (
	"context"
	"github.com/vladazn/dhq/service/app/domain"
	"github.com/vladazn/dhq/service/app/pkg/mariadb"
)

type Answers interface {
	GetHistory(ctx context.Context, userId int, key string) ([]domain.Action, error)
	GetLastAction(ctx context.Context, userId int, key string) (*domain.Action, error)
	AddAction(ctx context.Context, userId int, action *domain.Action) error
	Migrate(ctx context.Context) error
}

type Repositories struct {
	Answers Answers
}

func InitRepositories(mdb mariadb.MariaDB) *Repositories {
	return &Repositories{Answers: newAnswersRepo(mdb)}
}
