package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

type Event struct {
	Timestamp time.Time              `json:"Timestamp"`
	Data      map[string]interface{} `json:"Data"`
}

func main() {
	dsn := os.Getenv("UPTRACE_DSN")
	if len(dsn) == 0 {
		log.Fatalf("UPTRACE_DSN is empty\n")
	}
	secret := os.Getenv("ANALYTICS_SECRET")
	if len(secret) == 0 {
		log.Fatalf("ANALYTICS_SECRET is empty\n")
	}
	err := InitOtel(dsn)
	if err != nil {
		log.Fatalf("Failed to init otel: %s\n", err.Error())
	}

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != secret {
			log.Printf("Invalid X-Api-Key: %s\n", r.Header.Get("X-Api-Key"))
			return
		}
		buff, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read body: %s\n", err.Error())
			return
		}
		buff = bytes.TrimPrefix(buff, []byte("\xef\xbb\xbf"))
		events := make([]Event, 0)
		err = json.Unmarshal(buff, &events)
		if err != nil {
			log.Printf("Failed to unmarshal event: %s\n", err.Error())
		}

		for _, event := range events {
			eventType := event.Data["type"].(string)
			rawEvent, err := json.Marshal(event)
			if err != nil {
				log.Printf("Failed to marshal raw event: %s\n", err.Error())
				continue
			}
			//Form attributes
			attr := make([]attribute.KeyValue, 0)
			attr = append(attr, attribute.Int64("timestamp", event.Timestamp.Unix()))
			attr = append(attr, attribute.String("raw", string(rawEvent)))
			for k, v := range event.Data {
				switch v.(type) {
				case string:
					attr = append(attr, attribute.String(k, v.(string)))
					break
				}
			}
			_, span := DefaultTracer.Start(r.Context(), eventType, trace.WithAttributes(attr...))
			span.End()
			log.Printf("Saved event of type: %s\n", eventType)
		}

	})
	http.ListenAndServe(":5555", nil)
}
