package di

import (
	"log"

	"github.com/wakuwaku3/account-book.api/src/usecases/queries"

	"github.com/tampopos/dijct"
	"github.com/wakuwaku3/account-book.api/src/ctrls"
	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/services"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/auth"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/crypt"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/env"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store/repos"
	"github.com/wakuwaku3/account-book.api/src/usecases"
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
	if err := container.Register(crypt.NewCrypt, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(auth.NewJwt, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}

	// controllers
	if err := container.Register(ctrls.NewAccounts); err != nil {
		return nil, err
	}

	// usecases
	if err := container.Register(usecases.NewAccounts); err != nil {
		return nil, err
	}

	// queries
	if err := container.Register(queries.NewAccounts); err != nil {
		return nil, err
	}

	// services
	if err := container.Register(services.NewAccounts); err != nil {
		return nil, err
	}

	// repos
	if err := container.Register(repos.NewUsers); err != nil {
		return nil, err
	}
	if err := container.Register(repos.NewAccounts); err != nil {
		return nil, err
	}

	// initialize
	if err := container.Invoke(initialize); err != nil {
		return nil, err
	}

	return container, nil
}
func initialize(envService domains.Env, storeProvider store.Provider) {
	err := envService.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
	err = storeProvider.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
}
