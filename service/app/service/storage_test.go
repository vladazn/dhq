package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/vladazn/dhq/service/app/domain"
	"github.com/vladazn/dhq/service/app/pkg/mariadb"
	"github.com/vladazn/dhq/service/app/pkg/redis"
	"github.com/vladazn/dhq/service/app/repository"
	"github.com/vladazn/dhq/service/config"
	"os"
	"testing"
)

const configPath = "../../config/config.test.yml"

func shouldProcess() bool {
	return os.Getenv("MANUAL") != ""
}

func TestStorageService(t *testing.T) {
	if !shouldProcess() {
		t.Skip()
	}

	ctx := context.Background()
	_ = ctx

	var err error
	configs, err := config.New(configPath)
	require.NoError(t, err)

	r, err := redis.NewRedisConnection(configs.Redis)
	require.NoError(t, err)

	mdb, err := mariadb.NewMariadbConnection(configs.MariaDb)
	require.NoError(t, err)

	repo := repository.InitRepositories(mdb)

	err = repo.Answers.Migrate(ctx)
	require.NoError(t, err)

	s := InitServices(repo, r)

	answer, err := s.Storage.Get(ctx, 1, "aaa")
	require.Error(t, err)
	require.Empty(t, answer)

	success, err := s.Storage.Create(ctx, 1, &domain.Answer{
		Key:   "aaa",
		Value: "bbb",
	})
	require.NoError(t, err)
	require.True(t, success)

	answer, err = s.Storage.Get(ctx, 1, "aaa")
	require.NoError(t, err)
	require.NotEmpty(t, answer)
	require.Equal(t, "bbb", answer.Value)

	success, err = s.Storage.Create(ctx, 1, &domain.Answer{
		Key:   "aaa",
		Value: "ccc",
	})
	require.Error(t, err)
	require.False(t, success)

	success, err = s.Storage.Update(ctx, 1, &domain.Answer{
		Key:   "aaa",
		Value: "ccc",
	})
	require.NoError(t, err)
	require.True(t, success)

	answer, err = s.Storage.Get(ctx, 1, "aaa")
	require.NoError(t, err)
	require.NotEmpty(t, answer)
	require.Equal(t, "ccc", answer.Value)

	answer, err = s.Storage.Get(ctx, 1, "aaa")
	require.NoError(t, err)
	require.NotEmpty(t, answer)
	require.Equal(t, "ccc", answer.Value)

	success, err = s.Storage.Delete(ctx, 1, "aaa")
	require.NoError(t, err)
	require.True(t, success)

	answer, err = s.Storage.Get(ctx, 1, "aaa")
	require.Error(t, err)
	require.Empty(t, answer)

	history, err := s.Storage.History(ctx, 1, "aaa")
	require.NoError(t, err)
	require.Equal(t, 3, len(history))

}
