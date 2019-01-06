package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	logger.Info(fmt.Sprintf("%s: %s", "Environment", config.Environment()))
	logger.Info(fmt.Sprintf("%s: %s", "ShielderAddress", config.ShielderAddress()))
	logger.Info(fmt.Sprintf("%s: %s", "PorterV1Address", config.PorterV1Address()))
}

var environment string

func process() {

	if environment == config.EnvironmentProd {
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
		config.PorterV1Address(),
		[]grpc.DialOption{
			grpc.WithInsecure(),
		},
	)
	if err != nil {
		logger.Panic(fmt.Sprintf("pb.RegisterPorterV1HandlerFromEndpoint %v", err))
	}

	http.ListenAndServe(config.ShielderAddress(), mux)
}

func main() {

	logger.Info("Starting shielder...")

	// environment
	switch environment = config.Environment(); environment {
	case config.EnvironmentDev, config.EnvironmentTest, config.EnvironmentStg, config.EnvironmentProd:
	default:
		logger.Panic("Unknown environment")
	}

	for {

		process()

		time.Sleep(3 * time.Second)

		if environment != config.EnvironmentProd {
			break
		}
	}
}
