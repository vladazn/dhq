package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/vladazn/dhq/api/config"
	"github.com/vladazn/dhq/proto/gen/go/proto/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func Run(configPath string) {

	ctx := context.Background()
	_ = ctx

	var err error
	configs, err := config.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	customMatcher := func(key string) (string, bool) {
		switch key {
		case "Auth":
			return "auth", true
		default:
			return runtime.DefaultHeaderMatcher(key)
		}
	}

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customMatcher),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = storage.RegisterStorageHandlerFromEndpoint(
		ctx,
		mux,
		configs.Service.Host,
		opts,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	handler := cors.AllowAll().Handler(mux)

	fmt.Printf("serving api at :%v\n", configs.Api.Port)

	err = http.ListenAndServe(
		fmt.Sprintf(":%v", configs.Api.Port),
		handler,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

}
