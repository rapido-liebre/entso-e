package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"go-fiber-api-docker/pkg/archiver"
	"go-fiber-api-docker/pkg/common/config"
	"go-fiber-api-docker/pkg/generator"
	"go-fiber-api-docker/pkg/watcher"

	"github.com/onsi/ginkgo/v2"
)

func Test_main_integration(t *testing.T) {
	defer ginkgo.GinkgoRecover()

	//given
	ch := config.GetChannels()
	cfg, err := config.GetConfig()
	if err != nil {
		t.Errorf("Loading config error = %v, wantErr %v", err, false)
	}
	time.Sleep(2 * time.Second) //wait for cfg is loaded
	gen := generator.NewService(cfg, &ch, "2022-08-09 09.37.csv")

	var wg sync.WaitGroup
	wg.Add(2)

	go gen.Run()
	arch := archiver.NewService(cfg, &ch)
	go arch.Run(&wg)
	watch := watcher.NewService(cfg, &ch, 3*time.Second)
	go watch.Run(&wg)

	// when
	ch.RunGenerate <- true //start generate files
	time.Sleep(4 * time.Second)
	ch.RunWatch <- true // start watching next files

	// then
	timer1 := time.NewTimer(500 * time.Second)
	done := make(chan bool)
	go func() {
		<-timer1.C
		fmt.Println("Timer 1 fired")
		done <- true
	}()
	<-done
	close(ch.Quit)
	wg.Done()
}
