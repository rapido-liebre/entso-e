package config

import "entso-e_reports/pkg/common/models"

// DBAction defines action performed on database
type DBAction struct {
	Publish        bool              `json:"publish"`         // is report to be published after send to db
	TestData       bool              `json:"test_data"`       // send to db testing report, depends of report type
	ConnectionOnly bool              `json:"connection_only"` // is only test for connection
	ReportType     models.ReportType `json:"report_type"`     // report type
	Payload        string            `json:"payload"`         // payload containing full data report (data + payload)
}
