package main

import (
	"fmt"
	"os"
	"strconv"

	infweb "github.com/wakuwaku3/account-book.api/src/adapter/web"
)

func main() {
	web, err := infweb.NewWeb()
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
