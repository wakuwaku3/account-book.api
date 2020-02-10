package di

import (
	"log"
	"reflect"

	"github.com/wakuwaku3/account-book.api/src/adapter/crypt"
	"github.com/wakuwaku3/account-book.api/src/adapter/event"
	handler "github.com/wakuwaku3/account-book.api/src/adapter/event/handlers"
	"github.com/wakuwaku3/account-book.api/src/adapter/mails/sendgrid"
	"github.com/wakuwaku3/account-book.api/src/adapter/system"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"

	"github.com/wakuwaku3/account-book.api/src/application/queries"

	"github.com/tampopos/dijct"
	"github.com/wakuwaku3/account-book.api/src/adapter/auth"
	"github.com/wakuwaku3/account-book.api/src/adapter/env"
	"github.com/wakuwaku3/account-book.api/src/adapter/mails"
	"github.com/wakuwaku3/account-book.api/src/adapter/store"
	"github.com/wakuwaku3/account-book.api/src/adapter/store/repos"
	"github.com/wakuwaku3/account-book.api/src/adapter/web/ctrls"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/services"
	"github.com/wakuwaku3/account-book.api/src/application/usecases"
	accountbook "github.com/wakuwaku3/account-book.api/src/enterprise/domains/accountBook"
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
	if err := container.Register(core.NewClock, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	ifs := []reflect.Type{reflect.TypeOf((*core.ClaimsProvider)(nil)).Elem()}
	if err := container.Register(auth.NewAnonymousClaimsProvider(), dijct.RegisterOptions{Interfaces: ifs}); err != nil {
		return nil, err
	}
	if err := container.Register(sendgrid.NewHelper, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(event.NewProvider, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(event.NewPublisher, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(system.NewGuidFactory, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(event.NewRouter, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}
	if err := container.Register(event.NewSubscriber, dijct.RegisterOptions{LifetimeScope: dijct.ContainerManaged}); err != nil {
		return nil, err
	}

	// mails
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
	if err := container.Register(ctrls.NewNotificationRules); err != nil {
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
	if err := container.Register(usecases.NewNotificationRules); err != nil {
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
	if err := container.Register(queries.NewNotificationRules); err != nil {
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
	if err := container.Register(repos.NewNotificationRules); err != nil {
		return nil, err
	}

	//events
	if err := container.Register(accountbook.NewAssetsChangedEvent); err != nil {
		return nil, err
	}

	// handler
	if err := container.Register(handler.NewAlert); err != nil {
		return nil, err
	}

	// initialize
	if err := container.Invoke(initialize); err != nil {
		return nil, err
	}

	return container, nil
}
func initialize(
	envService application.Env,
	storeProvider store.Provider,
	eventProvider event.Provider,
) {
	if err := envService.Initialize(); err != nil {
		log.Fatalln(err)
	}
	if err := storeProvider.Initialize(); err != nil {
		log.Fatalln(err)
	}
	if err := eventProvider.Initialize(); err != nil {
		log.Fatalln(err)
	}
}
