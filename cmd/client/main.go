package main

import (
	"fmt"
	"os"

	"github.com/aifedorov/gophkeeper/internal/client/app"
)

func main() {
	a := app.NewApp()
	if err := a.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
		os.Exit(1)
	}
}
