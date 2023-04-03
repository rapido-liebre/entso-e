package api

import (
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"fmt"
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

func (h handler) GetKjczReport(rd models.ReportData) (models.KjczReport, error) {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       false,
		ConnectionOnly: false,
		ReportType:     models.PR_SO_KJCZ,
		ReportData:     rd,
		Payload:        "",
	}
	report := <-h.channels.KjczReport
	return report, nil
}

func (h handler) GetPzrrReport(rd models.ReportData) (models.PzrrReport, error) {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       false,
		ConnectionOnly: false,
		ReportType:     models.PD_BI_PZRR,
		ReportData:     rd,
		Payload:        "",
	}
	report := <-h.channels.PzrrReport
	return report, nil
}

func (h handler) GetPzfrrReport(rd models.ReportData) (models.PzfrrReport, error) {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       false,
		ConnectionOnly: false,
		ReportType:     models.PD_BI_PZFRR,
		ReportData:     rd,
		Payload:        "",
	}
	report := <-h.channels.PzfrrReport
	return report, nil
}

func getCommonReportData(ctx *fiber.Ctx) (models.ReportData, error) {
	params, err := models.ParseQueryParams(ctx, 2)
	if err != nil {
		return models.ReportData{}, err
	}
	fmt.Println(params)

	//extract the dates
	isLastDay := false
	dateFrom, err := models.ExtractDate(params[0], isLastDay)
	if err != nil {
		return models.ReportData{}, err
	}
	isLastDay = true
	dateTo, err := models.ExtractDate(params[1], isLastDay)
	if err != nil {
		return models.ReportData{}, err
	}

	rd := models.ReportData{
		Start: dateFrom,
		End:   dateTo,
	}
	rd.MonthsDuration, err = rd.GetDurationInMonths()
	if err != nil {
		return models.ReportData{}, err
	}

	return rd, nil
}

func (h handler) GetKjcz(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	report, err := h.GetKjczReport(rd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) GetPzrr(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	report, err := h.GetPzrrReport(rd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) GetPzfrr(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	report, err := h.GetPzfrrReport(rd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(report)
}
