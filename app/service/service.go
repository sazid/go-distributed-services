package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"sazid.github.io/distributed_systems/app/registry"
)

// Start starts a new service and registers it to the service registry.
func Start(
	ctx context.Context,
	host string,
	port string,
	reg registry.Registration,
	registerHandlersFunc func(),
) (context.Context, error) {

	registerHandlersFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(
	ctx context.Context,
	serviceName registry.ServiceName,
	host string,
	port string,
) context.Context {

	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("%v started on %s:%s. Press any key to stop.\n", serviceName, host, port)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
