package db_connector

import (
	"bufio"
	"entso-e_reports/pkg/common/models"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	Data []models.LfcAce
}

// Parse parses data dump having PG format
func (p *Parser) Parse(filePath string) {

	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	var rowsFetched int

	for scanner.Scan() {
		line := scanner.Text()
		//valid line with data starts with date like '| 2022-12-01 '
		if strings.Contains(line, "| 2") {
			//fmt.Println(line)
			p.parseLine(line)
			rowsFetched++
		}
		//if strings.Contains(line, "rows fetched") {
		//	fmt.Println(scanner.Text())
		//}
	}
	fmt.Println(rowsFetched)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (p *Parser) parseLine(line string) {
	var measure models.LfcAce
	var result []string

	//line = strings.Replace(line, "|", " ", -1)
	split := strings.Split(line, "|")
	for _, s := range split {
		strings.Trim(s, " ")
		if len(s) != 0 {
			//fmt.Println(s)
			result = append(result, strings.Trim(s, " "))
		}
	}
	//6 values expected
	if len(result) != 6 {
		log.Fatalf("parsing line failed. Line: %s", line)
	}

	for idx, val := range result {
		var err error
		switch models.LfcAceIndex(idx) {
		case models.AvgTime:
			measure.AvgTime, err = time.Parse(time.DateTime, val)
			if err != nil {
				log.Fatal("Could not parse avg time:", err)
			}
		case models.SaveTime:
			measure.SaveTime, err = time.Parse(time.DateTime, val)
			if err != nil {
				log.Fatal("Could not parse save time:", val, err)
			}
		case models.AvgName:
			measure.AvgName = val
		case models.AvgValue:
			measure.AvgValue, err = strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal("Could not parse float value:", val, err)
			}
		case models.AvgStatus:
			measure.AvgStatus, err = strconv.Atoi(val)
			if err != nil {
				log.Fatal("Could not parse status value:", val, err)
			}
		case models.SystemSite:
			if len(val) != 1 {
				log.Fatal("Could not parse byte value:", val, err)
			}
			measure.SystemSite = val
		}
	}
	p.Data = append(p.Data, measure)
}

// Parse2 parses data dump having PM format
func (p *Parser) Parse2(filePath string) {

	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	var rowsFetched int

	rt := models.FETCH_15_MIN
	if strings.Contains(filePath, "1m") {
		rt = models.FETCH_1_MIN
	}

	for scanner.Scan() {
		line := scanner.Text()
		//valid line with data looks like '2022-12-31 23:00:00;-159.312042;B'
		if strings.Contains(line, ":") {
			//fmt.Println(line)
			p.parseLine2(line, rt)
			rowsFetched++
		}
		//if strings.Contains(line, "rows fetched") {
		//	fmt.Println(scanner.Text())
		//}
	}
	fmt.Println(rowsFetched)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (p *Parser) parseLine2(line string, rt models.ReportType) {
	var measure models.LfcAce
	var result []string

	//line = strings.Replace(line, "|", " ", -1)
	split := strings.Split(line, ";")
	for _, s := range split {
		strings.Trim(s, " ")
		if len(s) != 0 {
			//fmt.Println(s)
			result = append(result, strings.Trim(s, " "))
		}
	}
	//6 values expected
	if len(result) != 3 {
		log.Fatalf("parsing line failed. Line: %s", line)
	}

	for idx, val := range result {
		var err error
		switch models.LfcAceIndex(idx) {
		case 0:
			measure.AvgTime, err = time.Parse(time.DateTime, val)
			if err != nil {
				log.Fatal("Could not parse avg time:", err)
			}
		case 1:
			measure.AvgValue, err = strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal("Could not parse float value:", val, err)
			}
		case 2:
			//if len(val) != 1 {
			//	log.Fatal("Could not parse byte value:", val, err)
			//}
			measure.SystemSite = "B"
		}
	}
	measure.SaveTime = measure.AvgTime
	measure.AvgName = rt.String()
	measure.AvgStatus = 0
	p.Data = append(p.Data, measure)
}
