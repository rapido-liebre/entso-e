package db_connector

import (
	"database/sql"
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"fmt"
	"log"
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
	data      map[models.Year]map[time.Month][]models.LfcAce
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
		case <-dbc.channels.RunDBConn:
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
				log.Fatalf("Connect to DB failed, err: %v\n", err)
				return
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
	t := time.Now()
	dbc.connectToDB()
	fmt.Println("Time Elapsed", time.Now().Sub(t).Milliseconds())

	//if len(p.data) == 0 {
	//	p.errch <- errors.New("no data for processing")
	//}

	dbc.status = Ready
	dbc.errch <- nil
}

func (p *dbConnector) connectToDB() {
	cfg := p.config.Params
	connectionString := "oracle://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBServer + ":" + cfg.DBPort + "/" + cfg.DBService
	//if val, ok := dbParams["walletLocation"]; ok && val != "" {
	//	connectionString += "?TRACE FILE=trace.log&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(dbParams["walletLocation"])
	//}
	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}

	someAdditionalActions(db)
}

//const createTableStatement = "CREATE TABLE TEMP_TABLE ( NAME VARCHAR2(100), CREATION_TIME TIMESTAMP DEFAULT SYSTIMESTAMP, VALUE  NUMBER(5))"
//const dropTableStatement = "DROP TABLE TEMP_TABLE PURGE"
//const insertStatement = "INSERT INTO TEMP_TABLE ( NAME , VALUE) VALUES (:name, :value)"

func someAdditionalActions(db *sql.DB) {

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
	_ = callPutKjczReport(db)
}

func callPutKjczReport(db *sql.DB) error {
	t := time.Now()

	tStart, _ := time.Parse(time.DateOnly, "2023-01-01")
	tEnd, _ := time.Parse(time.DateOnly, "2023-03-01")

	rdata := models.ReportData{
		Creator: "Janko Muzykant",
		Start:   tStart,
		End:     tEnd,
	}
	var reportId int64
	_, err := db.Exec(models.GetPutKjczReportBody(rdata), sql.Out{Dest: &reportId})

	if err != nil {
		return err
	}

	var rpayloads []models.ReportPayload
	rpayloads = append(rpayloads, models.ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        "C66",
		FlowDirection:       "A03",
		QuantityMeasureunit: "MAW",
		Position:            1,
		Quantity:            3.309,
		SecondaryQuantity:   nil,
	})
	rpayloads = append(rpayloads, models.ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        "C66",
		FlowDirection:       "A03",
		QuantityMeasureunit: "MAW",
		Position:            2,
		Quantity:            1.388,
		SecondaryQuantity:   nil,
	})
	rpayloads = append(rpayloads, models.ReportPayload{
		ReportId:            reportId,
		MrId:                1,
		BusinessType:        "C66",
		FlowDirection:       "A03",
		QuantityMeasureunit: "MAW",
		Position:            2,
		Quantity:            1.941,
		SecondaryQuantity:   nil,
	})

	for _, payload := range rpayloads {
		_, err := db.Exec(models.GetAddPayloadEntryBody(payload))
		if err != nil {
			return err
		}
	}

	fmt.Println("Finish call store procedure: ", time.Now().Sub(t))

	return nil
}
