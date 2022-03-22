package repository

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/vladazn/dhq/service/app/domain"
	"github.com/vladazn/dhq/service/app/pkg/mariadb"
	"github.com/vladazn/dhq/service/config"
	"os"
	"testing"
)

const configPath = "../../config/config.test.yml"

func shouldProcess() bool {
	return os.Getenv("MANUAL") != ""
}

func TestAnswersRepo(t *testing.T) {
	if !shouldProcess() {
		t.Skip()
	}

	ctx := context.Background()
	_ = ctx

	var err error
	configs, err := config.New(configPath)
	require.NoError(t, err)

	mdb, err := mariadb.NewMariadbConnection(configs.MariaDb)
	require.NoError(t, err)

	h := newAnswersRepo(mdb)

	err = h.Migrate(ctx)
	require.NoError(t, err)

	err = h.AddAction(ctx, 1, &domain.Action{
		Event: "create",
		Data: &domain.Answer{
			Key:   "aaa",
			Value: "bbb",
		},
	})
	require.NoError(t, err)

	err = h.AddAction(ctx, 1, &domain.Action{
		Event: "create",
		Data: &domain.Answer{
			Key:   "aaa",
			Value: "ccc",
		},
	})
	require.NoError(t, err)

	action, err := h.GetLastAction(ctx, 1, "bbb")
	require.NoError(t, err)
	require.Empty(t, action)

	action, err = h.GetLastAction(ctx, 1, "aaa")
	require.NoError(t, err)
	require.NotEmpty(t, action)
	require.NotEmpty(t, action.Data)
	require.Equal(t, "ccc", action.Data.Value)

	history, err := h.GetHistory(ctx, 1, "aaa")
	require.NoError(t, err)
	require.Equal(t, 2, len(history))

	history2, err := h.GetHistory(ctx, 1, "ccc")
	require.NoError(t, err)
	require.Equal(t, 0, len(history2))

}
