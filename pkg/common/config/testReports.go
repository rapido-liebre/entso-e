package config

import (
	"entso-e_reports/pkg/common/models"
	"time"
)

func TestReportData() models.ReportData {
	tStart, _ := time.Parse(time.DateOnly, "2022-01-31")
	tEnd, _ := time.Parse(time.DateOnly, "2022-03-31")

	return models.ReportData{
		Creator: "Reksio",
		Start:   tStart,
		End:     tEnd,
	}
}

func TestKjczReportBody(reportId int64, data models.ReportData) models.KjczReport {

	var report models.KjczReport

	report.Data = data

	report.MeanValue = append(report.MeanValue, models.ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        models.MeanValue.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            3.309,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        models.MeanValue.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            1.388,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        models.MeanValue.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            -1.941,
		SecondaryQuantity:   nil,
	})

	report.StandardDeviation = append(report.StandardDeviation, models.ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        models.StandardDeviation.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            56.739,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        models.StandardDeviation.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            61.257,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                2,
		BusinessType:        models.StandardDeviation.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            58.645,
		SecondaryQuantity:   nil,
	})

	var sq *int
	sq = new(int)
	*sq = 1
	report.Percentile1 = append(report.Percentile1, models.ReportPayload{
		ReportId:            reportId,
		MrId:                3,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                3,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                3,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq,
	})

	*sq = 5
	report.Percentile5 = append(report.Percentile5, models.ReportPayload{
		ReportId:            reportId,
		MrId:                4,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                4,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                4,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq,
	})

	*sq = 10
	report.Percentile10 = append(report.Percentile10, models.ReportPayload{
		ReportId:            reportId,
		MrId:                5,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                5,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                5,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq,
	})

	*sq = 90
	report.Percentile90 = append(report.Percentile90, models.ReportPayload{
		ReportId:            reportId,
		MrId:                6,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                6,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                6,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq,
	})

	*sq = 95
	report.Percentile95 = append(report.Percentile95, models.ReportPayload{
		ReportId:            reportId,
		MrId:                7,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                7,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                7,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq,
	})

	*sq = 99
	report.Percentile99 = append(report.Percentile99, models.ReportPayload{
		ReportId:            reportId,
		MrId:                8,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            1,
		Quantity:            -132.749,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                8,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            2,
		Quantity:            -154.430,
		SecondaryQuantity:   sq,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                8,
		BusinessType:        models.Percentile.String(),
		FlowDirection:       models.UpAndDown.String(),
		QuantityMeasureUnit: models.MAW.String(),
		Position:            3,
		Quantity:            -162.567,
		SecondaryQuantity:   sq,
	})

	report.FRCEOutsideLevel1RangeUp = append(report.FRCEOutsideLevel1RangeUp, models.ReportPayload{
		ReportId:            reportId,
		MrId:                9,
		BusinessType:        models.FRCEOutsideLevel1Range.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            1,
		Quantity:            64,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                9,
		BusinessType:        models.FRCEOutsideLevel1Range.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            2,
		Quantity:            39,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                9,
		BusinessType:        models.FRCEOutsideLevel1Range.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            3,
		Quantity:            32,
		SecondaryQuantity:   nil,
	})

	report.FRCEOutsideLevel1RangeDown = append(report.FRCEOutsideLevel1RangeDown, models.ReportPayload{
		ReportId:            reportId,
		MrId:                10,
		BusinessType:        models.FRCEOutsideLevel1Range.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            1,
		Quantity:            28,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                10,
		BusinessType:        models.FRCEOutsideLevel1Range.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            2,
		Quantity:            50,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                10,
		BusinessType:        models.FRCEOutsideLevel1Range.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            3,
		Quantity:            51,
		SecondaryQuantity:   nil,
	})

	report.FRCEOutsideLevel2RangeUp = append(report.FRCEOutsideLevel2RangeUp, models.ReportPayload{
		ReportId:            reportId,
		MrId:                11,
		BusinessType:        models.FRCEOutsideLevel2Range.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            1,
		Quantity:            6,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                11,
		BusinessType:        models.FRCEOutsideLevel2Range.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            2,
		Quantity:            8,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                11,
		BusinessType:        models.FRCEOutsideLevel2Range.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            3,
		Quantity:            0,
		SecondaryQuantity:   nil,
	})

	report.FRCEOutsideLevel2RangeDown = append(report.FRCEOutsideLevel2RangeDown, models.ReportPayload{
		ReportId:            reportId,
		MrId:                12,
		BusinessType:        models.FRCEOutsideLevel2Range.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            1,
		Quantity:            3,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                12,
		BusinessType:        models.FRCEOutsideLevel2Range.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            2,
		Quantity:            8,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                12,
		BusinessType:        models.FRCEOutsideLevel2Range.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            3,
		Quantity:            10,
		SecondaryQuantity:   nil,
	})

	report.FRCEExceeded60PercOfFRRCapacityUp = append(report.FRCEExceeded60PercOfFRRCapacityUp, models.ReportPayload{
		ReportId:            reportId,
		MrId:                13,
		BusinessType:        models.FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            1,
		Quantity:            7,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                13,
		BusinessType:        models.FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            2,
		Quantity:            2,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                13,
		BusinessType:        models.FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       models.Up.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            3,
		Quantity:            0,
		SecondaryQuantity:   nil,
	})

	report.FRCEExceeded60PercOfFRRCapacityDown = append(report.FRCEExceeded60PercOfFRRCapacityDown, models.ReportPayload{
		ReportId:            reportId,
		MrId:                14,
		BusinessType:        models.FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            1,
		Quantity:            1,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                14,
		BusinessType:        models.FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            2,
		Quantity:            3,
		SecondaryQuantity:   nil,
	}, models.ReportPayload{
		ReportId:            reportId,
		MrId:                14,
		BusinessType:        models.FRCEExceeded60PercOfFRRCapacity.String(),
		FlowDirection:       models.Down.String(),
		QuantityMeasureUnit: models.C62.String(),
		Position:            3,
		Quantity:            6,
		SecondaryQuantity:   nil,
	})

	return report
}
