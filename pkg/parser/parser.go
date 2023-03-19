package parser

import (
	"bufio"
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Parser interface {
	Run(wg *sync.WaitGroup)
	parse()
}

type Status int

const (
	Ready Status = iota
	Processing
)

type parser struct {
	isRunning bool
	config    config.Config
	channels  *config.Channels
	errch     chan error
	status    Status
	data      map[models.Year]map[time.Month][]models.LfcAce
}

// NewService returns new Parser instance
func NewService(cfg config.Config, ch *config.Channels) Parser {
	return &parser{
		config:   cfg,
		channels: ch,
		errch:    make(chan error, 1),
		status:   Ready,
		data:     make(map[models.Year]map[time.Month][]models.LfcAce),
	}
}

func (p *parser) Run(wg *sync.WaitGroup) {
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
		case <-p.channels.RunParse:
			if p.isRunning { //TODO check if parser is ready
				//log.Printf("-= Parse %d files =-", len(p.filenames))
				if p.status == Ready {
					p.status = Processing
					go p.parse()
				}
			}
		case p.isRunning = <-p.channels.ParserIsRunning:
			log.Printf("Parser is running: %v\n", p.isRunning)

		case err := <-p.errch:
			if err != nil {
				log.Fatalf("Parsing input data failed, err: %v\n", err)
				return
			}
			//
			log.Printf("Parsing completed isRunning:%v  status:%v", p.isRunning, p.status)
			if p.isRunning && p.status == Ready {
				p.channels.RunProcess <- p.data
			}
		//case <-p.channels.CfgUpdate:
		//	// TODO config update
		//	log.Println("Parser updates config")
		case p.isRunning = <-p.channels.Quit:
			// TODO should wait until parser completes its job
			log.Printf("Parser says Bye bye.. status:%v", p.status)
			return
		}
	}
}

func (p *parser) parse() {
	filenames, err := p.collectFilenames()
	if err != nil {
		log.Fatal(err)
		//p.errch <- err
	}

	for _, filename := range filenames {
		if strings.HasPrefix(filename, ".") {
			continue
		}
		f, err := os.Open(filepath.Join(p.config.Params.InputDir, filename))

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
			//log.Fatal(err)
			p.errch <- err
		}
	}

	p.status = Ready
	p.errch <- nil
}

func (p *parser) collectFilenames() ([]string, error) {
	var filenames []string

	inputDir := p.config.Params.InputDir
	c, err := os.ReadDir(inputDir)
	if err != nil {
		//log.Fatalf("Can't read input dir %s", inputDir)
		return filenames, err
	}
	log.Println("Listing input dir")
	for _, entry := range c {
		fmt.Println(" ", entry.Name())

		// collect filename
		filenames = append(filenames, entry.Name())
		continue
	}
	return filenames, nil
}

func (p *parser) parseLine(line string) {
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
	//var measures []models.LfcAce
	//measures = append(measures, measure)

	_, exists := p.data[models.Year(measure.AvgTime.Year())]
	if !exists {
		p.data[models.Year(measure.AvgTime.Year())] = make(map[time.Month][]models.LfcAce)
	}

	p.data[models.Year(measure.AvgTime.Year())][measure.AvgTime.Month()] =
		append(p.data[models.Year(measure.AvgTime.Year())][measure.AvgTime.Month()], measure)
}
