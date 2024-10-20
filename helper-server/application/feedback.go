package application

import (
	"encoding/json"
	"helper-server/internal/models"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Application) feedbackHandler(c echo.Context) error {
	payload := map[string]interface{}{}
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)
	if err != nil {
		log.Printf("Failed to parse body: %s\n", err.Error())
		return err
	}
	userid, err := strconv.ParseInt(payload["userid"].([]string)[0], 10, 64)
	if err != nil {
		log.Printf("Failed to parse body: %s\n", err.Error())
		return err
	}
	data := payload["data"].([]string)[0]

	dataMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		log.Printf("Failed to parse body: %s\n", err.Error())
		return err
	}

	sub := dataMap["Subject"].(string)
	msg := dataMap["Message"].(string)
	typ := int(dataMap["Type"].(float64))

	ply := &models.Player{
		SteamID: uint64(userid),
	}
	err = h.repo.AttachPlayer(ply)
	if err != nil {
		log.Printf("Failed to attach player: %s\n", err.Error())
		return err
	}

	err = h.repo.CreateFeedback(&models.PlayerFeedback{
		Player:       ply,
		Subject:      sub,
		Content:      msg,
		FeedbackType: typ,
	})
	if err != nil {
		log.Printf("Failed to create report: %s\n", err.Error())
		return err
	}

	return nil
}
