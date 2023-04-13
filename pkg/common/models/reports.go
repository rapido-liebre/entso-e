package models

import (
	"errors"
	"time"
)

type Reporter interface {
	GetAllPayloads() []ReportPayload
	Save(cd CursorData, cps []CursorPayload)
	//GetTestReport(reportId int64, data ReportData) any
}

type ReportData struct {
	Creator        string
	Revision       int64
	Start          time.Time
	End            time.Time
	Created        time.Time
	Saved          time.Time
	Reported       time.Time
	MonthsDuration int64
}

func (rd ReportData) GetDurationInMonths() (int64, error) {
	duration := rd.End.Sub(rd.Start)
	difference := int64(duration.Hours() / 24 / 30)
	if difference <= 0 {
		return 0, errors.New("invalid date range")
	}
	return difference, nil
}

type ReportPayload struct {
	ReportId            int64
	MrId                int
	BusinessType        string
	FlowDirection       string
	QuantityMeasureUnit string
	Position            int
	Quantity            float64
	SecondaryQuantity   *int
}

type KjczReport struct {
	Data                                ReportData
	MeanValue                           []ReportPayload
	StandardDeviation                   []ReportPayload
	Percentile1                         []ReportPayload
	Percentile5                         []ReportPayload
	Percentile10                        []ReportPayload
	Percentile90                        []ReportPayload
	Percentile95                        []ReportPayload
	Percentile99                        []ReportPayload
	FRCEOutsideLevel1RangeUp            []ReportPayload
	FRCEOutsideLevel1RangeDown          []ReportPayload
	FRCEOutsideLevel2RangeUp            []ReportPayload
	FRCEOutsideLevel2RangeDown          []ReportPayload
	FRCEExceeded60PercOfFRRCapacityUp   []ReportPayload
	FRCEExceeded60PercOfFRRCapacityDown []ReportPayload
}

func (r *KjczReport) GetAllPayloads() (payloads []ReportPayload) {
	payloads = append(payloads, r.MeanValue...)
	payloads = append(payloads, r.StandardDeviation...)
	payloads = append(payloads, r.Percentile1...)
	payloads = append(payloads, r.Percentile5...)
	payloads = append(payloads, r.Percentile10...)
	payloads = append(payloads, r.Percentile90...)
	payloads = append(payloads, r.Percentile95...)
	payloads = append(payloads, r.Percentile99...)
	payloads = append(payloads, r.FRCEOutsideLevel1RangeUp...)
	payloads = append(payloads, r.FRCEOutsideLevel1RangeDown...)
	payloads = append(payloads, r.FRCEOutsideLevel2RangeUp...)
	payloads = append(payloads, r.FRCEOutsideLevel2RangeDown...)
	payloads = append(payloads, r.FRCEExceeded60PercOfFRRCapacityUp...)
	payloads = append(payloads, r.FRCEExceeded60PercOfFRRCapacityDown...)

	return
}

func (r *KjczReport) Save(cd CursorData, cps []CursorPayload) {
	r.Data.Creator = cd.Creator
	r.Data.Revision = cd.Revision
	r.Data.Start = cd.Start
	r.Data.End = cd.End
	r.Data.Created = cd.Created
	r.Data.Saved = cd.Saved
	r.Data.Reported = cd.Reported

	for _, cp := range cps {
		switch cp.MrId {
		case 1:
			r.MeanValue = append(r.MeanValue, cp.getReportPayload())
		case 2:
			r.StandardDeviation = append(r.StandardDeviation, cp.getReportPayload())
		case 3:
			r.Percentile1 = append(r.Percentile1, cp.getReportPayload())
		case 4:
			r.Percentile5 = append(r.Percentile5, cp.getReportPayload())
		case 5:
			r.Percentile10 = append(r.Percentile10, cp.getReportPayload())
		case 6:
			r.Percentile90 = append(r.Percentile90, cp.getReportPayload())
		case 7:
			r.Percentile95 = append(r.Percentile95, cp.getReportPayload())
		case 8:
			r.Percentile99 = append(r.Percentile99, cp.getReportPayload())
		case 9:
			r.FRCEOutsideLevel1RangeUp = append(r.FRCEOutsideLevel1RangeUp, cp.getReportPayload())
		case 10:
			r.FRCEOutsideLevel1RangeDown = append(r.FRCEOutsideLevel1RangeDown, cp.getReportPayload())
		case 11:
			r.FRCEOutsideLevel2RangeUp = append(r.FRCEOutsideLevel2RangeUp, cp.getReportPayload())
		case 12:
			r.FRCEOutsideLevel2RangeDown = append(r.FRCEOutsideLevel2RangeDown, cp.getReportPayload())
		case 13:
			r.FRCEExceeded60PercOfFRRCapacityUp = append(r.FRCEExceeded60PercOfFRRCapacityUp, cp.getReportPayload())
		case 14:
			r.FRCEExceeded60PercOfFRRCapacityDown = append(r.FRCEExceeded60PercOfFRRCapacityDown, cp.getReportPayload())
		}
	}
}

func updatePayloads(dest []ReportPayload, src []BodyReportPayload) {
	for _, d := range dest {
		for _, s := range src {
			if d.Position == s.Position {
				d.Quantity = s.Quantity
			}
		}
	}
}

func (r *KjczReport) Update(payload any) {
	p := payload.(KjczBody)

	r.Data.Creator = p.Data.Creator
	r.Data.Start, _ = FirstDayDate(p.Data.Start)
	r.Data.End, _ = LastDayDate(p.Data.End)

	updatePayloads(r.MeanValue, p.MeanValue)
	updatePayloads(r.StandardDeviation, p.StandardDeviation)
	updatePayloads(r.Percentile1, p.Percentile1)
	updatePayloads(r.Percentile5, p.Percentile5)
	updatePayloads(r.Percentile10, p.Percentile10)
	updatePayloads(r.Percentile90, p.Percentile90)
	updatePayloads(r.Percentile95, p.Percentile95)
	updatePayloads(r.Percentile99, p.Percentile99)
	updatePayloads(r.FRCEOutsideLevel1RangeUp, p.FrceOutsideLevel1RangeUp)
	updatePayloads(r.FRCEOutsideLevel1RangeDown, p.FrceOutsideLevel1RangeDown)
	updatePayloads(r.FRCEOutsideLevel2RangeUp, p.FrceOutsideLevel2RangeUp)
	updatePayloads(r.FRCEOutsideLevel2RangeDown, p.FrceOutsideLevel2RangeDown)
	updatePayloads(r.FRCEExceeded60PercOfFRRCapacityUp, p.FrceExceeded60PercOfFRRCapacityUp)
	updatePayloads(r.FRCEExceeded60PercOfFRRCapacityDown, p.FrceExceeded60PercOfFRRCapacityDown)
}

//func (r KjczReport) GetTestReport(reportId int64, data ReportData) any {
//	return GetTestKjczReportBody(reportId, data)
//}

type PzrrReport struct {
	Data                   ReportData
	ForecastedCapacityUp   []ReportPayload
	ForecastedCapacityDown []ReportPayload
}

func (r *PzrrReport) GetAllPayloads() (payloads []ReportPayload) {
	payloads = append(payloads, r.ForecastedCapacityUp...)
	payloads = append(payloads, r.ForecastedCapacityDown...)

	return
}

func (r *PzrrReport) Save(cd CursorData, cps []CursorPayload) {
	r.Data.Creator = cd.Creator
	r.Data.Revision = cd.Revision
	r.Data.Start = cd.Start
	r.Data.End = cd.End
	r.Data.Created = cd.Created
	r.Data.Saved = cd.Saved
	r.Data.Reported = cd.Reported

	for _, cp := range cps {
		switch cp.MrId {
		case 1:
			r.ForecastedCapacityUp = append(r.ForecastedCapacityUp, cp.getReportPayload())
		case 2:
			r.ForecastedCapacityDown = append(r.ForecastedCapacityDown, cp.getReportPayload())
		}
	}
}

func (r *PzrrReport) Update(payload any) {
	p := payload.(PzrrBody)

	r.Data.Creator = p.Data.Creator
	r.Data.Start, _ = FirstDayDate(p.Data.Start)
	r.Data.End, _ = LastDayDate(p.Data.End)

	updatePayloads(r.ForecastedCapacityUp, p.ForecastedCapacityUp)
	updatePayloads(r.ForecastedCapacityDown, p.ForecastedCapacityDown)
}

//func (r PzrrReport) GetTestReport(reportId int64, data ReportData) any {
//	return GetTestPzrrReportBody(reportId, data)
//}

type PzfrrReport struct {
	Data                   ReportData
	ForecastedCapacityUp   []ReportPayload
	ForecastedCapacityDown []ReportPayload
}

func (r *PzfrrReport) GetAllPayloads() (payloads []ReportPayload) {
	payloads = append(payloads, r.ForecastedCapacityUp...)
	payloads = append(payloads, r.ForecastedCapacityDown...)

	return
}

func (r *PzfrrReport) Save(cd CursorData, cps []CursorPayload) {
	r.Data.Creator = cd.Creator
	r.Data.Revision = cd.Revision
	r.Data.Start = cd.Start
	r.Data.End = cd.End
	r.Data.Created = cd.Created
	r.Data.Saved = cd.Saved
	r.Data.Reported = cd.Reported

	for _, cp := range cps {
		switch cp.MrId {
		case 1:
			r.ForecastedCapacityUp = append(r.ForecastedCapacityUp, cp.getReportPayload())
		case 2:
			r.ForecastedCapacityDown = append(r.ForecastedCapacityDown, cp.getReportPayload())
		}
	}
}

func (r *PzfrrReport) Update(payload any) {
	p := payload.(PzfrrBody)

	r.Data.Creator = p.Data.Creator
	r.Data.Start, _ = FirstDayDate(p.Data.Start)
	r.Data.End, _ = LastDayDate(p.Data.End)

	updatePayloads(r.ForecastedCapacityUp, p.ForecastedCapacityUp)
	updatePayloads(r.ForecastedCapacityDown, p.ForecastedCapacityDown)
}

//func (r PzfrrReport) GetTestReport(reportId int64, data ReportData) any {
//	return GetTestPzfrrReportBody(reportId, data)
//}

type FlowDirection int

const (
	Up FlowDirection = iota
	Down
	UpAndDown
)

func (fd FlowDirection) String() string {
	return []string{"A01", "A02", "A03"}[fd]
}

type BusinessType int

const (
	MeanValue BusinessType = iota
	StandardDeviation
	Percentile
	FRCEOutsideLevel1Range
	FRCEOutsideLevel2Range
	FRCEExceeded60PercOfFRRCapacity
	ForecastedCapacity
)

func (bt BusinessType) String() string {
	return []string{"C66", "C67", "C68", "C71", "C72", "C73", "C76"}[bt]
}

type QuantityMeasureUnit int

const (
	MAW QuantityMeasureUnit = iota
	C62
)

func (qmu QuantityMeasureUnit) String() string {
	return []string{"MAW", "C62"}[qmu]
}
