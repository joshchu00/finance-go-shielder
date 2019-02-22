package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/joshchu00/finance-go-common/config"
	"github.com/joshchu00/finance-go-common/logger"
	pb "github.com/joshchu00/finance-protobuf/porter"
	"google.golang.org/grpc"
)

func init() {

	// config
	config.Init()

	// logger
	logger.Init(config.LogDirectory(), "shielder")

	// log config
	logger.Info(fmt.Sprintf("%s: %s", "EnvironmentName", config.EnvironmentName()))
	logger.Info(fmt.Sprintf("%s: %s", "ShielderPort", config.ShielderPort()))
	logger.Info(fmt.Sprintf("%s: %s", "ShielderCORSMethods", config.ShielderCORSMethods()))
	logger.Info(fmt.Sprintf("%s: %s", "ShielderCORSOrigins", config.ShielderCORSOrigins()))
	logger.Info(fmt.Sprintf("%s: %s", "PorterV1Host", config.PorterV1Host()))
	logger.Info(fmt.Sprintf("%s: %s", "PorterV1Port", config.PorterV1Port()))
}

var environmentName string

func process() {

	if environmentName == config.EnvironmentNameProd {
		defer func() {
			if err := recover(); err != nil {
				logger.Panic(fmt.Sprintf("recover %v", err))
			}
		}()
	}

	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()

	err = pb.RegisterPorterV1HandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("%s:%s", config.PorterV1Host(), config.PorterV1Port()),
		[]grpc.DialOption{
			grpc.WithInsecure(),
		},
	)
	if err != nil {
		logger.Panic(fmt.Sprintf("pb.RegisterPorterV1HandlerFromEndpoint %v", err))
	}

	http.ListenAndServe(
		fmt.Sprintf(":%s", config.ShielderPort()),
		handlers.CORS(
			handlers.AllowedMethods(config.ShielderCORSMethods()),
			handlers.AllowedOrigins(config.ShielderCORSOrigins()),
		)(mux),
	)
}

func main() {

	logger.Info("Starting shielder...")

	// environment name
	switch environmentName = config.EnvironmentName(); environmentName {
	case config.EnvironmentNameDev, config.EnvironmentNameTest, config.EnvironmentNameStg, config.EnvironmentNameProd:
	default:
		logger.Panic("Unknown environment name")
	}

	for {

		process()

		time.Sleep(3 * time.Second)

		if environmentName != config.EnvironmentNameProd {
			break
		}
	}
}
