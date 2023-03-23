package api

import (
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"github.com/gofiber/fiber/v2"
)

func (h handler) ConnectToDB(ctx *fiber.Ctx) error {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       false,
		ConnectionOnly: true,
		ReportType:     0,
		Payload:        "",
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h handler) SendTestKjcz(ctx *fiber.Ctx) error {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       true,
		ConnectionOnly: false,
		ReportType:     models.Kjcz,
		Payload:        "",
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h handler) SendTestKjczAndPublish(ctx *fiber.Ctx) error {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        true,
		TestData:       true,
		ConnectionOnly: false,
		ReportType:     models.Kjcz,
		Payload:        "",
	}
	return ctx.SendStatus(fiber.StatusOK)
}
