package db_connector

import (
	"database/sql"
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"fmt"
	go_ora "github.com/sijms/go-ora/v2"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

type DBConnector interface {
	Run(wg *sync.WaitGroup)
	connect()
}

type Status int

const (
	Ready Status = iota
	Processing
)

type dbConnector struct {
	isRunning bool
	status    Status
	config    config.Config
	channels  *config.Channels
	errch     chan error
	data      config.DBAction
	db        *sql.DB
}

// NewService returns new DBConnector instance
func NewService(cfg config.Config, ch *config.Channels) DBConnector {
	return &dbConnector{
		config:   cfg,
		channels: ch,
		errch:    make(chan error, 1),
		status:   Ready,
	}
}

func (dbc *dbConnector) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	// proceed in infinite loop
	for {
		select {
		//case filename := <-p.channels.Filename:
		//	//using map here instead of slice for easier lookup during processing
		//	if len(filename) > 0 {
		//		p.filenames[filename] = true
		//		log.Printf("Queued: %s", filename)
		//	}
		case dbc.data = <-dbc.channels.RunDBConn:
			if dbc.isRunning { //TODO check if dbConnector is ready
				if dbc.status == Ready {
					dbc.status = Processing
					go dbc.connect()
				}
			}
		case dbc.isRunning = <-dbc.channels.DBConnectorIsRunning:
			log.Printf("DBConnector is running: %v\n", dbc.isRunning)

		case err := <-dbc.errch:
			if err != nil {
				log.Printf("Database connection failed, err: %v\n", err)
			}
			//
			log.Printf("Connect to DB successful isRunning:%v  status:%v", dbc.isRunning, dbc.status)
			//if dbc.isRunning && dbc.status == Ready {
			//	dbc.channels.RunProcess <- dbc.data
			//}
		//case <-dbc.channels.CfgUpdate:
		//	// TODO config update
		//	log.Println("Processor updates config")
		case dbc.isRunning = <-dbc.channels.Quit:
			// TODO should wait until dbConnector completes its job
			log.Printf("DBConnector says Bye bye.. status:%v", dbc.status)
			return
		}
	}
}

func (dbc *dbConnector) connect() {
	fmt.Println("*** Using only go_ora package (no additional client software)")
	fmt.Println("Local Database, simple connect string ")

	os := runtime.GOOS
	switch os {
	case "windows":
		fmt.Println("This service is not dedicated to running on Windows")
	case "darwin":
		fmt.Println("This service can be run on Mac but then connection to oracle db is skipped")
		dbc.getReport()
		return
	case "linux":
		fallthrough
	default:
		goto linux
	}

linux:
	defer func() {
		err := dbc.db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()

	t := time.Now()
	if dbc.data.ConnectionOnly {
		err := dbc.connectToDB()
		if err != nil {
			dbc.errch <- err
		}
		goto end
	}
	if dbc.data.TestData {
		if dbc.data.Publish {
			err := dbc.testDataAndPublish()
			if err != nil {
				dbc.errch <- err
			}
			goto end
		}
		err := dbc.testData()
		if err != nil {
			dbc.errch <- err
		}
	} else {
		//normal usage
		dbc.getReport()
	}

	//if len(p.data) == 0 {
	//	p.errch <- errors.New("no data for processing")
	//}
end:
	fmt.Println("Time Elapsed", time.Now().Sub(t).Milliseconds())
	dbc.status = Ready
	dbc.errch <- nil
}

func (dbc *dbConnector) connectToDB() error {
	cfg := dbc.config.Params
	//connectionString := "oracle://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBServer + ":" + cfg.DBPort + "/" + cfg.DBService
	//if cfg.DBDSN != "" {
	//	connectionString += "?TRACE FILE=trace.log&SSL=enable&SSL Verify=false&WALLET=" + cfg.DBDSN //url.QueryEscape(dbParams["walletLocation"])
	//}

	port, _ := strconv.Atoi(cfg.DBPort)

	urlOptions := map[string]string{
		"TRACE FILE": "trace.log",
		"AUTH TYPE":  "TCPS",
		"SSL":        "TRUE",
		"SSL VERIFY": "FALSE",
		"WALLET":     cfg.DBWallet,
	}
	connectionString := go_ora.BuildUrl(cfg.DBServer, port, cfg.DBService, "", "", urlOptions)

	if len(cfg.ConnString) > 0 {
		fmt.Println("Using provided connection string")
		connectionString = cfg.ConnString
	}

	//"oracle://10.69.9.32:1522/OSP&&wallet=/usr/lib/oracle/18.3/client64/network/wallet"
	//"oracle://10.69.9.32:1522:OSP&&wallet=/usr/lib/oracle/18.3/client64/network/wallet"

	fmt.Println(connectionString)
	var err error
	dbc.db, err = sql.Open("oracle", connectionString)
	if err != nil {
		return fmt.Errorf("error in sql.Open: %w", err)
	}
	//defer func() {
	//	err = db.Close()
	//	if err != nil {
	//		fmt.Println("Can't close connection: ", err)
	//	}
	//}()
	dbc.db.SetConnMaxLifetime(time.Minute * 5)
	err = dbc.db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging db: %w", err)
	}
	return nil
}

//const createTableStatement = "CREATE TABLE TEMP_TABLE ( NAME VARCHAR2(100), CREATION_TIME TIMESTAMP DEFAULT SYSTIMESTAMP, VALUE  NUMBER(5))"
//const dropTableStatement = "DROP TABLE TEMP_TABLE PURGE"
//const insertStatement = "INSERT INTO TEMP_TABLE ( NAME , VALUE) VALUES (:name, :value)"

func (dbc *dbConnector) testData() error {
	if err := dbc.connectToDB(); err != nil {
		return err
	}
	if _, err := dbc.callPutReport(); err != nil {
		return err
	}

	return nil
}

func (dbc *dbConnector) testDataAndPublish() error {
	if err := dbc.connectToDB(); err != nil {
		return err
	}
	rdata, err := dbc.callPutReport()
	if err != nil {
		return err
	}
	return dbc.callInicjujPozyskanie(rdata)
}

func (dbc *dbConnector) callPutReport() (models.ReportData, error) {
	t := time.Now()

	data := models.TestReportData(dbc.data.ReportType)

	var reportId int64
	_, err := dbc.db.Exec(models.GetPutReportBody(data, dbc.data.ReportType), sql.Out{Dest: &reportId})

	if err != nil {
		return data, err
	}

	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		report := models.GetTestKjczReportBody(reportId, data)
		for _, payload := range report.GetAllPayloads() {
			_, err := dbc.db.Exec(models.GetAddPayloadEntryBody(payload))
			if err != nil {
				return data, err
			}
		}
	case models.PD_BI_PZRR:
		report := models.GetTestPzrrReportBody(reportId, data)
		for _, payload := range report.GetAllPayloads() {
			_, err := dbc.db.Exec(models.GetAddPayloadEntryBody(payload))
			if err != nil {
				return data, err
			}
		}
	case models.PD_BI_PZFRR:
		report := models.GetTestPzfrrReportBody(reportId, data)
		for _, payload := range report.GetAllPayloads() {
			_, err := dbc.db.Exec(models.GetAddPayloadEntryBody(payload))
			if err != nil {
				return data, err
			}
		}
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))

	return data, nil
}

func (dbc *dbConnector) callInicjujPozyskanie(rdata models.ReportData) error {
	t := time.Now()

	_, err := dbc.db.Exec(models.GetInicjujPozyskanie(dbc.data.ReportType, rdata))
	if err != nil {
		return err
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))

	return nil
}

func (dbc *dbConnector) getReport() {
	data := models.TestReportData(dbc.data.ReportType)
	//data.Start = dbc.data.ReportData.Start
	//data.End = dbc.data.ReportData.End
	data.MonthsDuration = dbc.data.ReportData.MonthsDuration

	var reportId int64
	reportId = 0

	dbc.status = Ready

	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		report := models.GetTestKjczReportBody(reportId, data)
		dbc.channels.KjczReport <- report
	case models.PD_BI_PZRR:
		report := models.GetTestPzrrReportBody(reportId, data)
		dbc.channels.PzrrReport <- report
	case models.PD_BI_PZFRR:
		report := models.GetTestPzfrrReportBody(reportId, data)
		dbc.channels.PzfrrReport <- report
	default:
		fmt.Println("getReport() fatal error! Unknown report type")
	}
}

//func someAdditionalActions(_ *sql.DB) {

//var queryResultColumnOne string
//row := db.QueryRow("SELECT systimestamp FROM dual")
//err := row.Scan(&queryResultColumnOne)
//if err != nil {
//	panic(fmt.Errorf("error scanning db: %w", err))
//}
//fmt.Println("The time in the database ", queryResultColumnOne)
//_, err = db.Exec(createTableStatement)
//handleError("create table", err)
//defer db.Exec(dropTableStatement)
//stmt, err := db.Prepare(insertStatement)
//handleError("prepare insert statement", err)
//sqlresult, err := stmt.Exec("John", 42)
//handleError("execute insert statement", err)
//rowCount, _ := sqlresult.RowsAffected()
//fmt.Println("Inserted number of rows = ", rowCount)
//
//var queryResultName string
//var queryResultTimestamp string
//var queryResultValue int32
//row = db.QueryRow("SELECT name, creation_time, value FROM temp_table")
//err = row.Scan(&queryResultName, &queryResultTimestamp, &queryResultValue)
//handleError("query single row", err)
//if err != nil {
//	panic(fmt.Errorf("error scanning db: %w", err))
//}
//fmt.Println(fmt.Sprintf("The name: %s, time: %s, value:%d ", queryResultName, queryResultTimestamp, queryResultValue))
//
//_, err = stmt.Exec("Jane", 69)
//handleError("execute insert statement", err)
//_, err = stmt.Exec("Malcolm", 13)
//handleError("execute insert statement", err)
//
//// fetching multiple rows
//theRows, err := db.Query("select name, value from TEMP_TABLE")
//handleError("Query for multiple rows", err)
//defer theRows.Close()
//var (
//	name  string
//	value int32
//)
//for theRows.Next() {
//	err := theRows.Scan(&name, &value)
//	handleError("next row in multiple rows", err)
//	fmt.Println(fmt.Sprintf("The name: %s and value:%d ", name, value))
//}
//err = theRows.Err()
//handleError("next row in multiple rows", err)

//_ = callPutKjczReport(db)
//}
