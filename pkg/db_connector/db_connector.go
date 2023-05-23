package db_connector

import (
	"database/sql"
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"errors"
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
				log.Printf("Error occured, err: %v\n", err)
				dbc.releaseChannel(err)
			}
			//
			log.Printf("Connect to DB successful isRunning:%v  status:%v\n", dbc.isRunning, dbc.status)
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
		dbc.getTestReport()
		return
	case "linux":
		fallthrough
	default:
		goto linux
	}

linux:
	defer func() {
		if dbc.db != nil {
			err := dbc.db.Close()
			if err != nil {
				fmt.Println("Can't close connection: ", err)
			}
		}
	}()

	t := time.Now()
	if dbc.data.ConnectionOnly {
		if dbc.data.ReportType == models.FETCH_15_MIN {
			if err := dbc.callFetchLfcAce(); err != nil {
				dbc.errch <- err
			}
		} else {
			if err := dbc.connectToDB(); err != nil {
				dbc.errch <- err
			}
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
		goto end
	}
	//normal usage
	if dbc.data.Payload != nil {
		if err := dbc.callSaveReport(); err != nil {
			dbc.errch <- err
		}
	} else {
		if err := dbc.callGetReport(); err != nil { //dbc.getTestReport()
			dbc.errch <- err
		}
	}

	//if len(p.data) == 0 {
	//	p.errch <- errors.New("no data for processing")
	//}
end:
	fmt.Println("Time Elapsed", time.Now().Sub(t).Milliseconds())
	dbc.status = Ready
	dbc.errch <- nil
}

func (dbc *dbConnector) releaseChannel(err error) {
	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		dbc.channels.KjczReport <- models.KjczReport{
			Data: models.ReportData{
				Error: err,
			}}
	case models.PD_BI_PZRR:
		dbc.channels.PzrrReport <- models.PzrrReport{
			Data: models.ReportData{
				Error: err,
			}}
	case models.PD_BI_PZFRR:
		dbc.channels.PzfrrReport <- models.PzfrrReport{
			Data: models.ReportData{
				Error: err,
			}}
	}
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

func (dbc *dbConnector) callFetchLfcAce() error {
	t := time.Now()

	if err := dbc.connectToDB(); err != nil {
		return err
	}
	dbc.status = Ready

	//fetch 15min data
	lfcAce15min, err := dbc.fetchRawLfcAce(models.FETCH_15_MIN)
	if err != nil {
		return err
	}

	//fetch 1min data
	lfcAce1min, err := dbc.fetchRawLfcAce(models.FETCH_1_MIN)
	if err != nil {
		return err
	}

	var rc models.ReportCalculator
	report := rc.Calculate(lfcAce15min, lfcAce1min)
	//fmt.Println(report)

	dbc.channels.KjczReport <- report

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))
	return nil
}

func (dbc *dbConnector) fetchRawLfcAce(rt models.ReportType) ([]models.LfcAce, error) {
	if rt < models.FETCH_15_MIN {
		return []models.LfcAce{}, errors.New(fmt.Sprintf("Wrong report type! Expected: %s or %s, got: %s",
			models.FETCH_15_MIN.String(), models.FETCH_1_MIN.String(), rt.String()))
	}
	t := time.Now()

	statement := models.GetFetchSourceData(dbc.data.ReportData, rt)
	fmt.Println(statement)

	// fetching multiple rows
	dataRows, err := dbc.db.Query(statement)
	if err != nil {
		return []models.LfcAce{}, err
	}
	defer dataRows.Close()

	var lfcAce []models.LfcAce

	for dataRows.Next() {
		var lfc models.LfcAce
		err = dataRows.Scan(&lfc.AvgTime, &lfc.SaveTime, &lfc.AvgName, &lfc.AvgValue, &lfc.AvgStatus, &lfc.SystemSite)
		if err != nil {
			return []models.LfcAce{}, err
		}
		lfcAce = append(lfcAce, lfc)
	}
	fmt.Printf("len(%s): %d\n", rt.Shortly(), len(lfcAce))
	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))

	return lfcAce, nil
}

func (dbc *dbConnector) testData() error {
	// fetch data from db only for kjcz
	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		if err := dbc.callFetchLfcAce(); err != nil {
			return err
		}
	case models.PD_BI_PZRR:
		fallthrough
	case models.PD_BI_PZFRR:
		if err := dbc.getTestReport(); err != nil {
			return err
		}
	}
	return nil
}

func (dbc *dbConnector) testDataAndPublish() error {
	if err := dbc.connectToDB(); err != nil {
		return err
	}
	rdata, err := dbc.callPutTestReport()
	if err != nil {
		return err
	}
	return dbc.callInicjujPozyskanie(rdata)
}

func (dbc *dbConnector) callPutTestReport() (models.ReportData, error) {
	t := time.Now()

	data := models.TestReportData(dbc.data.ReportType, dbc.data.ReportData.Start)

	var reportId int64
	statement := models.GetPutReportBody(data, dbc.data.ReportType)
	if _, err := dbc.db.Exec(statement, sql.Out{Dest: &reportId}); err != nil {
		return data, err
	}

	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		report := models.GetTestKjczReportBody(reportId, data)
		for _, payload := range report.GetAllPayloads() {
			statement = models.GetAddPayloadEntryBody(payload)
			_, err := dbc.db.Exec(statement)
			if err != nil {
				return data, err
			}
		}
	case models.PD_BI_PZRR:
		report := models.GetTestPzrrReportBody(reportId, data)
		for _, payload := range report.GetAllPayloads() {
			statement = models.GetAddPayloadEntryBody(payload)
			_, err := dbc.db.Exec(statement)
			if err != nil {
				return data, err
			}
		}
	case models.PD_BI_PZFRR:
		report := models.GetTestPzfrrReportBody(reportId, data)
		for _, payload := range report.GetAllPayloads() {
			statement = models.GetAddPayloadEntryBody(payload)
			_, err := dbc.db.Exec(statement)
			if err != nil {
				return data, err
			}
		}
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))

	return data, nil
}

func (dbc *dbConnector) callPutReport(report any) error {
	t := time.Now()

	var reportId int64

	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		r := report.(models.KjczReport)
		statement := models.GetPutReportBody(r.Data, dbc.data.ReportType)
		if _, err := dbc.db.Exec(statement, sql.Out{Dest: &reportId}); err != nil {
			return err
		}
		for _, payload := range r.GetAllPayloads() {
			payload.ReportId = reportId
			statement = models.GetAddPayloadEntryBody2(payload)
			if _, err := dbc.db.Exec(statement); err != nil {
				return err
			}
		}
	case models.PD_BI_PZRR:
		r := report.(models.PzrrReport)
		statement := models.GetPutReportBody(r.Data, dbc.data.ReportType)
		if _, err := dbc.db.Exec(statement, sql.Out{Dest: &reportId}); err != nil {
			return err
		}
		for _, payload := range r.GetAllPayloads() {
			payload.ReportId = reportId
			statement = models.GetAddPayloadEntryBody(payload)
			if _, err := dbc.db.Exec(statement); err != nil {
				return err
			}
		}
	case models.PD_BI_PZFRR:
		r := report.(models.PzfrrReport)
		statement := models.GetPutReportBody(r.Data, dbc.data.ReportType)
		if _, err := dbc.db.Exec(statement, sql.Out{Dest: &reportId}); err != nil {
			return err
		}
		for _, payload := range r.GetAllPayloads() {
			payload.ReportId = reportId
			statement = models.GetAddPayloadEntryBody(payload)
			if _, err := dbc.db.Exec(statement); err != nil {
				return err
			}
		}
	}

	if dbc.data.Publish {
		statement := models.GetSetReported(reportId)
		if _, err := dbc.db.Exec(statement); err != nil {
			return err
		}
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))

	return nil
}

func (dbc *dbConnector) callInicjujPozyskanie(rdata models.ReportData) error {
	t := time.Now()

	statement := models.GetInicjujPozyskanie(rdata, dbc.data.ReportType)
	fmt.Println(statement)
	_, err := dbc.db.Exec(statement)
	if err != nil {
		return err
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))

	return nil
}

func (dbc *dbConnector) getTestReport() error {
	data := models.TestReportData(dbc.data.ReportType, dbc.data.ReportData.Start)
	//data.Start = dbc.data.ReportData.Start
	//data.End = dbc.data.ReportData.End
	if len(dbc.data.ReportData.YearMonths) > 0 {
		data.YearMonths = dbc.data.ReportData.YearMonths
	}

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
	return nil
}

func (dbc *dbConnector) callSaveReport() error {
	t := time.Now()

	if err := dbc.connectToDB(); err != nil {
		return err
	}
	dbc.status = Ready

	var (
		cursorReport  go_ora.RefCursor
		cursorPayload go_ora.RefCursor
	)

	//get last report
	statement := models.GetLastReport(dbc.data.ReportData, dbc.data.ReportType)
	if _, err := dbc.db.Exec(statement, sql.Out{Dest: &cursorReport}, sql.Out{Dest: &cursorPayload}); err != nil {
		return err
	}
	defer cursorReport.Close()
	defer cursorPayload.Close()

	var (
		cd       models.CursorData
		cps      []models.CursorPayload
		reportId int64
	)

	//fetch report data
	dataRows, err := cursorReport.Query()
	if err != nil {
		return err
	}
	for dataRows.Next_() {
		err = dataRows.Scan(&reportId, &cd.ReportType, &cd.Revision, &cd.Creator, &cd.Created, &cd.Start, &cd.End, &cd.Saved, &cd.Reported)
		if err != nil {
			return err
		}
	}

	//fetch report payload
	payloadRows, err := cursorPayload.Query()
	if err != nil {
		return err
	}
	for payloadRows.Next_() {
		var cp models.CursorPayload
		err = payloadRows.Scan(&cp.MrId, &cp.BusinessType, &cp.FlowDirection, &cp.QuantityMeasurement, &cp.Position, &cp.Quantity, &cp.SecondaryQuantity)
		if err != nil {
			return err
		}
		cps = append(cps, cp)
	}

	//save to report
	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		report := models.KjczReport{}
		if cd.IsValid() {
			report.SaveCursors(cd, cps)
		}
		if areChanges := report.Update(dbc.data.Payload); areChanges == true {
			if err = dbc.callPutReport(report); err != nil {
				return err
			}
			//sync revision to corresponding in DB
			report.Data.Revision += 1
			report.Data.Saved = time.Now()
			if report.Data.Created.IsZero() {
				report.Data.Created = report.Data.Saved
			}
			report.Data.Reported = time.Time{}
			//reset reportId to avoid publish old revision of the report
			reportId = 0
		}
		if dbc.data.Publish {
			if reportId > 0 {
				statement = models.GetSetReported(reportId)
				if _, err = dbc.db.Exec(statement); err != nil {
					return err
				}
			}
			if err = dbc.callInicjujPozyskanie(report.Data); err != nil {
				if dbc.config.Params.FakePublish {
					fmt.Println("Fake publish triggered")
				} else {
					return err
				}
			}
			report.Data.Reported = time.Now()
		}
		dbc.channels.KjczReport <- report
	case models.PD_BI_PZRR:
		report := models.PzrrReport{}
		if cd.IsValid() {
			report.SaveCursors(cd, cps)
		}
		if areChanges := report.Update(dbc.data.Payload); areChanges == true {
			if err = dbc.callPutReport(report); err != nil {
				return err
			}
			//sync revision to corresponding in DB
			report.Data.Revision += 1
			report.Data.Saved = time.Now()
			if report.Data.Created.IsZero() {
				report.Data.Created = report.Data.Saved
			}
			report.Data.Reported = time.Time{}
			//reset reportId to avoid publish old revision of the report
			reportId = 0
		}
		if dbc.data.Publish {
			if reportId > 0 {
				statement = models.GetSetReported(reportId)
				if _, err = dbc.db.Exec(statement); err != nil {
					return err
				}
			}
			if err = dbc.callInicjujPozyskanie(report.Data); err != nil {
				if dbc.config.Params.FakePublish {
					fmt.Println("Fake publish triggered")
				} else {
					return err
				}
			}
			report.Data.Reported = time.Now()
		}
		dbc.channels.PzrrReport <- report
	case models.PD_BI_PZFRR:
		report := models.PzfrrReport{}
		if cd.IsValid() {
			report.SaveCursors(cd, cps)
		}
		if areChanges := report.Update(dbc.data.Payload); areChanges == true {
			if err = dbc.callPutReport(report); err != nil {
				return err
			}
			//sync revision to corresponding in DB
			report.Data.Revision += 1
			report.Data.Saved = time.Now()
			if report.Data.Created.IsZero() {
				report.Data.Created = report.Data.Saved
			}
			report.Data.Reported = time.Time{}
			//reset reportId to avoid publish old revision of the report
			reportId = 0
		}
		if dbc.data.Publish {
			if reportId > 0 {
				statement = models.GetSetReported(reportId)
				if _, err = dbc.db.Exec(statement); err != nil {
					return err
				}
			}
			if err = dbc.callInicjujPozyskanie(report.Data); err != nil {
				if dbc.config.Params.FakePublish {
					fmt.Println("Fake publish triggered")
				} else {
					return err
				}
			}
			report.Data.Reported = time.Now()
		}
		dbc.channels.PzfrrReport <- report
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))
	return nil
}

func (dbc *dbConnector) callGetReport() error {
	t := time.Now()

	if err := dbc.connectToDB(); err != nil {
		return err
	}
	dbc.status = Ready

	var (
		cursorReport  go_ora.RefCursor
		cursorPayload go_ora.RefCursor
	)

	//get last report
	statement := models.GetLastReport(dbc.data.ReportData, dbc.data.ReportType)
	if _, err := dbc.db.Exec(statement, sql.Out{Dest: &cursorReport}, sql.Out{Dest: &cursorPayload}); err != nil {
		return err
	}
	defer cursorReport.Close()
	defer cursorPayload.Close()

	var (
		cd       models.CursorData
		cps      []models.CursorPayload
		reportId int64
	)

	//fetch report data
	dataRows, err := cursorReport.Query()
	if err != nil {
		return err
	}
	for dataRows.Next_() {
		err = dataRows.Scan(&reportId, &cd.ReportType, &cd.Revision, &cd.Creator, &cd.Created, &cd.Start, &cd.End, &cd.Saved, &cd.Reported)
		if err != nil {
			return err
		}
		//fmt.Println(cd)
	}

	//fetch report payload
	payloadRows, err := cursorPayload.Query()
	if err != nil {
		return err
	}
	for payloadRows.Next_() {
		var cp models.CursorPayload
		err = payloadRows.Scan(&cp.MrId, &cp.BusinessType, &cp.FlowDirection, &cp.QuantityMeasurement, &cp.Position, &cp.Quantity, &cp.SecondaryQuantity)
		if err != nil {
			return err
		}
		//fmt.Println(cp)
		cps = append(cps, cp)
	}

	//save to report
	switch dbc.data.ReportType {
	case models.PR_SO_KJCZ:
		report := models.KjczReport{}
		if cd.IsValid() {
			report.SaveCursors(cd, cps)
		}
		dbc.channels.KjczReport <- report
	case models.PD_BI_PZRR:
		report := models.PzrrReport{}
		if cd.IsValid() {
			report.SaveCursors(cd, cps)
		}
		dbc.channels.PzrrReport <- report
	case models.PD_BI_PZFRR:
		report := models.PzfrrReport{}
		if cd.IsValid() {
			report.SaveCursors(cd, cps)
		}
		dbc.channels.PzfrrReport <- report
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))
	return nil
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
