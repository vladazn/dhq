package app

import (
	"context"
	"fmt"
	grpcserver "github.com/vladazn/dhq/service/app/pkg/gprc/server"
	"github.com/vladazn/dhq/service/app/pkg/mariadb"
	"github.com/vladazn/dhq/service/app/pkg/redis"
	"github.com/vladazn/dhq/service/app/repository"
	"github.com/vladazn/dhq/service/app/service"
	"github.com/vladazn/dhq/service/config"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	ctx := context.Background()

	var err error
	configs, err := config.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%#v", configs)

	r, err := redis.NewRedisConnection(configs.Redis)
	if err != nil {
		fmt.Println(err)
		return
	}

	mdb, err := mariadb.NewMariadbConnection(configs.MariaDb)
	if err != nil {
		fmt.Println(err)
		return
	}

	rp := repository.InitRepositories(mdb)

	err = rp.Answers.Migrate(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	services := service.InitServices(rp, r)

	server := grpcserver.NewGrpcServer(configs.Grpc.Host, configs.Grpc.Port, services)

	errCh := make(chan error)
	go func() {
		err := server.Serve()

		if err != nil {
			errCh <- err
		}
	}()

	fmt.Printf("grpc server started on port: %v\n", configs.Grpc.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-errCh:
		fmt.Println(err)
	case <-quit:
		fmt.Println("quit call")
	}

	fmt.Println("stopping")

	_ = r.Close()
	server.Stop()

}
