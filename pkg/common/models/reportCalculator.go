package models

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

type ReportCalculator struct {
	data15min map[Year]map[time.Month][]LfcAce
	data1min  map[Year]map[time.Month][]LfcAce
}

func (rc *ReportCalculator) Calculate(lfcAce15 []LfcAce, lfcAce1 []LfcAce) KjczReport {
	if rc.data15min == nil {
		rc.data15min = make(map[Year]map[time.Month][]LfcAce)
	}
	if rc.data1min == nil {
		rc.data1min = make(map[Year]map[time.Month][]LfcAce)
	}

	rc.splitToYearsMonths(lfcAce15, FETCH_15_MIN)
	rc.splitToYearsMonths(lfcAce1, FETCH_1_MIN)

	kjczBody := &KjczBody{}

	for year, months := range rc.data15min {
		position := 1
		for month, measurements := range months {
			yearMonth := strings.Join([]string{string(year), month.String()}, "-")
			CalculateReportData15min(measurements, position, kjczBody, yearMonth)
			CalculateReportData1min(measurements, position, kjczBody, yearMonth)
			position += 1
		}
	}

	var report KjczReport
	report.Data = TestReportData(PR_SO_KJCZ)
	var secondaryQuantity *int
	secondaryQuantity = new(int)
	report.MeanValue = getReportPayload(1, MeanValue.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.MeanValue)
	report.StandardDeviation = getReportPayload(2, StandardDeviation.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.StandardDeviation)
	*secondaryQuantity = 1
	report.Percentile1 = getReportPayload(3, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.Percentile1)
	*secondaryQuantity = 5
	report.Percentile5 = getReportPayload(4, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.Percentile5)
	*secondaryQuantity = 10
	report.Percentile10 = getReportPayload(5, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.Percentile10)
	*secondaryQuantity = 90
	report.Percentile90 = getReportPayload(6, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.Percentile90)
	*secondaryQuantity = 95
	report.Percentile95 = getReportPayload(7, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.Percentile95)
	*secondaryQuantity = 99
	report.Percentile99 = getReportPayload(8, Percentile.String(), UpAndDown.String(), MAW.String(), secondaryQuantity, kjczBody.Percentile99)
	secondaryQuantity = new(int)
	report.FRCEOutsideLevel1RangeUp = getReportPayload(9, FRCEOutsideLevel1Range.String(), Up.String(), C62.String(), secondaryQuantity, kjczBody.FrceOutsideLevel1RangeUp)
	report.FRCEOutsideLevel1RangeDown = getReportPayload(10, FRCEOutsideLevel1Range.String(), Down.String(), C62.String(), secondaryQuantity, kjczBody.FrceOutsideLevel1RangeDown)
	report.FRCEOutsideLevel2RangeUp = getReportPayload(11, FRCEOutsideLevel2Range.String(), Up.String(), C62.String(), secondaryQuantity, kjczBody.FrceOutsideLevel2RangeUp)
	report.FRCEOutsideLevel2RangeDown = getReportPayload(12, FRCEOutsideLevel2Range.String(), Down.String(), C62.String(), secondaryQuantity, kjczBody.FrceOutsideLevel2RangeDown)
	report.FRCEExceeded60PercOfFRRCapacityUp = getReportPayload(13, FRCEExceeded60PercOfFRRCapacity.String(), Up.String(), C62.String(), secondaryQuantity, kjczBody.FrceExceeded60PercOfFRRCapacityUp)
	report.FRCEExceeded60PercOfFRRCapacityDown = getReportPayload(14, FRCEExceeded60PercOfFRRCapacity.String(), Down.String(), C62.String(), secondaryQuantity, kjczBody.FrceExceeded60PercOfFRRCapacityDown)

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

func CalculateReportData15min(lfcAce15 []LfcAce, position int, body *KjczBody, yearMonth string) {
	const level1 = 124.964
	const level2 = 236.326
	totalCount := float64(len(lfcAce15))
	fmt.Println(totalCount)

	var (
		sum, sumq, lv1pos, lv1neg, lv2pos, lv2neg, perc1, perc5, perc10, perc90, perc95, perc99 float64
	)

	for _, v := range lfcAce15 {
		//fmt.Println(k, v)
		sum += v.AvgValue
		sumq += math.Pow(v.AvgValue, 2)

		if -v.AvgValue > level1 {
			lv1pos += 1
		}
		if -v.AvgValue < -level1 {
			lv1neg += 1
		}
		if -v.AvgValue > level2 {
			lv2pos += 1
		}
		if -v.AvgValue < -level2 {
			lv2neg += 1
		}

	}
	avg := sum / totalCount
	avg = -avg
	fmt.Println(avg)

	dev := math.Sqrt(sumq / totalCount)
	fmt.Println(dev)

	fmt.Println("Level1 +:", lv1pos)
	fmt.Println("Level1 -:", lv1neg)
	fmt.Println("Level2 +:", lv2pos)
	fmt.Println("Level2 -:", lv2neg)

	perc1 = totalCount / float64(100*(100-1))
	perc5 = totalCount / float64(100*(100-5))
	perc10 = totalCount / float64(100*(100-10))
	perc90 = totalCount / float64(100*(100-90))
	perc95 = totalCount / float64(100*(100-95))
	perc99 = totalCount / float64(100*(100-99))
	fmt.Println("Perc 1:", perc1)
	fmt.Println("Perc 5:", perc5)
	fmt.Println("Perc 10:", perc10)
	fmt.Println("Perc 90:", perc90)
	fmt.Println("Perc 95:", perc95)
	fmt.Println("Perc 99:", perc99)

	//out0 := totalCount
	//out1 := avg
	//out2 := sumq
	//out3 := perc1
	//out4 := perc1
	//out5 := perc1
	//out6 := perc1
	//out7 := perc1
	//out8 := perc1
	//out9 := lv1pos
	//out10 := lv1neg
	//out11 := lv2pos
	//out12 := lv2neg
	//out13 := 0 //data od

	getBRPayload := func(quantity float64) BodyReportPayload {
		return BodyReportPayload{
			Position:  position,
			Quantity:  quantity,
			YearMonth: yearMonth,
		}
	}

	body.MeanValue = append(body.MeanValue, getBRPayload(avg))
	body.StandardDeviation = append(body.StandardDeviation, getBRPayload(sumq))
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

func CalculateReportData1min(lfcAce1 []LfcAce, position int, body *KjczBody, yearMonth string) {
	const FRR = 1075
	const FRR60 = FRR * 0.6
	const FRR15 = FRR * 0.15

	totalCount := float64(len(lfcAce1))
	fmt.Println(totalCount)

	var (
		exceeding, exceedingTime, plus, minus int
	)

	for k, v := range lfcAce1 {
		fmt.Println(k, v)
		val := v.AvgValue
		{
			if -val < FRR15 && exceeding == 1 {
				if exceedingTime > 14 {
					plus += 1
				}
				exceeding = 0
				exceedingTime = 0
			}

			if -val > -FRR15 && exceeding == -1 {
				if exceedingTime > 14 {
					minus += 1
				}
				exceeding = 0
				exceedingTime = 0
			}

			if -val > FRR60 && exceeding == 0 {
				exceeding = 1
			}

			if -val < -FRR60 && exceeding == 0 {
				exceeding = -1
			}

			if exceeding != 0 {
				exceedingTime += 1
			}
		}
	}

	//out14 := plus
	//out15 := minus
	//out16 := totalCount

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
