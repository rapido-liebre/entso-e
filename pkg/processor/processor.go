package processor

import (
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/common/models"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

type Processor interface {
	Run(wg *sync.WaitGroup)
	process()
}

type Status int

const (
	Ready Status = iota
	Processing
)

const FRRP = 1075
const FRRM = -1075

const level1 = 124.964
const level2 = 236.326

type processor struct {
	isRunning bool
	status    Status
	config    config.Config
	channels  *config.Channels
	errch     chan error
	data      map[models.Year]map[time.Month][]models.LfcAce
}

// NewService returns new Processor instance
func NewService(cfg config.Config, ch *config.Channels) Processor {
	return &processor{
		config:   cfg,
		channels: ch,
		errch:    make(chan error, 1),
		status:   Ready,
	}
}

func (p *processor) Run(wg *sync.WaitGroup) {
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
		case p.data = <-p.channels.RunProcess:
			if p.isRunning { //TODO check if parser is ready
				//log.Printf("-= Parse %d files =-", len(p.filenames))
				if p.status == Ready {
					p.status = Processing
					go p.process()
				}
			}
		case p.isRunning = <-p.channels.ProcessorIsRunning:
			log.Printf("Processor is running: %v\n", p.isRunning)

		case err := <-p.errch:
			if err != nil {
				log.Fatalf("Processing data failed, err: %v\n", err)
				return
			}
			//
			log.Printf("Processing data completed isRunning:%v  status:%v", p.isRunning, p.status)
			//if p.isRunning && p.status == Ready {
			//	p.channels.RunProcess <- p.data
			//}
		//case <-p.channels.CfgUpdate:
		//	// TODO config update
		//	log.Println("Processor updates config")
		case p.isRunning = <-p.channels.Quit:
			// TODO should wait until processor completes its job
			log.Printf("Processor says Bye bye.. status:%v", p.status)
			return
		}
	}
}

func (p *processor) process() {
	if len(p.data) == 0 {
		//TODO uncomment if ready for processing
		//p.errch <- errors.New("no data for processing")

		p.status = Ready
		p.errch <- nil
	}

	var (
		sum, sumq, lv1pos, lv1neg, lv2pos, lv2neg, perc1 /*, perc5, perc10, perc90, perc95, perc99*/ float64
	)

	for k, v := range p.data[2022][6] {
		fmt.Println(k, v)
		sum += v.AvgValue
		sumq += math.Pow(v.AvgValue, 2)

		if -v.AvgValue > level1 {
			lv1pos += 1
		}
		if -v.AvgValue < -level1 {
			lv1neg += 1
		}
		if -v.AvgValue > level2 {
			lv2pos += 1
		}
		if -v.AvgValue < -level2 {
			lv2neg += 1
		}

	}
	avg := sum / float64(len(p.data[2022][6]))
	avg = -avg
	fmt.Println(avg)

	dev := math.Sqrt(sumq / float64(len(p.data[2022][6])))
	fmt.Println(dev)

	fmt.Println("Level1 +:", lv1pos)
	fmt.Println("Level1 -:", lv1neg)
	fmt.Println("Level2 +:", lv2pos)
	fmt.Println("Level2 -:", lv2neg)

	perc1 = float64(len(p.data[2022][6])) / float64(100*(100-1))
	fmt.Println("Perc 1:", perc1)

	//for _, filename := range filenames {
	//	if strings.HasPrefix(filename, ".") {
	//		continue
	//	}
	//	f, err := os.Open(filepath.Join(p.config.Params.InputDir, filename))
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	defer f.Close()
	//
	//	scanner := bufio.NewScanner(f)
	//	var rowsFetched int
	//
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//		//valid line with data starts with date like '| 2022-12-01 '
	//		if strings.Contains(line, "| 2") {
	//			//fmt.Println(line)
	//			p.parseLine(line)
	//			rowsFetched++
	//		}
	//		//if strings.Contains(line, "rows fetched") {
	//		//	fmt.Println(scanner.Text())
	//		//}
	//	}
	//	fmt.Println(rowsFetched)
	//
	//	if err := scanner.Err(); err != nil {
	//		//log.Fatal(err)
	//		p.errch <- err
	//	}
	//}

	p.status = Ready
	p.errch <- nil
}

//func (p *processor) collectFilenames() ([]string, error) {
//	var filenames []string
//
//	inputDir := p.config.Params.InputDir
//	c, err := os.ReadDir(inputDir)
//	if err != nil {
//		//log.Fatalf("Can't read input dir %s", inputDir)
//		return filenames, err
//	}
//	log.Println("Listing input dir")
//	for _, entry := range c {
//		fmt.Println(" ", entry.Name())
//
//		// collect filename
//		filenames = append(filenames, entry.Name())
//		continue
//	}
//	return filenames, nil
//}
//
//func (p *processor) parseLine(line string) {
//	var measure models.LfcAce
//	var result []string
//
//	//line = strings.Replace(line, "|", " ", -1)
//	split := strings.Split(line, "|")
//	for _, s := range split {
//		strings.Trim(s, " ")
//		if len(s) != 0 {
//			fmt.Println(s)
//			result = append(result, strings.Trim(s, " "))
//		}
//	}
//	//6 values expected
//	if len(result) != 6 {
//		log.Fatalf("parsing line failed. Line: %s", line)
//	}
//
//	for idx, val := range result {
//		var err error
//		switch models.LfcAceIndex(idx) {
//		case models.AvgTime:
//			measure.AvgTime, err = time.Parse(time.DateTime, val)
//			if err != nil {
//				log.Fatal("Could not parse avg time:", err)
//			}
//		case models.SaveTime:
//			measure.SaveTime, err = time.Parse(time.DateTime, val)
//			if err != nil {
//				log.Fatal("Could not parse save time:", val, err)
//			}
//		case models.AvgName:
//			measure.AvgName = val
//		case models.AvgValue:
//			measure.AvgValue, err = strconv.ParseFloat(val, 64)
//			if err != nil {
//				log.Fatal("Could not parse float value:", val, err)
//			}
//		case models.AvgStatus:
//			measure.AvgStatus, err = strconv.Atoi(val)
//			if err != nil {
//				log.Fatal("Could not parse status value:", val, err)
//			}
//		case models.SystemSite:
//			if len(val) != 1 {
//				log.Fatal("Could not parse byte value:", val, err)
//			}
//			measure.SystemSite = val
//		}
//	}
//	p.data = append(p.data, measure)
//}
