package generator

import (
	"fmt"
	"log"
	"time"

	"go-fiber-api-docker/pkg/common/config"
)

func RunGenerator(cfg config.Config, ch *config.Channels) {
	log.Println("RunGenerator()")
	gen := NewService(cfg, ch, "2022-08-09 09.37.csv")
	go gen.Run()

	// when
	//dc.Filename <- "2022-08-09 09.37.csv"
	ch.RunGenerate <- true //start generate files

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

	log.Println("RunGenerator says Bye bye")
}
