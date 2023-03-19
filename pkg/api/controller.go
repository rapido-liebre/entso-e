package api

import (
	"entso-e_reports/pkg/common/config"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	// here you can set database instance if needed
	channels *config.Channels
	//cfgPath  string
}

func RegisterRoutes(app *fiber.App, ch *config.Channels) {
	h := &handler{
		channels: ch,
	}

	// routes configuration
	// example curl:
	// curl --request GET --url http://localhost:<port>/api/config
	routes := app.Group("/api")

	routes.Get("/update_config", h.UpdateConfig)

	routes.Get("/start", h.RunServices)

	routes.Get("/stop", h.StopServices)

	routes.Get("/config", h.GetConfig)

	routes.Get("/quit", h.QuitServices)
}

func (h handler) ChannelIsClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
