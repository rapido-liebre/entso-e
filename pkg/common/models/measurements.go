package models

import (
	"time"
)

// this model can be used when the configuration is sent in the request body

//type Config struct {
//	TimeInterval int    `json:"time_interval"`  // Time interval [minutes] between processing of successive archives
//	WarningSize  int    `json:"warning_size"`   // Minimum amount of free disk space [GB] before sending the alert
//	RedAlertSize int    `json:"red_alert_size"` // Minimum amount of free disk space [GB] before running the emergency plan
//	InputDir     string `json:"input_dir"`      // Input data directory for archiving
//	OutputDir    string `json:"output_dir"`     // Directory for holding output archive data
//	Port         string `json:"port"`           // The localhost port on which HTTP requests are listened
//}

type Year int

type LfcAce struct {
	AvgTime    time.Time
	SaveTime   time.Time
	AvgName    string
	AvgValue   float64
	AvgStatus  int
	SystemSite string
}

type LfcAceIndex int

const (
	AvgTime LfcAceIndex = iota
	SaveTime
	AvgName
	AvgValue
	AvgStatus
	SystemSite
)
