package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ReportType int

const (
	Kjcz ReportType = iota
	Pzrr
	Pzfrr
)

func (rt ReportType) String() string {
	return []string{"kjcz", "pzrr", "pzfrr"}[rt]
}

type Ekstrakt int

const (
	PD_BI_PZFRR Ekstrakt = iota
	PD_BI_PZRR
	PR_SO_KJCZ
)

func (e Ekstrakt) String() string {
	return []string{"PD_BI_PZFRR", "PD_BI_PZRR", "PR_SO_KJCZ"}[e]
}

type Resolution int

const (
	P1Y Resolution = iota
	P1M
	P1D
	PT60M
	PT30M
	PT15M
	PT1M
)

func (r Resolution) String() string {
	return []string{"P1Y", "P1M", "P1D", "PT60M", "PT30M", "PT15M", "PT1M"}[r]
}

func GetPutKjczReportBody(data ReportData) string {
	return GetPutReportBody(data, Kjcz)
}

func GetPutPzrrReportBody(data ReportData) string {
	return GetPutReportBody(data, Pzrr)
}

func GetPutPzfrrReportBody(data ReportData) string {
	return GetPutReportBody(data, Pzfrr)
}

func GetAddPayloadEntryBody(payload ReportPayload) string {
	rdata := fmt.Sprintf("hl_entsoe_reports_pk.add_payload_entry("+
		"p_reportid => %d, "+
		"p_mrid => %d, "+
		"p_businesstype => '%s', "+
		"p_flowdirection => '%s', "+
		"p_quantitymeasureunit => '%s', "+
		"p_position => %d, "+
		"p_quantity => %.3f, "+
		"p_secondaryquantity => %s);", payload.ReportId, payload.MrId, payload.BusinessType, payload.FlowDirection,
		payload.QuantityMeasureUnit, payload.Position, payload.Quantity, getSecondaryQuantityString(payload.SecondaryQuantity))

	s := strings.Join([]string{"begin", rdata, "end;"}, " ")
	return s
}

func getSecondaryQuantityString(secondaryQuantity *int) string {
	if secondaryQuantity == nil {
		return "null"
	}
	return strconv.Itoa(*secondaryQuantity)
}

func GetPutReportBody(data ReportData, rt ReportType) string {
	rdata := fmt.Sprintf(":1 := hl_entsoe_reports_pk.put_%s_report("+
		"p_creator => '%s', "+
		"p_report_start => date '%s', "+
		"p_report_end   => date '%s');", rt.String(), data.Creator, data.Start.Format(time.DateOnly), data.End.Format(time.DateOnly))

	return strings.Join([]string{"begin", rdata, "end;"}, " ")
}

func GetInicjujPozyskanie(rt ReportType, rd ReportData) string {
	rdata := fmt.Sprintf("inicjujPozyskanie("+
		"p_ekstrakt => '%s', "+
		"p_zakresOd => null, "+
		"p_zakresDo => null, "+
		"p_ziarno => '%s', "+
		"p_dataOd => ad_czas.podajCzasUTC(to_date('%s','yyyy-mm-dd'),'N'), "+
		"p_dataDo => ad_czas.podajCzasUTC(to_date('%s','yyyy-mm-dd'),'N'), "+
		"p_zrodlo => null, "+
		"p_obiekt_danych => null);", PR_SO_KJCZ.String(), P1M.String(), rd.Start.Format(time.DateOnly), rd.End.Format(time.DateOnly))

	s := strings.Join([]string{"begin", rdata, "end;"}, " ")
	return s
}
