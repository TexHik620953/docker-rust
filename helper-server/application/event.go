package application

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
)

type event struct {
	Timestamp time.Time              `json:"Timestamp"`
	Data      map[string]interface{} `json:"Data"`
}

func (h *Application) eventHandler(c echo.Context) error {
	key := c.Request().Header.Get("X-Api-Key")
	if key != h.cfg.AnalyticsSecret {
		log.Printf("Invalid X-Api-Key: %s\n", key)
		return c.String(http.StatusUnauthorized, "")
	}
	payload := []event{}
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)
	if err != nil {
		log.Printf("Failed to parse body: %s\n", err.Error())
		return err
	}
	//buff = bytes.TrimPrefix(buff, []byte("\xef\xbb\xbf"))

	for _, event := range payload {
		//eventType := event.Data["type"].(string)
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
			}
		}
		//_, span := DefaultTracer.Start(r.Context(), eventType, trace.WithAttributes(attr...))
		//span.End()
	}
	return nil
}
