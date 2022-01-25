package main

import (
	"context"
	"fmt"
	stlog "log"

	"sazid.github.io/distributed_systems/app/grades"
	"sazid.github.io/distributed_systems/app/log"
	"sazid.github.io/distributed_systems/app/registry"
	"sazid.github.io/distributed_systems/app/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.Registration
	r.ServiceName = registry.GradingService
	r.ServiceURL = serviceAddress
	r.RequiredServices = []registry.ServiceName{
		registry.LogService,
	}
	r.ServiceUpdateUrl = r.ServiceURL + "/services"

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

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %v\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
