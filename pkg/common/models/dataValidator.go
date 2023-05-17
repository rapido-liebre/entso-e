package models

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

func ParseQueryParams(ctx *fiber.Ctx, expectedParamsCount int) ([]string, error) {
	fmt.Println(string(ctx.Request().URI().QueryString()))
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
	year, month, err := splitYearMonth(params[1])
	if err != nil {
		return time.Time{}, err
	}

	currentLocation := time.Now().Location()
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0)
	fmt.Println(firstOfMonth)
	fmt.Println(lastOfMonth)

	if isLastDay {
		return lastOfMonth, nil
	}
	return firstOfMonth, nil
}

func LocalTimeAsUTC(local time.Time) time.Time {
	_, offset := local.Zone()
	fmt.Println(offset)

	return local.Add(time.Second * time.Duration(offset))
}

func GetReportDates(startDate, endDate string) (startDt, endDt time.Time, err error) {
	if startDt, err = firstDayDate(startDate); err != nil {
		return time.Time{}, time.Time{}, err
	}
	if endDt, err = lastDayDate(endDate); err != nil {
		return time.Time{}, time.Time{}, err
	}
	//extend end date for 1 day forward
	endDt = endDt.AddDate(0, 0, 1)
	return
}

func firstDayDate(yearMonth string) (time.Time, error) {
	year, month, err := splitYearMonth(yearMonth)
	if err != nil {
		return time.Time{}, err
	}

	dt := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	return dt, nil
}

func lastDayDate(yearMonth string) (time.Time, error) {
	year, month, err := splitYearMonth(yearMonth)
	if err != nil {
		return time.Time{}, err
	}

	var lastDay = 31
	if isEven(month) {
		lastDay = 30
		if month == 2 {
			lastDay = 28
			if isLeapYear(year) {
				lastDay = 29
			}
		}
	}
	dt := time.Date(year, time.Month(month), lastDay, 0, 0, 0, 0, time.Now().Location())
	return dt, nil
}

func splitYearMonth(yearMonth string) (year, month int, err error) {
	ym := strings.Split(yearMonth, "-")
	if len(ym) != 2 {
		return 0, 0, fmt.Errorf("invalid yearMonth param %s", yearMonth)
	}
	if year, err = strconv.Atoi(ym[0]); err != nil {
		return 0, 0, fmt.Errorf("invalid year %s err:%s", ym[0], err)
	}
	if month, err = strconv.Atoi(ym[1]); err != nil {
		return 0, 0, fmt.Errorf("invalid month %s err:%s", ym[1], err)
	}
	return year, month, nil
}

func isEven(n int) bool {
	months := []int{2, 4, 6, 9, 11}

	if contains(months, n) {
		return true //has 30 days
	}
	return false
}

func isLeapYear(y int) bool {
	if y%4 == 0 && y%100 != 0 || y%400 == 0 {
		return true //is leap
	}
	return false
}

func contains(s []int, n int) bool {
	for _, v := range s {
		if v == n {
			return true
		}
	}
	return false
}

func GetSecondaryQuantityString(secondaryQuantity *int) string {
	if secondaryQuantity == nil {
		return "null"
	}
	return strconv.Itoa(*secondaryQuantity)
}
