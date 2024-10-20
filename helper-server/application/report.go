package application

import (
	"helper-server/internal/models"
	"log"

	"github.com/labstack/echo/v4"
)

type report struct {
	Reporter uint64 `json:"reporter_steamid"`
	Target   uint64 `json:"target_steamid"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func (h *Application) reportHandler(c echo.Context) error {
	payload := &report{}
	err := (&echo.DefaultBinder{}).BindBody(c, payload)
	if err != nil {
		log.Printf("Failed to parse body: %s\n", err.Error())
		return err
	}

	reporter := &models.Player{
		SteamID: payload.Reporter,
	}
	err = h.repo.AttachPlayer(reporter)
	if err != nil {
		log.Printf("Failed to attach player: %s\n", err.Error())
		return err
	}

	target := &models.Player{
		SteamID: payload.Target,
	}
	err = h.repo.AttachPlayer(target)
	if err != nil {
		log.Printf("Failed to attach player: %s\n", err.Error())
		return err
	}

	err = h.repo.CreateReport(&models.PlayerReport{
		SourcePlayer: reporter,
		TargetPlayer: target,
		Subject:      payload.Subject,
		Content:      payload.Message,
	})
	if err != nil {
		log.Printf("Failed to create report: %s\n", err.Error())
		return err
	}

	return nil
}
