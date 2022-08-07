package main

import (
	"fmt"
	"jadwalin-auth-events-integrator/internal/controller"
	"jadwalin-auth-events-integrator/internal/di"
	"jadwalin-auth-events-integrator/internal/shared"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var container = di.Container

	err := container.Invoke(func(deps shared.Dependencies, ch controller.Holder) error {
		var sig = make(chan os.Signal, 1)

		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		ch.RegisterRoutes()

		go func() {
			deps.Logger.Infof("Starting server on port %d", deps.Config.HTTPServerPort)
			if err := deps.Echo.Start(fmt.Sprintf(":%d", deps.Config.HTTPServerPort)); err != nil {
				deps.Logger.Errorf("Failed to start server: %s", err)
				sig <- syscall.SIGTERM
			}
		}()

		<-sig
		deps.Logger.Infof("Shutting down server...")
		deps.Close()

		return nil
	})

	if err != nil {
		panic(err)
	}

}
