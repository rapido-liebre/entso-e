package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//type storedProc struct {
//	name string
//	in   []string
//	out  []string
//}

type ReportData struct {
	Creator string
	Start   time.Time //time.Parse(time.DateOnly, val)
	End     time.Time
}

type ReportPayload struct {
	ReportId            int64
	MrId                int
	BusinessType        string
	FlowDirection       string
	QuantityMeasureunit string
	Position            int
	Quantity            float64
	SecondaryQuantity   *int
}

func GetPutKjczReportBody(data ReportData) string {
	rdata := fmt.Sprintf(":1 := hl_entsoe_reports_pk.put_kjcz_report("+
		"p_creator => '%s', "+
		"p_report_start => date '%s', "+
		"p_report_end   => date '%s');", data.Creator, data.Start.Format(time.DateOnly), data.End.Format(time.DateOnly))

	return strings.Join([]string{"begin", rdata, "end;"}, " ")
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
		payload.QuantityMeasureunit, payload.Position, payload.Quantity, getSecondaryQuantityString(payload.SecondaryQuantity))

	s := strings.Join([]string{"begin", rdata, "end;"}, " ")
	return s
}

/*
function put_kjcz_report(p_creator      in hl_entsoe_reports.creator%type,
p_report_start in hl_entsoe_reports.report_start%type,
p_report_end   in hl_entsoe_reports.report_end%type)
return hl_entsoe_reports.id%type;


procedure add_payload_entry(
p_reportid            in hl_entsoe_report_payloads.reportid%type,
p_mrid                in hl_entsoe_report_payloads.mrid%type,
p_businesstype        in hl_entsoe_report_payloads.businesstype%type,
p_flowdirection       in hl_entsoe_report_payloads.flowdirection%type,
p_quantitymeasureunit in hl_entsoe_report_payloads.quantitymeasureunit%type,
p_position            in hl_entsoe_report_payloads.position%type,
p_quantity            in hl_entsoe_report_payloads.quantity%type,
p_secondaryquantity   in hl_entsoe_report_payloads.secondaryquantity%type);



begin

l_rep_id := hl_entsoe_reports_pk.put_kjcz_report(p_creator      => 'MARCIN',
p_report_start => date
'2023-03-01',
p_report_end   => date
'2023-05-01');



hl_entsoe_reports_pk.add_payload_entry(
p_reportid            => l_rep_id,
p_mrid                => 1,
p_businesstype        => 'A01',
p_flowdirection       => 'G',
p_quantitymeasureunit => 'MAW',
p_position            => 1,
p_quantity            => 10,
p_secondaryquantity   => null);



hl_entsoe_reports_pk.add_payload_entry(p_reportid            => l_rep_id,

p_mrid                => 1,

p_businesstype        => 'A01',

p_flowdirection       => 'G',

p_quantitymeasureunit => 'MAW',

p_position            => 2,

p_quantity            => 12.5,

p_secondaryquantity   => null);

hl_entsoe_reports_pk.add_payload_entry(p_reportid            => l_rep_id,

p_mrid                => 1,

p_businesstype        => 'A01',

p_flowdirection       => 'G',

p_quantitymeasureunit => 'MAW',

p_position            => 3,

p_quantity            => 15,

p_secondaryquantity   => null);

end;

*/

func getSecondaryQuantityString(secondaryQuantity *int) string {
	if secondaryQuantity == nil {
		return "null"
	}
	return strconv.Itoa(*secondaryQuantity)
}
