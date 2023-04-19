package models

import (
	"math/rand"
	"time"
)

type Report interface {
	KjczReport | PzrrReport | PzfrrReport
}

var creators = []string{
	"Mikołaj Kopernik",
	"Maria Curie-Skłodowska",
	"Jan Heweliusz",
	"Henryk Arctowski",
	"Ernest Malinowski",
	"Kazimierz Funk",
	"Ludwig Zamenhoff",
	"Korczak Ziółkowski",
}
var dateRanges = []string{
	"2022-01-01",
	"2022-03-31",
	"2022-12-01",
	"2023-12-31",
}

func TestReportData(rt ReportType) ReportData {
	index := 0
	monthsDuration := 3

	if rt != PR_SO_KJCZ {
		index = 2
		monthsDuration = 4
	}
	tStart, _ := time.Parse(time.DateOnly, dateRanges[index])
	tEnd, _ := time.Parse(time.DateOnly, dateRanges[index+1])

	return ReportData{
		Creator:        creators[rand.Intn(len(creators))],
		Revision:       0,
		Start:          tStart,
		End:            tEnd,
		Created:        time.Time{},
		Saved:          time.Time{},
		Reported:       time.Time{},
		MonthsDuration: int64(monthsDuration),
	}
}

//func GetTestReportBody(reportId int64, data ReportData, rt ReportType) REPORT {
//	switch rt {
//	case PR_SO_KJCZ:
//		var kjcz KjczReport
//		return kjcz.GetTestReport(reportId, data)
//	case PD_BI_PZRR:
//		var pzrr PzrrReport
//		return pzrr.GetTestReport(reportId, data)
//	case PD_BI_PZFRR:
//		var pzfrr PzfrrReport
//		return pzfrr.GetTestReport(reportId, data)
//	}
//	var r REPORT
//	return r
//}

func GetTestKjczReportBody(reportId int64, data ReportData) KjczReport {

	var report KjczReport
	report.Data = data

	var sq1, sq5, sq10, sq90, sq95, sq99, dummy *int
	sq1 = new(int)
	sq5 = new(int)
	sq10 = new(int)
	sq90 = new(int)
	sq95 = new(int)
	sq99 = new(int)
	dummy = new(int)

	report.MeanValue = append(report.MeanValue, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        MeanValue.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            3.309,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        MeanValue.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            1.388,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        MeanValue.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            -1.941,
		SecondaryQuantity:   dummy,
	})

	report.StandardDeviation = append(report.StandardDeviation, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        StandardDeviation.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            56.739,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        StandardDeviation.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            61.257,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        StandardDeviation.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            58.645,
		SecondaryQuantity:   dummy,
	})

	*sq1 = 1
	report.Percentile1 = append(report.Percentile1, ReportPayload{
		ReportId:            reportId,
		MrId:                3,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq1,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                3,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq1,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                3,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq1,
	})

	*sq5 = 5
	report.Percentile5 = append(report.Percentile5, ReportPayload{
		ReportId:            reportId,
		MrId:                4,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq5,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                4,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq5,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                4,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq5,
	})

	*sq10 = 10
	report.Percentile10 = append(report.Percentile10, ReportPayload{
		ReportId:            reportId,
		MrId:                5,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq10,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                5,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq10,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                5,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq10,
	})

	*sq90 = 90
	report.Percentile90 = append(report.Percentile90, ReportPayload{
		ReportId:            reportId,
		MrId:                6,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq90,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                6,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq90,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                6,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq90,
	})

	*sq95 = 95
	report.Percentile95 = append(report.Percentile95, ReportPayload{
		ReportId:            reportId,
		MrId:                7,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq95,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                7,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq95,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                7,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq95,
	})

	*sq99 = 99
	report.Percentile99 = append(report.Percentile99, ReportPayload{
		ReportId:            reportId,
		MrId:                8,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq99,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                8,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq99,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                8,
		BusinessType:        Percentile.String(),
		FlowDirection:       UpAndDown.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq99,
	})

	report.FRCEOutsideLevel1RangeUp = append(report.FRCEOutsideLevel1RangeUp, ReportPayload{
		ReportId:            reportId,
		MrId:                9,
		BusinessType:        FRCEOutsideLevel1Range.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            1,
		Quantity:            64,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                9,
		BusinessType:        FRCEOutsideLevel1Range.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            2,
		Quantity:            39,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                9,
		BusinessType:        FRCEOutsideLevel1Range.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            3,
		Quantity:            32,
		SecondaryQuantity:   dummy,
	})

	report.FRCEOutsideLevel1RangeDown = append(report.FRCEOutsideLevel1RangeDown, ReportPayload{
		ReportId:            reportId,
		MrId:                10,
		BusinessType:        FRCEOutsideLevel1Range.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            1,
		Quantity:            28,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                10,
		BusinessType:        FRCEOutsideLevel1Range.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            2,
		Quantity:            50,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                10,
		BusinessType:        FRCEOutsideLevel1Range.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            3,
		Quantity:            51,
		SecondaryQuantity:   dummy,
	})

	report.FRCEOutsideLevel2RangeUp = append(report.FRCEOutsideLevel2RangeUp, ReportPayload{
		ReportId:            reportId,
		MrId:                11,
		BusinessType:        FRCEOutsideLevel2Range.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            1,
		Quantity:            6,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                11,
		BusinessType:        FRCEOutsideLevel2Range.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            2,
		Quantity:            8,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                11,
		BusinessType:        FRCEOutsideLevel2Range.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            3,
		Quantity:            0,
		SecondaryQuantity:   dummy,
	})

	report.FRCEOutsideLevel2RangeDown = append(report.FRCEOutsideLevel2RangeDown, ReportPayload{
		ReportId:            reportId,
		MrId:                12,
		BusinessType:        FRCEOutsideLevel2Range.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            1,
		Quantity:            3,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                12,
		BusinessType:        FRCEOutsideLevel2Range.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            2,
		Quantity:            8,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                12,
		BusinessType:        FRCEOutsideLevel2Range.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            3,
		Quantity:            10,
		SecondaryQuantity:   dummy,
	})

	report.FRCEExceeded60PercOfFRRCapacityUp = append(report.FRCEExceeded60PercOfFRRCapacityUp, ReportPayload{
		ReportId:            reportId,
		MrId:                13,
		BusinessType:        FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            1,
		Quantity:            7,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                13,
		BusinessType:        FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            2,
		Quantity:            2,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                13,
		BusinessType:        FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            3,
		Quantity:            0,
		SecondaryQuantity:   dummy,
	})

	report.FRCEExceeded60PercOfFRRCapacityDown = append(report.FRCEExceeded60PercOfFRRCapacityDown, ReportPayload{
		ReportId:            reportId,
		MrId:                14,
		BusinessType:        FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            1,
		Quantity:            1,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                14,
		BusinessType:        FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            2,
		Quantity:            3,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                14,
		BusinessType:        FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: C62.String(),
		Position:            3,
		Quantity:            6,
		SecondaryQuantity:   dummy,
	})

	return report
}

func GetTestPzrrReportBody(reportId int64, data ReportData) PzrrReport {

	var report PzrrReport
	report.Data = data

	var dummy *int
	dummy = new(int)

	report.ForecastedCapacityUp = append(report.ForecastedCapacityUp, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            1500.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            1500.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            1500.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            4,
		Quantity:            1500.0,
		SecondaryQuantity:   dummy,
	})

	report.ForecastedCapacityDown = append(report.ForecastedCapacityDown, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            0.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            0.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            0.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            4,
		Quantity:            0.0,
		SecondaryQuantity:   dummy,
	})

	return report
}

func GetTestPzfrrReportBody(reportId int64, data ReportData) PzfrrReport {

	var report PzfrrReport
	report.Data = data

	var dummy *int
	dummy = new(int)

	report.ForecastedCapacityUp = append(report.ForecastedCapacityUp, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            1075.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            1075.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            1075.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Up.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            4,
		Quantity:            1075.0,
		SecondaryQuantity:   dummy,
	})

	report.ForecastedCapacityDown = append(report.ForecastedCapacityDown, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            1,
		Quantity:            600.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            2,
		Quantity:            600.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            3,
		Quantity:            600.0,
		SecondaryQuantity:   dummy,
	}, ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        ForecastedCapacity.String(),
		FlowDirection:       Down.String(),
		QuantityMeasureUnit: MAW.String(),
		Position:            4,
		Quantity:            600.0,
		SecondaryQuantity:   dummy,
	})

	return report
}
