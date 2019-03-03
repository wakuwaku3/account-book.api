package di

import (
	"log"

	"github.com/tampopos/dijct"
	"github.com/wakuwaku3/account-book.api/src/1-application-business-rules/usecases"
	"github.com/wakuwaku3/account-book.api/src/2-interface-adapters/ctrl"
	"github.com/wakuwaku3/account-book.api/src/3-framework-and-drivers/env"
	"github.com/wakuwaku3/account-book.api/src/3-framework-and-drivers/store"
)

// CreateContainer はDIContainerを生成します
func CreateContainer() (dijct.Container, error) {
	container := dijct.NewContainer()

	// settings
	if err := container.Register(env.NewEnv, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(store.NewProvider, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}

	// repositories
	if err := container.Register(store.CreateUsersRepository); err != nil {
		return nil, err
	}

	// controllers
	if err := container.Register(ctrl.NewHome); err != nil {
		return nil, err
	}

	// initialize
	if err := container.Invoke(initialize); err != nil {
		return nil, err
	}

	return container, nil
}
func initialize(envService usecases.Env, storeProvider store.Provider) {
	err := envService.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
	err = storeProvider.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
}
