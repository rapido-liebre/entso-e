package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ReportType int

const (
	PD_BI_PZFRR ReportType = iota
	PD_BI_PZRR
	PR_SO_KJCZ
	FETCH_15_MIN
	FETCH_1_MIN
)

func (rt ReportType) String() string {
	return []string{"PD_BI_PZFRR", "PD_BI_PZRR", "PR_SO_KJCZ", "RC_AVG15m_LFC_ACE_PL", "RC_AVG1M_LFC_ACE_PL"}[rt]
}
func (rt ReportType) shortly() string {
	return []string{"pzfrr", "pzrr", "kjcz", "avg_15", "avg_1"}[rt]
}

type Resolution int

const (
	P1Y Resolution = iota
	P3M
	P1M
	P1D
	PT60M
	PT30M
	PT15M
	PT1M
)

func (r Resolution) String() string {
	return []string{"P1Y", "P3M", "P1M", "P1D", "PT60M", "PT30M", "PT15M", "PT1M"}[r]
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

func GetAddPayloadEntryBody2(payload ReportPayload) string {
	rdata := fmt.Sprintf("insert into hl_entsoe_report_payloads("+
		"reportid, mrid, businesstype, flowdirection, quantitymeasureunit, position, quantity, secondaryquantity) "+
		"values (%d, %d, '%s', '%s', '%s', %d, %.3f, %s);", payload.ReportId, payload.MrId, payload.BusinessType, payload.FlowDirection,
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

func GetPutReportBody(rd ReportData, rt ReportType) string {
	rdata := fmt.Sprintf(":1 := hl_entsoe_reports_pk.put_%s_report("+
		"p_creator => '%s', "+
		"p_report_start => date '%s', "+
		"p_report_end   => date '%s');", rt.shortly(), rd.Creator, rd.Start.Format(time.DateOnly), rd.End.Format(time.DateOnly))

	return strings.Join([]string{"begin", rdata, "end;"}, " ")
}

func GetLastReport(rd ReportData, rt ReportType) string {
	rdata := fmt.Sprintf("hl_entsoe_reports_pk.get_last_%s("+
		"date '%s', "+
		"date '%s', "+
		":1, "+
		":2);", rt.shortly(), rd.Start.Format(time.DateOnly), rd.End.Format(time.DateOnly))

	return strings.Join([]string{"begin", rdata, "end;"}, " ")
}

func GetSetReported(reportId int64) string {
	rdata := fmt.Sprintf("hl_entsoe_reports_pk.set_reported(p_id => %d);", reportId)

	return strings.Join([]string{"begin", rdata, "end;"}, " ")
}

func GetInicjujPozyskanie(rd ReportData, rt ReportType) string {
	rdata := fmt.Sprintf("CN_INT_STERUJ_POZYSKANIEM_PK.inicjujPozyskanie("+
		"v_ekstrakt, "+
		"null, "+
		"null, "+
		"'%s', "+
		"ad_czas.podajCzasUTC(to_date('%s','yyyy-mm-dd'),'N'), "+
		"ad_czas.podajCzasUTC(to_date('%s','yyyy-mm-dd'),'N'), "+
		"null, "+
		"v_dane);", getResolution(rt), rd.Start.Format(time.DateOnly), rd.End.Format(time.DateOnly))

	s := strings.Join([]string{"declare v_ekstrakt ad_util.t_tDane; v_dane ad_util.t_tDane; begin",
		fmt.Sprintf("v_dane(1) := null; v_ekstrakt(1) := '%s'; ", rt.String()), rdata, "end;"}, " ")
	return s
}

func getResolution(rt ReportType) string {
	if rt == PR_SO_KJCZ {
		return P1M.String()
	}
	return P3M.String()
}

func GetFetchSourceData(rd ReportData, rt ReportType) string {
	rdata := fmt.Sprintf("SELECT avg_time, save_time, avg_name, avg_value, avg_status, system_site "+
		"FROM %s WHERE avg_time >= to_date('%s','yyyy-mm-dd') AND avg_time < to_date('%s','yyyy-mm-dd') AND avg_name = '%s'", //ORDER BY avg_value
		rt.shortly(), rd.Start.Format(time.DateOnly), rd.End.Format(time.DateOnly), rt.String())

	return rdata //strings.Join([]string{"begin", rdata, "end;"}, " ")
}

func GetFetchSourceData1min(rd ReportData) string {
	rdata := fmt.Sprintf("SELECT avg_value, avg_time FROM avg_1 WHERE avg_time >= '%s' AND avg_time < '%s' "+
		"AND avg_name = 'RC_AVG1M_LFC_ACE_PL' ORDER BY avg_time;", rd.Start.Format(time.DateOnly), rd.End.Format(time.DateOnly))

	return strings.Join([]string{"begin", rdata, "end;"}, " ")
}

//SELECT avg_value FROM avg_1mon WHERE avg_time = '$bar-01 00:00:00' AND avg_name = 'RC_AVG1MON_LFC_Pp_max';
