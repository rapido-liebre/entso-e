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

	routes.Get("/test_conn", h.ConnectToDB)

	routes.Get("/test_fetch_15", h.Fetch15min)

	routes.Get("/test_kjcz", h.SendTestKjcz)

	routes.Get("/test_kjcz_publish", h.SendTestKjczAndPublish)

	routes.Get("/test_pzrr", h.SendTestPzrr)

	routes.Get("/test_pzrr_publish", h.SendTestPzrrAndPublish)

	routes.Get("/test_pzfrr", h.SendTestPzfrr)

	routes.Get("/test_pzfrr_publish", h.SendTestPzfrrAndPublish)

	routes.Get("/get_kjcz", h.GetKjcz)

	routes.Get("/get_pzrr", h.GetPzrr)

	routes.Get("/get_pzfrr", h.GetPzfrr)

	routes.Post("/save_kjcz", h.SaveKjcz)

	routes.Post("/save_kjcz_publish", h.SaveKjczAndPublish)

	routes.Post("/save_pzrr", h.SavePzrr)

	routes.Post("/save_pzrr_publish", h.SavePzrrAndPublish)

	routes.Post("/save_pzfrr", h.SavePzfrr)

	routes.Post("/save_pzfrr_publish", h.SavePzfrrAndPublish)
}

func (h handler) ChannelIsClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
