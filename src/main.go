package main

import (
	"fmt"
	"os"

	infweb "github.com/wakuwaku3/account-book.api/src/3-framework-and-drivers/web"
)

func main() {
	web, err := infweb.NewWeb()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	web.Start()
}
