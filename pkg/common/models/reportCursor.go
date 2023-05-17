package models

import "time"

type CursorData struct {
	ReportType string
	Revision   int64
	Creator    string
	Created    time.Time
	Start      time.Time
	End        time.Time
	Saved      time.Time
	Reported   time.Time
}

func (cd CursorData) IsValid() bool {
	if len(cd.Creator) == 0 {
		return false
	}
	if cd.Start.IsZero() {
		return false
	}
	if cd.End.IsZero() {
		return false
	}
	return true
}

type CursorPayload struct {
	MrId                int64
	BusinessType        string
	FlowDirection       string
	QuantityMeasurement string
	Position            int64
	Quantity            float64
	SecondaryQuantity   int64
}

//func (cp CursorPayload) SaveToKjczReport(report *KjczReport) {
//	switch cp.MrId {
//	case 1:
//		report.MeanValue = append(report.MeanValue, cp.getReportPayload())
//	case 2:
//		report.StandardDeviation = append(report.StandardDeviation, cp.getReportPayload())
//	case 3:
//		report.Percentile1 = append(report.Percentile1, cp.getReportPayload())
//	case 4:
//		report.Percentile5 = append(report.Percentile5, cp.getReportPayload())
//	case 5:
//		report.Percentile10 = append(report.Percentile10, cp.getReportPayload())
//	case 6:
//		report.Percentile90 = append(report.Percentile90, cp.getReportPayload())
//	case 7:
//		report.Percentile95 = append(report.Percentile95, cp.getReportPayload())
//	case 8:
//		report.Percentile99 = append(report.Percentile99, cp.getReportPayload())
//	case 9:
//		report.FRCEOutsideLevel1RangeUp = append(report.FRCEOutsideLevel1RangeUp, cp.getReportPayload())
//	case 10:
//		report.FRCEOutsideLevel1RangeDown = append(report.FRCEOutsideLevel1RangeDown, cp.getReportPayload())
//	case 11:
//		report.FRCEOutsideLevel2RangeUp = append(report.FRCEOutsideLevel2RangeUp, cp.getReportPayload())
//	case 12:
//		report.FRCEOutsideLevel2RangeDown = append(report.FRCEOutsideLevel2RangeDown, cp.getReportPayload())
//	case 13:
//		report.FRCEExceeded60PercOfFRRCapacityUp = append(report.FRCEExceeded60PercOfFRRCapacityUp, cp.getReportPayload())
//	case 14:
//		report.FRCEExceeded60PercOfFRRCapacityDown = append(report.FRCEExceeded60PercOfFRRCapacityDown, cp.getReportPayload())
//	}
//}

func (cp CursorPayload) getReportPayload() ReportPayload {
	var sq *int
	sq = new(int)
	*sq = int(cp.SecondaryQuantity)

	return ReportPayload{
		MrId:                int(cp.MrId),
		BusinessType:        cp.BusinessType,
		FlowDirection:       cp.FlowDirection,
		QuantityMeasureUnit: cp.QuantityMeasurement,
		Position:            int(cp.Position),
		Quantity:            cp.Quantity,
		SecondaryQuantity:   sq,
	}
}
