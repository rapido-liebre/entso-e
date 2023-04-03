package models

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

func ParseQueryParams(ctx *fiber.Ctx, expectedParamsCount int) ([]string, error) {
	params := strings.Split(string(ctx.Request().URI().QueryString()), "&")
	if len(params) != expectedParamsCount {
		return []string{}, fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("Invalid request params %s", ctx.Request().URI().QueryString()))
	}
	return params, nil
}

func ExtractDate(param string, isLastDay bool) (time.Time, error) {
	params := strings.Split(param, "=")
	if len(params) != 2 {
		return time.Time{}, fmt.Errorf("invalid request params %s", param)
	}
	yearMonth := strings.Split(params[1], "-")
	if len(yearMonth) != 2 {
		return time.Time{}, fmt.Errorf("invalid request params %s", params[1])
	}
	currentLocation := time.Now().Location()
	year, err := strconv.Atoi(yearMonth[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid request params %s err:%s", yearMonth[0], err)
	}
	month, err := strconv.Atoi(yearMonth[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid request params %s err:%s", yearMonth[1], err)
	}

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	//date, err := time.Parse(time.DateOnly, lastOfMonth)
	//if err != nil {
	//	return time.Time{}, err
	//}
	if isLastDay {
		return lastOfMonth, nil
	}
	return firstOfMonth, nil
}
