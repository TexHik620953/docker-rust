package application

import (
	"context"
	"helper-server/internal/config"
	"helper-server/internal/repos"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	ctx context.Context
	cfg *config.AppConfig

	repo *repos.Repository

	api *echo.Echo
}

func New(ctx context.Context, cfg *config.AppConfig) (*Application, error) {
	h := &Application{
		ctx: ctx,
		cfg: cfg,
		api: echo.New(),
	}
	var err error
	h.repo, err = repos.New(cfg)
	if err != nil {
		return nil, err
	}

	h.api.Use(middleware.Recover())
	h.api.Any("/event", h.eventHandler)
	h.api.Any("/feedback", h.feedbackHandler)
	h.api.Any("/report", h.reportHandler)

	h.api.Static("/static", "./static")

	{
		balance_group := h.api.Group("/balance")
		balance_group.POST("/:steam_id/withdraw", h.withdrawBalance)
		balance_group.POST("/:steam_id/deposit", h.depositBalance)
		balance_group.GET("/:steam_id", h.getBalance)
	}

	return h, nil
}

func (h *Application) Start() {
	err := h.api.Start(":5555")
	if err != nil {
		log.Fatal(err)
	}
}
