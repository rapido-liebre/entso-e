package api

import (
	"encoding/json"
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
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

func (h handler) Fetch15min(ctx *fiber.Ctx) error {
	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       false,
		ConnectionOnly: true,
		ReportType:     models.FETCH_15_MIN,
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
	rd, err := getCommonReportData(ctx, 6)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       true,
		ConnectionOnly: false,
		ReportType:     models.PR_SO_KJCZ,
		ReportData:     rd,
		Payload:        "",
	}
	report := <-h.channels.KjczReport

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) SendTestPzrr(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx, 2)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       true,
		ConnectionOnly: false,
		ReportType:     models.PD_BI_PZRR,
		ReportData:     rd,
		Payload:        "",
	}
	report := <-h.channels.PzrrReport

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) SendTestPzfrr(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx, 2)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.channels.RunDBConn <- config.DBAction{
		Publish:        false,
		TestData:       true,
		ConnectionOnly: false,
		ReportType:     models.PD_BI_PZFRR,
		ReportData:     rd,
		Payload:        "",
	}
	report := <-h.channels.PzfrrReport

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
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
		Payload:        nil,
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
		Payload:        nil,
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
		Payload:        nil,
	}
	report := <-h.channels.PzfrrReport
	return report, nil
}

func getCommonReportData(ctx *fiber.Ctx, expectedParamsCount int) (models.ReportData, error) {
	params, err := models.ParseQueryParams(ctx, expectedParamsCount)
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
		Start:       dateFrom,
		End:         dateTo,
		ExtraParams: map[string]string{},
	}
	for i := 2; i < len(params); i++ {
		ep := strings.Split(params[i], "=")
		rd.ExtraParams[ep[0]] = ep[1]
	}

	return rd, nil
}

func (h handler) GetKjcz(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx, 2)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	report, err := h.GetKjczReport(rd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if report.Data.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, report.Data.Error.Error())
	}

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) GetPzrr(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx, 2)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	report, err := h.GetPzrrReport(rd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if report.Data.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, report.Data.Error.Error())
	}

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) GetPzfrr(ctx *fiber.Ctx) error {
	rd, err := getCommonReportData(ctx, 2)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	report, err := h.GetPzfrrReport(rd)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if report.Data.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, report.Data.Error.Error())
	}

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) saveKjczReport(ctx *fiber.Ctx, publish bool) error {
	fmt.Println(string(ctx.Body()))
	var (
		body   models.KjczBody
		tStart time.Time
		tEnd   time.Time
		err    error
	)

	if err = ctx.BodyParser(&body); err != nil {
		return err
	}
	if tStart, tEnd, err = models.GetReportDates(body.Data.Start, body.Data.End); err != nil {
		return err
	}

	rd := models.ReportData{
		Creator: body.Data.Creator,
		Start:   tStart,
		End:     tEnd,
	}

	h.channels.RunDBConn <- config.DBAction{
		Publish:        publish,
		TestData:       false,
		ConnectionOnly: false,
		ReportType:     models.PR_SO_KJCZ,
		ReportData:     rd,
		Payload:        body,
	}

	report := <-h.channels.KjczReport
	if report.Data.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, report.Data.Error.Error())
	}

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) SaveKjcz(ctx *fiber.Ctx) error {
	return h.saveKjczReport(ctx, false)
}

func (h handler) SaveKjczAndPublish(ctx *fiber.Ctx) error {
	return h.saveKjczReport(ctx, true)
}

func (h handler) savePzrrReport(ctx *fiber.Ctx, publish bool) error {
	fmt.Println(string(ctx.Body()))
	var (
		body   models.PzrrBody
		tStart time.Time
		tEnd   time.Time
		err    error
	)

	if err = ctx.BodyParser(&body); err != nil {
		return err
	}
	if tStart, tEnd, err = models.GetReportDates(body.Data.Start, body.Data.End); err != nil {
		return err
	}

	rd := models.ReportData{
		Creator: body.Data.Creator,
		Start:   tStart,
		End:     tEnd,
	}

	h.channels.RunDBConn <- config.DBAction{
		Publish:        publish,
		TestData:       false,
		ConnectionOnly: false,
		ReportType:     models.PD_BI_PZRR,
		ReportData:     rd,
		Payload:        body,
	}

	report := <-h.channels.PzrrReport
	if report.Data.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, report.Data.Error.Error())
	}

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) SavePzrr(ctx *fiber.Ctx) error {
	return h.savePzrrReport(ctx, false)
}

func (h handler) SavePzrrAndPublish(ctx *fiber.Ctx) error {
	return h.savePzrrReport(ctx, true)
}

func (h handler) savePzfrrReport(ctx *fiber.Ctx, publish bool) error {
	fmt.Println(string(ctx.Body()))
	var (
		body   models.PzfrrBody
		tStart time.Time
		tEnd   time.Time
		err    error
	)

	if err = ctx.BodyParser(&body); err != nil {
		return err
	}
	if tStart, tEnd, err = models.GetReportDates(body.Data.Start, body.Data.End); err != nil {
		return err
	}

	rd := models.ReportData{
		Creator: body.Data.Creator,
		Start:   tStart,
		End:     tEnd,
	}

	h.channels.RunDBConn <- config.DBAction{
		Publish:        publish,
		TestData:       false,
		ConnectionOnly: false,
		ReportType:     models.PD_BI_PZFRR,
		ReportData:     rd,
		Payload:        body,
	}

	report := <-h.channels.PzfrrReport
	if report.Data.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, report.Data.Error.Error())
	}

	rs, _ := json.Marshal(report)
	fmt.Println(fmt.Sprintf("%s", rs))
	return ctx.Status(fiber.StatusOK).JSON(report)
}

func (h handler) SavePzfrr(ctx *fiber.Ctx) error {
	return h.savePzfrrReport(ctx, false)
}

func (h handler) SavePzfrrAndPublish(ctx *fiber.Ctx) error {
	return h.savePzfrrReport(ctx, true)
}
