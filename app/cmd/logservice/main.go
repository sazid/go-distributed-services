package main

import (
	"context"
	"fmt"
	stlog "log"

	"sazid.github.io/distributed_systems/app/log"
	"sazid.github.io/distributed_systems/app/registry"
	"sazid.github.io/distributed_systems/app/service"
)

func main() {
	log.Run("./app.log")

	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	reg := registry.Registration{
		ServiceName: registry.LogService,
		ServiceURL:  serviceAddress,
	}

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		reg,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()

	fmt.Println("Shutting down log service")
}
