package models

import (
	"reflect"
	"time"
)

type Reporter interface {
	GetAllPayloads() []ReportPayload
	SaveCursors(cd CursorData, cps []CursorPayload)
	//GetTestReport(reportId int64, data ReportData) any
}

type ReportData struct {
	Creator     string
	Revision    int64
	Start       time.Time
	End         time.Time
	Created     time.Time
	Saved       time.Time
	Reported    time.Time
	YearMonths  []string
	Error       error
	ExtraParams map[string]string
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

func payloadsAreEqual(payload1, payload2 []ReportPayload) bool {
	if len(payload1) != len(payload2) {
		return false
	}

	for _, p1 := range payload1 {
		for _, p2 := range payload2 {
			if p2.Position != p1.Position {
				continue
			}
			p1.ReportId = 0
			p2.ReportId = 0
			p1.SecondaryQuantity = nil
			p2.SecondaryQuantity = nil

			if !reflect.DeepEqual(p1, p2) {
				return false
			}
		}
	}

	return true
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

func (r *KjczReport) SaveCursors(cd CursorData, cps []CursorPayload) {
	r.Data.Creator = cd.Creator
	r.Data.Revision = cd.Revision
	r.Data.Start = LocalTimeAsUTC(cd.Start)
	r.Data.End = LocalTimeAsUTC(cd.End)
	r.Data.Created = cd.Created
	r.Data.Saved = cd.Saved
	if cd.Reported.Compare(time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)) != 0 {
		r.Data.Reported = cd.Reported
	}
	r.Data.YearMonths = calculateYearMonths(PR_SO_KJCZ, r.Data.Start)

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

func updatePayloads(dest *[]ReportPayload, src []BodyReportPayload) {
	for _, d := range *dest {
		for _, s := range src {
			if d.Position == s.Position {
				d.Quantity = s.Quantity
			}
		}
	}
}

func getPayloads(template []ReportPayload, src []BodyReportPayload) (dest []ReportPayload) {
	getSecQuantity := func(sq *int) *int {
		if sq == nil {
			return nil
		}
		if *sq == 0 {
			return nil
		}
		return sq
	}

	for _, s := range src {
		dest = append(dest, ReportPayload{
			ReportId:            0,
			MrId:                template[0].MrId,
			BusinessType:        template[0].BusinessType,
			FlowDirection:       template[0].FlowDirection,
			QuantityMeasureUnit: template[0].QuantityMeasureUnit,
			Position:            s.Position,
			Quantity:            s.Quantity,
			SecondaryQuantity:   getSecQuantity(template[0].SecondaryQuantity),
		})
	}
	return
}

func (r *KjczReport) Update(payload any) (areChanges bool) {
	p := payload.(KjczBody)

	if r.Data.Creator != p.Data.Creator {
		r.Data.Creator = p.Data.Creator
		areChanges = true
	}

	//if r.Data.Start.IsZero() || r.Data.End.IsZero() {
	dateStart, dateEnd, _ := GetReportDates(p.Data.Start, p.Data.End)
	if r.Data.Start != dateStart || r.Data.End != dateEnd {
		r.Data.Start = dateStart
		r.Data.End = dateEnd
		areChanges = true
	}
	r.Data.YearMonths = calculateYearMonths(PR_SO_KJCZ, r.Data.Start)
	//}

	t := GetKjczReportTemplate(r.Data)

	updatePayloads := func(template []ReportPayload, src []BodyReportPayload, dst *[]ReportPayload) (diff bool) {
		payloads := getPayloads(template, src)
		if !payloadsAreEqual(*dst, payloads) {
			diff = true
		}
		*dst = payloads
		return
	}

	if updatePayloads(t.MeanValue, p.MeanValue, &r.MeanValue) {
		areChanges = true
	}
	if updatePayloads(t.StandardDeviation, p.StandardDeviation, &r.StandardDeviation) {
		areChanges = true
	}
	if updatePayloads(t.Percentile1, p.Percentile1, &r.Percentile1) {
		areChanges = true
	}
	if updatePayloads(t.Percentile5, p.Percentile5, &r.Percentile5) {
		areChanges = true
	}
	if updatePayloads(t.Percentile10, p.Percentile10, &r.Percentile10) {
		areChanges = true
	}
	if updatePayloads(t.Percentile90, p.Percentile90, &r.Percentile90) {
		areChanges = true
	}
	if updatePayloads(t.Percentile95, p.Percentile95, &r.Percentile95) {
		areChanges = true
	}
	if updatePayloads(t.Percentile99, p.Percentile99, &r.Percentile99) {
		areChanges = true
	}
	if updatePayloads(t.FRCEOutsideLevel1RangeUp, p.FrceOutsideLevel1RangeUp, &r.FRCEOutsideLevel1RangeUp) {
		areChanges = true
	}
	if updatePayloads(t.FRCEOutsideLevel1RangeDown, p.FrceOutsideLevel1RangeDown, &r.FRCEOutsideLevel1RangeDown) {
		areChanges = true
	}
	if updatePayloads(t.FRCEOutsideLevel2RangeUp, p.FrceOutsideLevel2RangeUp, &r.FRCEOutsideLevel2RangeUp) {
		areChanges = true
	}
	if updatePayloads(t.FRCEOutsideLevel2RangeDown, p.FrceOutsideLevel2RangeDown, &r.FRCEOutsideLevel2RangeDown) {
		areChanges = true
	}
	if updatePayloads(t.FRCEExceeded60PercOfFRRCapacityUp, p.FrceExceeded60PercOfFRRCapacityUp, &r.FRCEExceeded60PercOfFRRCapacityUp) {
		areChanges = true
	}
	if updatePayloads(t.FRCEExceeded60PercOfFRRCapacityDown, p.FrceExceeded60PercOfFRRCapacityDown, &r.FRCEExceeded60PercOfFRRCapacityDown) {
		areChanges = true
	}

	return
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

func (r *PzrrReport) SaveCursors(cd CursorData, cps []CursorPayload) {
	r.Data.Creator = cd.Creator
	r.Data.Revision = cd.Revision
	r.Data.Start = LocalTimeAsUTC(cd.Start)
	r.Data.End = LocalTimeAsUTC(cd.End)
	r.Data.Created = cd.Created
	r.Data.Saved = cd.Saved
	if cd.Reported.Compare(time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)) != 0 {
		r.Data.Reported = cd.Reported
	}
	r.Data.YearMonths = calculateYearMonths(PD_BI_PZRR, r.Data.Start)

	for _, cp := range cps {
		switch cp.MrId {
		case 1:
			r.ForecastedCapacityUp = append(r.ForecastedCapacityUp, cp.getReportPayload())
		case 2:
			r.ForecastedCapacityDown = append(r.ForecastedCapacityDown, cp.getReportPayload())
		}
	}
}

func (r *PzrrReport) Update(payload any) (areChanges bool) {
	p := payload.(PzrrBody)

	if r.Data.Creator != p.Data.Creator {
		r.Data.Creator = p.Data.Creator
		areChanges = true
	}

	//if r.Data.Start.IsZero() || r.Data.End.IsZero() {
	dateStart, dateEnd, _ := GetReportDates(p.Data.Start, p.Data.End)
	if r.Data.Start != dateStart || r.Data.End != dateEnd {
		r.Data.Start = dateStart
		r.Data.End = dateEnd
		areChanges = true
	}
	r.Data.YearMonths = calculateYearMonths(PD_BI_PZRR, r.Data.Start)
	//}

	t := GetPzrrReportTemplate(r.Data)

	forecastedCapacityUp := getPayloads(t.ForecastedCapacityUp, p.ForecastedCapacityUp)
	if !payloadsAreEqual(r.ForecastedCapacityUp, forecastedCapacityUp) {
		areChanges = true
	}
	r.ForecastedCapacityUp = forecastedCapacityUp

	forecastedCapacityDown := getPayloads(t.ForecastedCapacityDown, p.ForecastedCapacityDown)
	if !payloadsAreEqual(r.ForecastedCapacityDown, forecastedCapacityDown) {
		areChanges = true
	}
	r.ForecastedCapacityDown = forecastedCapacityDown

	return
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

func (r *PzfrrReport) SaveCursors(cd CursorData, cps []CursorPayload) {
	r.Data.Creator = cd.Creator
	r.Data.Revision = cd.Revision
	r.Data.Start = LocalTimeAsUTC(cd.Start)
	r.Data.End = LocalTimeAsUTC(cd.End)
	r.Data.Created = cd.Created
	r.Data.Saved = cd.Saved
	if cd.Reported.Compare(time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local)) != 0 {
		r.Data.Reported = cd.Reported
	}
	r.Data.YearMonths = calculateYearMonths(PD_BI_PZFRR, r.Data.Start)

	for _, cp := range cps {
		switch cp.MrId {
		case 1:
			r.ForecastedCapacityUp = append(r.ForecastedCapacityUp, cp.getReportPayload())
		case 2:
			r.ForecastedCapacityDown = append(r.ForecastedCapacityDown, cp.getReportPayload())
		}
	}
}

func (r *PzfrrReport) Update(payload any) (areChanges bool) {
	p := payload.(PzfrrBody)

	if r.Data.Creator != p.Data.Creator {
		r.Data.Creator = p.Data.Creator
		areChanges = true
	}

	if r.Data.Start.IsZero() || r.Data.End.IsZero() {
		dateStart, dateEnd, _ := GetReportDates(p.Data.Start, p.Data.End)
		if r.Data.Start != dateStart || r.Data.End != dateEnd {
			r.Data.Start = dateStart
			r.Data.End = dateEnd
			areChanges = true
		}
		r.Data.YearMonths = calculateYearMonths(PD_BI_PZFRR, r.Data.Start)
	}

	t := GetPzrrReportTemplate(r.Data)

	forecastedCapacityUp := getPayloads(t.ForecastedCapacityUp, p.ForecastedCapacityUp)
	if !payloadsAreEqual(r.ForecastedCapacityUp, forecastedCapacityUp) {
		areChanges = true
	}
	r.ForecastedCapacityUp = forecastedCapacityUp

	forecastedCapacityDown := getPayloads(t.ForecastedCapacityDown, p.ForecastedCapacityDown)
	if !payloadsAreEqual(r.ForecastedCapacityDown, forecastedCapacityDown) {
		areChanges = true
	}
	r.ForecastedCapacityDown = forecastedCapacityDown

	return
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
