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

func (h handler) SendTest(ctx *fiber.Ctx, rt models.ReportType, publish bool) error {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        publish,
		TestData:       true,
		ConnectionOnly: false,
		ReportType:     rt,
		Payload:        "",
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h handler) SendTestKjcz(ctx *fiber.Ctx) error {
	return h.SendTest(ctx, models.PR_SO_KJCZ, false)
}

func (h handler) SendTestPzrr(ctx *fiber.Ctx) error {
	return h.SendTest(ctx, models.PD_BI_PZRR, false)
}

func (h handler) SendTestPzfrr(ctx *fiber.Ctx) error {
	return h.SendTest(ctx, models.PD_BI_PZFRR, false)
}

func (h handler) SendTestKjczAndPublish(ctx *fiber.Ctx) error {
	return h.SendTest(ctx, models.PR_SO_KJCZ, true)
}

func (h handler) SendTestPzrrAndPublish(ctx *fiber.Ctx) error {
	return h.SendTest(ctx, models.PD_BI_PZRR, true)
}

func (h handler) SendTestPzfrrAndPublish(ctx *fiber.Ctx) error {
	return h.SendTest(ctx, models.PD_BI_PZFRR, true)
}
