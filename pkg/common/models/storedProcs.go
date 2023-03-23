package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type reportType int

const (
	kjcz reportType = iota
	pzrr
	pzfrr
)

func (rt reportType) String() string {
	return []string{"kjcz", "pzrr", "pzfrr"}[rt]
}

func GetPutKjczReportBody(data ReportData) string {
	return GetPutReportBody(data, kjcz)
}

func GetPutPzrrReportBody(data ReportData) string {
	return GetPutReportBody(data, pzrr)
}

func GetPutPzfrrReportBody(data ReportData) string {
	return GetPutReportBody(data, pzfrr)
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

func GetPutReportBody(data ReportData, rt reportType) string {
	rdata := fmt.Sprintf(":1 := hl_entsoe_reports_pk.put_%s_report("+
		"p_creator => '%s', "+
		"p_report_start => date '%s', "+
		"p_report_end   => date '%s');", rt.String(), data.Creator, data.Start.Format(time.DateOnly), data.End.Format(time.DateOnly))

	return strings.Join([]string{"begin", rdata, "end;"}, " ")
}

func getSecondaryQuantityString(secondaryQuantity *int) string {
	if secondaryQuantity == nil {
		return "null"
	}
	return strconv.Itoa(*secondaryQuantity)
}
