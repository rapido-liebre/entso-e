package generator

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-fiber-api-docker/pkg/common/config"
)

// Generator Package
type Generator interface {
	Run()
	generate()
}

type generator struct {
	config   config.Config
	channels *config.Channels
	filename string
}

// NewService returns new Generator instance
func NewService(cfg config.Config, ch *config.Channels, filename string) Generator {
	return generator{config: cfg, channels: ch, filename: filename}
}

func (g generator) Run() {
	// proceed in infinite loop
	for {
		select {
		//case g.filename = <-g.dataChan.Filename:
		//	log.Printf("Queued: %s", g.filename)

		case <-g.channels.RunGenerate:
			log.Printf("-= Start generate data files =-")
			go g.generate()
			//if err := g.generate(); err != nil {
			//	// TODO what if an error occures during processing?
			//	log.Fatalln("Failed at generating files", err)
			//}

		case <-g.channels.Quit:
			log.Println("Files generator says Bye bye..")
			return
		}
	}
}

// generate generates data files in destination folder
func (g generator) generate() {
	if len(g.filename) == 0 {
		return
	}

	srcDir := g.config.Params.InputDir
	srcDir = strings.Replace(srcDir, "input", "src", 1)
	dstDir := g.config.Params.InputDir
	src := filepath.Join(srcDir, g.filename)
	dt := time.Now()

	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		dst := filepath.Join(dstDir, dt.Format("2006-01-02 15.04")) + ".csv"
		fmt.Printf("Tock %s\n", dst)
		//"2022-08-09 09.37.csv"

		if err := copyFile(src, dst); err != nil {
			log.Fatal(err)
			return
		}
		dt = dt.Add(1 * time.Minute)
	}
	// run as infinitive loop
	select {}
}

func copyFile(src, dst string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		return err
	}
	return nil
}
