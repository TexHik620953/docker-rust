package application

import (
	"helper-server/internal/models"
	"helper-server/utils"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getBalanceRequest struct {
	SteamID uint64 `param:"steam_id"`
}
type getBalanceResponse struct {
	Balance int64 `json:"balance"`
}

func (h *Application) getBalance(c echo.Context) error {
	rq := &getBalanceRequest{}
	err := (&echo.DefaultBinder{}).BindPathParams(c, rq)
	if err != nil {
		return err
	}

	ply := &models.Player{
		SteamID: rq.SteamID,
	}

	err = h.repo.AttachPlayer(ply)
	if err != nil {
		log.Printf("Failed to attach player: %s\n", err.Error())
		return err
	}
	return c.JSON(http.StatusOK, &getBalanceResponse{
		Balance: ply.Balance,
	})
}

type withdrawBalanceRequest struct {
	SteamID uint64 `param:"steam_id"`
	Amount  uint64 `query:"amount"`
}
type withdrawBalanceResponse struct {
	Error      *string `json:"error,omitempty"`
	ErrorCode  int     `json:"error_code,omitempty"`
	NewBalance int64   `json:"balance"`
}

func (h *Application) withdrawBalance(c echo.Context) error {
	rq := &withdrawBalanceRequest{}
	err := (&echo.DefaultBinder{}).BindPathParams(c, rq)
	if err != nil {
		return err
	}
	err = (&echo.DefaultBinder{}).BindQueryParams(c, rq)
	if err != nil {
		return err
	}

	ply := &models.Player{
		SteamID: rq.SteamID,
	}
	err = h.repo.AttachPlayer(ply)
	if err != nil {
		log.Printf("Failed to attach player: %s\n", err.Error())
		return err
	}

	ply.Balance -= int64(rq.Amount)
	if ply.Balance < 0 {
		return c.JSON(http.StatusOK, &withdrawBalanceResponse{
			Error:      utils.RefValue("Insufficient balance"),
			ErrorCode:  1,
			NewBalance: ply.Balance + int64(rq.Amount),
		})
	}

	err = h.repo.UpdateBalance(ply)
	if err != nil {
		log.Printf("Failed to update balance: %s\n", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, &withdrawBalanceResponse{
		NewBalance: ply.Balance,
	})
}

type depositBalanceRequest struct {
	SteamID uint64 `param:"steam_id"`
	Amount  uint64 `query:"amount"`
}
type depositBalanceResponse struct {
	NewBalance int64 `json:"balance"`
}

func (h *Application) depositBalance(c echo.Context) error {
	rq := &depositBalanceRequest{}
	err := (&echo.DefaultBinder{}).BindPathParams(c, rq)
	if err != nil {
		return err
	}
	err = (&echo.DefaultBinder{}).BindQueryParams(c, rq)
	if err != nil {
		return err
	}

	ply := &models.Player{
		SteamID: rq.SteamID,
	}
	err = h.repo.AttachPlayer(ply)
	if err != nil {
		log.Printf("Failed to attach player: %s\n", err.Error())
		return err
	}

	ply.Balance += int64(rq.Amount)

	err = h.repo.UpdateBalance(ply)
	if err != nil {
		log.Printf("Failed to update balance: %s\n", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, &depositBalanceResponse{
		NewBalance: ply.Balance,
	})
}
