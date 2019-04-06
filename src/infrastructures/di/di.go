package di

import (
	"log"
	"reflect"

	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/mails"

	"github.com/wakuwaku3/account-book.api/src/application/queries"

	"github.com/tampopos/dijct"
	"github.com/wakuwaku3/account-book.api/src/adapter/ctrls"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/services"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/auth"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/crypt"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/env"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store/repos"
	"github.com/wakuwaku3/account-book.api/src/application/usecases"
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
	if err := container.Register(cmn.NewClock, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	ifs := []reflect.Type{reflect.TypeOf((*application.ClaimsProvider)(nil)).Elem()}
	if err := container.Register(auth.NewAnonymousClaimsProvider(), dijct.RegisterOptions{Interfaces: ifs}); err != nil {
		return nil, err
	}

	// mails
	if err := container.Register(mails.NewHelper, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(mails.NewResetPassword, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(mails.NewUserCreation, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(mails.NewUserExisting, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}

	// controllers
	if err := container.Register(ctrls.NewAccounts); err != nil {
		return nil, err
	}
	if err := container.Register(ctrls.NewTransactions); err != nil {
		return nil, err
	}
	if err := container.Register(ctrls.NewPlans); err != nil {
		return nil, err
	}
	if err := container.Register(ctrls.NewDashboard); err != nil {
		return nil, err
	}
	if err := container.Register(ctrls.NewActual); err != nil {
		return nil, err
	}

	// usecases
	if err := container.Register(usecases.NewAccounts); err != nil {
		return nil, err
	}
	if err := container.Register(usecases.NewTransactions); err != nil {
		return nil, err
	}
	if err := container.Register(usecases.NewPlans); err != nil {
		return nil, err
	}
	if err := container.Register(usecases.NewDashboard); err != nil {
		return nil, err
	}
	if err := container.Register(usecases.NewActual); err != nil {
		return nil, err
	}

	// queries
	if err := container.Register(queries.NewAccounts); err != nil {
		return nil, err
	}
	if err := container.Register(queries.NewTransactions); err != nil {
		return nil, err
	}
	if err := container.Register(queries.NewPlans); err != nil {
		return nil, err
	}
	if err := container.Register(queries.NewDashboard); err != nil {
		return nil, err
	}
	if err := container.Register(queries.NewActual); err != nil {
		return nil, err
	}

	// services
	if err := container.Register(services.NewAccounts); err != nil {
		return nil, err
	}
	if err := container.Register(services.NewTransactions); err != nil {
		return nil, err
	}
	if err := container.Register(services.NewPlans); err != nil {
		return nil, err
	}
	if err := container.Register(services.NewDashboard); err != nil {
		return nil, err
	}
	if err := container.Register(services.NewActual); err != nil {
		return nil, err
	}

	// repos
	if err := container.Register(repos.NewUsers); err != nil {
		return nil, err
	}
	if err := container.Register(repos.NewAccounts); err != nil {
		return nil, err
	}
	if err := container.Register(repos.NewTransactions); err != nil {
		return nil, err
	}
	if err := container.Register(repos.NewPlans); err != nil {
		return nil, err
	}
	if err := container.Register(repos.NewDashboard); err != nil {
		return nil, err
	}

	// initialize
	if err := container.Invoke(initialize); err != nil {
		return nil, err
	}

	return container, nil
}
func initialize(envService application.Env, storeProvider store.Provider) {
	err := envService.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
	err = storeProvider.Initialize()
	if err != nil {
		log.Fatalln(err)
	}
}
