package application

import (
	"encoding/json"
	"fmt"
	"helper-server/internal/models"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

type report struct {
	Reporter string `json:"reporter_steamid"`
	Target   string `json:"target_steamid"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func (h *Application) reportHandler(c echo.Context) error {
	/*
		text, _ := io.ReadAll(c.Request().Body)
		t1 := string(text)
		fmt.Println(t1)
	*/

	payload := &report{}
	err := json.NewDecoder(c.Request().Body).Decode(payload)
	if err != nil {
		log.Printf("Failed to parse body: %s\n", err.Error())
		return err
	}
	fmt.Println(payload.Reporter)
	fmt.Println(payload.Target)
	fmt.Println(payload.Subject)
	fmt.Println(payload.Message)

	reporterId, err := strconv.ParseUint(payload.Reporter, 10, 64)
	if err != nil {
		log.Printf("Failed to parse reporter_steamid: %s\n", err.Error())
		return err
	}
	targetId, err := strconv.ParseUint(payload.Target, 10, 64)
	if err != nil {
		log.Printf("Failed to parse target_steamid: %s\n", err.Error())
		return err
	}

	fmt.Printf("received eport from %d to %d\n", reporterId, targetId)

	reporter := &models.Player{
		SteamID: reporterId,
	}
	err = h.repo.AttachPlayer(reporter)
	if err != nil {
		log.Printf("Failed to attach player: %s\n", err.Error())
		return err
	}

	target := &models.Player{
		SteamID: targetId,
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
