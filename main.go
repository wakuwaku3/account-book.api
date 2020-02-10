package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wakuwaku3/account-book.api/src/adapter/event"
	"github.com/wakuwaku3/account-book.api/src/adapter/system/di"
	infweb "github.com/wakuwaku3/account-book.api/src/adapter/web"
)

func main() {
	container, err := di.CreateContainer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	var subscriber event.Subscriber
	container.Invoke(func(s event.Subscriber) {
		subscriber = s
	})
	if err := subscriber.Subscribe(container); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	web, err := infweb.NewWeb(container)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	port := "8080"
	if len(os.Args) > 1 {
		if i, err := strconv.Atoi(os.Args[1]); err == nil && i >= 0 && i <= 65535 {
			port = os.Args[1]
		}
	}
	web.Start(port)
}
