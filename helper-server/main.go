package main

import (
	"context"
	"fmt"
	"helper-server/application"
	"helper-server/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/uptrace/uptrace-go/uptrace"
)

var (
	DefaultTracer trace.Tracer
	GlobalSpan    trace.Span
	ctx           context.Context
)

func InitOtel(dsn string) error {
	if DefaultTracer != nil {
		return fmt.Errorf("Otel already initialized")
	}
	ctx = context.Background()
	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(dsn),
	)
	uptrace.SetLogger(log.Default())
	// Create a tracer. Usually, tracer is a global variable.
	DefaultTracer = otel.Tracer("DEFAULT_GLOBAL")
	ctx, GlobalSpan = DefaultTracer.Start(ctx, "main-operation")
	return nil
}

func main() {

	appConfig, err := config.LoadLaunchConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	app, err := application.New(ctx, appConfig)
	if err != nil {
		log.Fatal(err)
	}

	app.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for killcode and stop application
	<-sigCh
	cancelFunc()
}
