package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println(getVersion())

	addr := flag.String("l", ":8081", "The listen address of service.")
	flag.StringVar(&host, "d", "mikanani.me", "The host of RSS.")
	flag.Parse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Add(echo.GET, "/", handler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(*addr); err != nil && !errors.Is(http.ErrServerClosed, err) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	e.Logger.Info("starting shutdown the server..")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("shutdown the server success")
}
