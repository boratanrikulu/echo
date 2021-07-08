package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/boratanrikulu/echo"
	"github.com/boratanrikulu/echo/config"
)

func main() {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	opts := config.Configure(fs, os.Args[1:])

	s := echo.NewServer().
		Address(opts.Addr).
		Banner(!opts.NoBanner).
		Verbose(opts.Verbose)

	ctx := context.Background()
	s.Run(ctx)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	fmt.Println("Server is gracefully shutdown.")
}
