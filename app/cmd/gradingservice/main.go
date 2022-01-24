package main

import (
	"context"
	"fmt"
	stlog "log"

	"sazid.github.io/distributed_systems/app/grades"
	"sazid.github.io/distributed_systems/app/registry"
	"sazid.github.io/distributed_systems/app/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.Registration
	r.ServiceName = registry.GradingService
	r.ServiceURL = serviceAddress

	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
		return
	}

	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
