package generator

import (
	"fmt"
	"testing"
	"time"

	"go-fiber-api-docker/pkg/common/config"

	"github.com/onsi/ginkgo/v2"
)

func Test_generator_generate(t *testing.T) {
	defer ginkgo.GinkgoRecover()

	//given
	dc := config.GetChannels()
	cfg, err := config.GetConfig()
	if err != nil {
		t.Errorf("Loading config error = %v, wantErr %v", err, false)
	}

	gen := NewService(cfg, &dc, "2022-08-09 09.37.csv")
	go gen.Run()

	// when
	//dc.Filename <- "2022-08-09 09.37.csv"
	dc.RunGenerate <- true //start generate files

	// then
	timer1 := time.NewTimer(500 * time.Second)
	done := make(chan bool)
	go func() {
		<-timer1.C
		fmt.Println("Timer 1 fired")
		done <- true
	}()
	<-done
	close(dc.Quit)
}
