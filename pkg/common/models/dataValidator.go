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

func FirstDayDate(yearMonth string) (time.Time, error) {
	return time.Parse(time.DateOnly, strings.Join([]string{yearMonth, "01"}, "-"))
}

func LastDayDate(yearMonth string) (time.Time, error) {
	var lastDay = "31"
	ym := strings.Split(yearMonth, "-")
	month, _ := strconv.Atoi(ym[1])

	if isEven(month) {
		lastDay = "30"
		if month == 2 {
			lastDay = "28"
			year, _ := strconv.Atoi(ym[0])
			if isLeapYear(year) {
				lastDay = "29"
			}
		}
	}

	return time.Parse(time.DateOnly, strings.Join([]string{yearMonth, lastDay}, "-"))
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
