package config

import (
	"entso-e_reports/pkg/common/models"
	"time"
)

type APIChannels struct {
	ParserIsRunning      chan bool
	ProcessorIsRunning   chan bool
	DBConnectorIsRunning chan bool
	Quit                 chan bool
	CfgUpdate            chan Config
}

type DataChannels struct {
	RunParse chan bool
	//RunGenerate chan bool
	RunProcess chan map[models.Year]map[time.Month][]models.LfcAce
	RunDBConn  chan bool
	Filename   chan string
}

type Channels struct {
	APIChannels
	DataChannels
}

func GetChannels() Channels {
	return Channels{
		APIChannels: APIChannels{
			ParserIsRunning:      make(chan bool),
			ProcessorIsRunning:   make(chan bool),
			DBConnectorIsRunning: make(chan bool),
			Quit:                 make(chan bool, 2),
			CfgUpdate:            make(chan Config),
		},
		DataChannels: DataChannels{
			RunParse: make(chan bool),
			//RunGenerate: make(chan bool),
			RunProcess: make(chan map[models.Year]map[time.Month][]models.LfcAce),
			RunDBConn:  make(chan bool),
			Filename:   make(chan string),
		},
	}
}
