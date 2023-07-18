package models

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"time"
)

type ReportCalculator struct {
	data15min map[Year]map[time.Month][]LfcAce
	data1min  map[Year]map[time.Month][]LfcAce
}

func (rc *ReportCalculator) Calculate(lfcAce15 []LfcAce, lfcAce1 []LfcAce, extraParams map[string]string) KjczReport {
	if rc.data15min == nil {
		rc.data15min = make(map[Year]map[time.Month][]LfcAce)
	}
	if rc.data1min == nil {
		rc.data1min = make(map[Year]map[time.Month][]LfcAce)
	}

	rc.splitToYearsMonths(lfcAce15, FETCH_15_MIN)
	rc.splitToYearsMonths(lfcAce1, FETCH_1_MIN)

	var level1, level2, excCapacityUp, excCapacityDown *float64
	if len(extraParams["level1"]) > 0 {
		level1 = new(float64)
		*level1, _ = strconv.ParseFloat(extraParams["level1"], 64)
	}
	if len(extraParams["level2"]) > 0 {
		level2 = new(float64)
		*level2, _ = strconv.ParseFloat(extraParams["level2"], 64)
	}
	if len(extraParams["excCapacityUp"]) > 0 {
		excCapacityUp = new(float64)
		*excCapacityUp, _ = strconv.ParseFloat(extraParams["excCapacityUp"], 64)
	}
	if len(extraParams["excCapacityDown"]) > 0 {
		excCapacityDown = new(float64)
		*excCapacityDown, _ = strconv.ParseFloat(extraParams["excCapacityDown"], 64)
	}

	// get year-month(s)
	var yearMonths []string
	for year, months := range rc.data15min {
		for month, _ := range months {
			yearMonth := fmt.Sprintf("%d-%02d", year, month)
			yearMonths = append(yearMonths, yearMonth)
		}
	}
	sort.Strings(yearMonths)

	calculations := make(map[string]*KjczBody)

	// calculate avg15m
	for year, months := range rc.data15min {
		position := 1
		for month, measurements := range months {
			kjczBody := &KjczBody{}
			yearMonth := fmt.Sprintf("%d-%02d", year, month)
			CalculateReportData15min(measurements, position, kjczBody, yearMonth, level1, level2)
			calculations[yearMonth] = kjczBody
			position += 1
		}
	}
	// calculate avg1m
	for year, months := range rc.data1min {
		position := 1
		for month, measurements := range months {
			yearMonth := fmt.Sprintf("%d-%02d", year, month)
			CalculateReportData1min(measurements, position, calculations[yearMonth], yearMonth, excCapacityUp, excCapacityDown)
			position += 1
		}
	}

	report := fillReport(calculations, yearMonths)

	return report
}

func fillReport(calculations map[string]*KjczBody, yearMonths []string) KjczReport {
	var report KjczReport
	startDate, _ := time.Parse(time.DateOnly, fmt.Sprintf("%s-01", yearMonths[0]))
	report.Data = TestReportData(PR_SO_KJCZ, startDate)
	report.Data.YearMonths = yearMonths

	for _, yearMonth := range yearMonths {
		var secondaryQuantity *int
		secondaryQuantity = new(int)
		report.MeanValue = append(report.MeanValue, getReportPayload(1, MeanValue.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].MeanValue)...)
		report.StandardDeviation = append(report.StandardDeviation, getReportPayload(2, StandardDeviation.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].StandardDeviation)...)
		*secondaryQuantity = 1
		report.Percentile1 = append(report.Percentile1, getReportPayload(3, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].Percentile1)...)
		*secondaryQuantity = 5
		report.Percentile5 = append(report.Percentile5, getReportPayload(4, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].Percentile5)...)
		*secondaryQuantity = 10
		report.Percentile10 = append(report.Percentile10, getReportPayload(5, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].Percentile10)...)
		*secondaryQuantity = 90
		report.Percentile90 = append(report.Percentile90, getReportPayload(6, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].Percentile90)...)
		*secondaryQuantity = 95
		report.Percentile95 = append(report.Percentile95, getReportPayload(7, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].Percentile95)...)
		*secondaryQuantity = 99
		report.Percentile99 = append(report.Percentile99, getReportPayload(8, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, calculations[yearMonth].Percentile99)...)
		secondaryQuantity = new(int)
		report.FRCEOutsideLevel1RangeUp = append(report.FRCEOutsideLevel1RangeUp, getReportPayload(9, FRCEOutsideLevel1Range.String(), Up.String(), C62.String(), secondaryQuantity, calculations[yearMonth].FrceOutsideLevel1RangeUp)...)
		report.FRCEOutsideLevel1RangeDown = append(report.FRCEOutsideLevel1RangeDown, getReportPayload(10, FRCEOutsideLevel1Range.String(), Down.String(), C62.String(), secondaryQuantity, calculations[yearMonth].FrceOutsideLevel1RangeDown)...)
		report.FRCEOutsideLevel2RangeUp = append(report.FRCEOutsideLevel2RangeUp, getReportPayload(11, FRCEOutsideLevel2Range.String(), Up.String(), C62.String(), secondaryQuantity, calculations[yearMonth].FrceOutsideLevel2RangeUp)...)
		report.FRCEOutsideLevel2RangeDown = append(report.FRCEOutsideLevel2RangeDown, getReportPayload(12, FRCEOutsideLevel2Range.String(), Down.String(), C62.String(), secondaryQuantity, calculations[yearMonth].FrceOutsideLevel2RangeDown)...)
		report.FRCEExceeded60PercOfFRRCapacityUp = append(report.FRCEExceeded60PercOfFRRCapacityUp, getReportPayload(13, FRCEExceeded60PercOfFRRCapacity.String(), Up.String(), C62.String(), secondaryQuantity, calculations[yearMonth].FrceExceeded60PercOfFRRCapacityUp)...)
		report.FRCEExceeded60PercOfFRRCapacityDown = append(report.FRCEExceeded60PercOfFRRCapacityDown, getReportPayload(14, FRCEExceeded60PercOfFRRCapacity.String(), Down.String(), C62.String(), secondaryQuantity, calculations[yearMonth].FrceExceeded60PercOfFRRCapacityDown)...)
	}
	return report
}

func (rc *ReportCalculator) splitToYearsMonths(lfcAce []LfcAce, rt ReportType) {
	switch rt {
	case FETCH_15_MIN:
		for _, v := range lfcAce {
			_, exists := rc.data15min[Year(v.AvgTime.Year())]
			if !exists {
				rc.data15min[Year(v.AvgTime.Year())] = make(map[time.Month][]LfcAce)
			}

			rc.data15min[Year(v.AvgTime.Year())][v.AvgTime.Month()] =
				append(rc.data15min[Year(v.AvgTime.Year())][v.AvgTime.Month()], v)
		}
	case FETCH_1_MIN:
		for _, v := range lfcAce {
			_, exists := rc.data1min[Year(v.AvgTime.Year())]
			if !exists {
				rc.data1min[Year(v.AvgTime.Year())] = make(map[time.Month][]LfcAce)
			}

			rc.data1min[Year(v.AvgTime.Year())][v.AvgTime.Month()] =
				append(rc.data1min[Year(v.AvgTime.Year())][v.AvgTime.Month()], v)
		}
	default:
		log.Fatalf("Wrong report type! Expected: %s or %s, got: %s", FETCH_15_MIN.String(), FETCH_1_MIN.String(), rt.String())
	}
}

func CalculateReportData15min(lfcAce15 []LfcAce, position int, body *KjczBody, yearMonth string, lev1, lev2 *float64) {
	var level1, level2 float64
	level1 = 124.964
	level2 = 236.326
	if lev1 != nil {
		level1 = *lev1
	}
	if lev2 != nil {
		level2 = *lev2
	}

	totalCount := float64(len(lfcAce15))
	fmt.Println(totalCount)

	var (
		sum, sumq, lv1pos, lv1neg, lv2pos, lv2neg    float64
		p1, p5, p10, p90, p95, p99                   float64
		perc1, perc5, perc10, perc90, perc95, perc99 float64
	)

	for _, v := range lfcAce15 {
		avgValue := -v.AvgValue
		sum += avgValue

		if avgValue > level1 {
			lv1pos += 1
		}
		if avgValue < -level1 {
			lv1neg += 1
		}
		if avgValue > level2 {
			lv2pos += 1
		}
		if avgValue < -level2 {
			lv2neg += 1
		}
	}

	avg := sum / totalCount
	fmt.Println("Avg:", avg)

	for _, v := range lfcAce15 {
		sumq += math.Pow(v.AvgValue-avg, 2)
	}

	dev := math.Sqrt(sumq / totalCount)
	fmt.Println("Deviation:", dev)

	fmt.Println("Level1 +:", lv1pos)
	fmt.Println("Level1 -:", lv1neg)
	fmt.Println("Level2 +:", lv2pos)
	fmt.Println("Level2 -:", lv2neg)

	var vals []float64
	for _, v := range lfcAce15 {
		vals = append(vals, v.AvgValue)
	}
	sort.Float64s(vals)

	// percentile indexes
	p1 = math.Round(1 * totalCount / float64(100))
	p5 = math.Round(5 * totalCount / float64(100))
	p10 = math.Round(10 * totalCount / float64(100))
	p90 = math.Round(90 * totalCount / float64(100))
	p95 = math.Round(95 * totalCount / float64(100))
	p99 = math.Round(99 * totalCount / float64(100))

	for i, v := range vals {
		if i == int(p1) {
			perc1 = v
		}
		if i == int(p5) {
			perc5 = v
		}
		if i == int(p10) {
			perc10 = v
		}
		if i == int(p90) {
			perc90 = v
		}
		if i == int(p95) {
			perc95 = v
		}
		if i == int(p99) {
			perc99 = v
		}
	}

	fmt.Println("Perc 1:", perc1, p1)
	fmt.Println("Perc 5:", perc5, p5)
	fmt.Println("Perc 10:", perc10, p10)
	fmt.Println("Perc 90:", perc90, p90)
	fmt.Println("Perc 95:", perc95, p95)
	fmt.Println("Perc 99:", perc99, p99)

	getBRPayload := func(quantity float64) BodyReportPayload {
		return BodyReportPayload{
			Position:  position,
			Quantity:  roundFloat(quantity, 3),
			YearMonth: yearMonth,
		}
	}

	body.MeanValue = append(body.MeanValue, getBRPayload(avg))
	body.StandardDeviation = append(body.StandardDeviation, getBRPayload(dev))
	body.Percentile1 = append(body.Percentile1, getBRPayload(perc1))
	body.Percentile5 = append(body.Percentile5, getBRPayload(perc5))
	body.Percentile10 = append(body.Percentile10, getBRPayload(perc10))
	body.Percentile90 = append(body.Percentile90, getBRPayload(perc90))
	body.Percentile95 = append(body.Percentile95, getBRPayload(perc95))
	body.Percentile99 = append(body.Percentile99, getBRPayload(perc99))
	body.FrceOutsideLevel1RangeUp = append(body.FrceOutsideLevel1RangeUp, getBRPayload(lv1pos))
	body.FrceOutsideLevel1RangeDown = append(body.FrceOutsideLevel1RangeDown, getBRPayload(lv1neg))
	body.FrceOutsideLevel2RangeUp = append(body.FrceOutsideLevel2RangeUp, getBRPayload(lv2pos))
	body.FrceOutsideLevel2RangeDown = append(body.FrceOutsideLevel2RangeDown, getBRPayload(lv2neg))
}

func CalculateReportData1min(lfcAce1 []LfcAce, position int, body *KjczBody, yearMonth string, excCapacityUp, excCapacityDown *float64) {
	//const FRR = 1075
	//const FRR60 = FRR * 0.6
	//const FRR15 = FRR * 0.15

	var FRRpos, FRRneg float64
	FRRpos = 1075
	FRRneg = -1075
	if excCapacityUp != nil {
		FRRpos = *excCapacityUp
	}
	if excCapacityDown != nil {
		FRRneg = *excCapacityDown
	}

	FRR60pos := FRRpos * 0.6
	FRR15pos := FRRpos * 0.15
	FRR60neg := FRRneg * 0.6
	FRR15neg := FRRneg * 0.15

	totalCount := float64(len(lfcAce1))
	fmt.Println(totalCount)

	var (
		exceeding, exceedingTime, plus, minus int
	)

	for _, v := range lfcAce1 {
		//fmt.Println(k, v)
		val := v.AvgValue
		{
			if -val < FRR15pos && exceeding == 1 {
				if exceedingTime > 14 {
					plus += 1
				}
				exceeding = 0
				exceedingTime = 0
			}

			if -val > FRR15neg && exceeding == -1 {
				if exceedingTime > 14 {
					minus += 1
				}
				exceeding = 0
				exceedingTime = 0
			}

			if -val > FRR60pos && exceeding == 0 {
				exceeding = 1
			}

			if -val < FRR60neg && exceeding == 0 {
				exceeding = -1
			}

			if exceeding != 0 {
				exceedingTime += 1
			}
		}
	}

	fmt.Println("60% Exc Plus:", plus)
	fmt.Println("60% Exc Minus:", minus)

	getBRPayload := func(quantity float64) BodyReportPayload {
		return BodyReportPayload{
			Position:  position,
			Quantity:  quantity,
			YearMonth: yearMonth,
		}
	}

	body.FrceExceeded60PercOfFRRCapacityUp = append(body.FrceExceeded60PercOfFRRCapacityUp, getBRPayload(float64(plus)))
	body.FrceExceeded60PercOfFRRCapacityDown = append(body.FrceExceeded60PercOfFRRCapacityDown, getBRPayload(float64(minus)))
}

func getReportPayload(mrId int, businessType, flowDirection, quantityMeasureUnit string, secondaryQuantity *int, bodyPayloads []BodyReportPayload) []ReportPayload {
	var reportPayloads []ReportPayload

	for _, bodyPayload := range bodyPayloads {
		reportPayloads = append(reportPayloads, ReportPayload{
			ReportId:            0,
			MrId:                mrId,
			BusinessType:        businessType,
			FlowDirection:       flowDirection,
			QuantityMeasureUnit: quantityMeasureUnit,
			Position:            bodyPayload.Position,
			Quantity:            bodyPayload.Quantity,
			SecondaryQuantity:   secondaryQuantity,
		})
	}

	return reportPayloads
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
