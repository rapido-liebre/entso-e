package archiver

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go-fiber-api-docker/pkg/common/config"
)

// Archiver Package
type Archiver interface {
	Run(wg *sync.WaitGroup)
	archive()
	zipSource(prefix, source, target string) error
}

type Status int

const (
	Ready Status = iota
	Processing
)

type archiver struct {
	isRunning bool
	cfg       config.Config
	channels  *config.Channels
	filenames map[string]bool
	errch     chan error
	status    Status
}

// NewService returns new Archiver instance
func NewService(cfg config.Config, ch *config.Channels) Archiver {
	return &archiver{
		cfg:       cfg,
		channels:  ch,
		filenames: make(map[string]bool),
		errch:     make(chan error, 1),
		status:    Ready,
	}
}

func (arch archiver) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	// proceed in infinite loop
	for {
		select {
		case filename := <-arch.channels.Filename:
			//using map here instead of slice for easier lookup during processing
			if len(filename) > 0 {
				arch.filenames[filename] = true
				log.Printf("Queued: %s", filename)
			}
		case <-arch.channels.RunArchive:
			if arch.isRunning { //TODO check if arch is ready
				log.Printf("-= Archive %d files =-", len(arch.filenames))
				if arch.status == Ready {
					go arch.archive()
				}
			}
		case arch.isRunning = <-arch.channels.ArchIsRunning:
			log.Printf("Archiver is running: %v\n", arch.isRunning)

		case err := <-arch.errch:
			if err != nil {
				log.Fatalf("Creating zip failed, err: %v\n", err)
				return
			}
			log.Printf("Zip created successfully isRunning:%v  status:%v", arch.isRunning, arch.status)
			if arch.isRunning && arch.status == Ready {
				arch.channels.RunWatch <- true
			}
		case <-arch.channels.CfgUpdate:
			// TODO config update
			log.Println("Archiver updates config")
		case arch.isRunning = <-arch.channels.Quit:
			// TODO should wait until archiver completes its job
			log.Printf("Archiver says Bye bye.. status:%v", arch.status)
			return
		}
	}
}

// archive zips collected files. https://gosamples.dev/zip-file/
func (arch *archiver) archive() {
	// return if nothing to do..
	if len(arch.filenames) == 0 {
		arch.errch <- nil
		return
	}
	arch.status = Processing

	srcDir := arch.cfg.Params.InputDir
	destDir := arch.cfg.Params.OutputDir
	prefix := "test/"
	zipName := "archive.zip"

	// extract valid zip name
	for filename, _ := range arch.filenames {
		if strings.HasSuffix(filename, ".csv") {
			prefix = filename[0:13] + "/"
			zipName = filename[0:13] + ".00.zip"
			//log.Println("prefix:", prefix)
		}
	}

	zipFile := filepath.Join(destDir, zipName)
	log.Printf("Creating zip archive %s", zipFile)

	if err := arch.zipSource(prefix, srcDir, zipFile); err != nil {
		arch.errch <- err
	} else {
		if config.UseMinio() {
			log.Println("arch.zipSource()")
			mc, err := config.GetMinioClient()
			if err != nil {
				arch.errch <- err
			}
			if err = mc.PutObject(arch.cfg.Params.BucketName, zipFile); err != nil {
				arch.errch <- err
			}
			//TODO - remove uploaded file from local disk
		}
	}
	arch.filenames = map[string]bool{}
	arch.status = Ready
	arch.errch <- nil
}

func (arch archiver) zipSource(prefix, source, target string) error {
	// create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip directories or files excluded from zip
		if info.IsDir() || !arch.filenames[info.Name()] {
			return nil
		}

		// create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// set the header name
		header.Name = prefix + info.Name()

		// create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		// copy content of the source file to archive
		f1, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f1.Close()

		_, err = io.Copy(headerWriter, f1)
		if err != nil {
			return err
		}

		return err
	})
}
