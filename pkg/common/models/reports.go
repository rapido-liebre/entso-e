package models

import "time"

type Reporter interface {
	GetAllPayloads() []ReportPayload
}

type ReportData struct {
	Creator string
	Start   time.Time
	End     time.Time
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

func (r KjczReport) GetAllPayloads() (payloads []ReportPayload) {
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

type PzrrReport struct {
	Data                   ReportData
	ForecastedCapacityUp   []ReportPayload
	ForecastedCapacityDown []ReportPayload
}

func (r PzrrReport) GetAllPayloads() (payloads []ReportPayload) {
	payloads = append(payloads, r.ForecastedCapacityUp...)
	payloads = append(payloads, r.ForecastedCapacityDown...)

	return
}

type PzfrrReport struct {
	Data                   ReportData
	ForecastedCapacityUp   []ReportPayload
	ForecastedCapacityDown []ReportPayload
}

func (r PzfrrReport) GetAllPayloads() (payloads []ReportPayload) {
	payloads = append(payloads, r.ForecastedCapacityUp...)
	payloads = append(payloads, r.ForecastedCapacityDown...)

	return
}

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
