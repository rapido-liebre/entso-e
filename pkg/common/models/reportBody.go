package models

type BodyData struct {
	Creator string `json:"creator"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

type BodyReportPayload struct {
	Position  int     `json:"position"`
	Quantity  float64 `json:"quantity"`
	YearMonth string  `json:"year_month"`
}

type KjczBody struct {
	Data                                BodyData            `json:"data"`
	MeanValue                           []BodyReportPayload `json:"meanValue"`
	StandardDeviation                   []BodyReportPayload `json:"standardDeviation"`
	Percentile1                         []BodyReportPayload `json:"percentile1"`
	Percentile5                         []BodyReportPayload `json:"percentile5"`
	Percentile10                        []BodyReportPayload `json:"percentile10"`
	Percentile90                        []BodyReportPayload `json:"percentile90"`
	Percentile95                        []BodyReportPayload `json:"percentile95"`
	Percentile99                        []BodyReportPayload `json:"percentile99"`
	FrceOutsideLevel1RangeUp            []BodyReportPayload `json:"frceOutsideLevel1RangeUp"`
	FrceOutsideLevel1RangeDown          []BodyReportPayload `json:"frceOutsideLevel1RangeDown"`
	FrceOutsideLevel2RangeUp            []BodyReportPayload `json:"frceOutsideLevel2RangeUp"`
	FrceOutsideLevel2RangeDown          []BodyReportPayload `json:"frceOutsideLevel2RangeDown"`
	FrceExceeded60PercOfFRRCapacityUp   []BodyReportPayload `json:"frceExceeded60PercOfFRRCapacityUp"`
	FrceExceeded60PercOfFRRCapacityDown []BodyReportPayload `json:"frceExceeded60PercOfFRRCapacityDown"`
}

type PzrrBody struct {
	Data                   BodyData            `json:"data"`
	ForecastedCapacityUp   []BodyReportPayload `json:"forecastedCapacityUp"`
	ForecastedCapacityDown []BodyReportPayload `json:"forecastedCapacityDown"`
}

type PzfrrBody struct {
	Data                   BodyData            `json:"data"`
	ForecastedCapacityUp   []BodyReportPayload `json:"forecastedCapacityUp"`
	ForecastedCapacityDown []BodyReportPayload `json:"forecastedCapacityDown"`
}
