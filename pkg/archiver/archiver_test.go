package archiver

//func Test_archiver_archive(t *testing.T) {
//	defer ginkgo.GinkgoRecover()
//
//	//given
//	dc := config.GetChannels()
//	cfg, err := config.GetConfig()
//	if err != nil {
//		t.Errorf("Loading config error = %v, wantErr %v", err, false)
//	}
//	cfg.Params.InputDir = strings.Replace(cfg.Params.InputDir, "input", "src", 1)
//
//	var wg sync.WaitGroup
//	wg.Add(1)
//
//	arch := NewService(cfg, &dc)
//	go arch.Run(&wg)
//
//	// when
//	dc.Filename <- "Stephenson Neal - Peanatema.pdf"
//	dc.Filename <- "latte-espresso.gif"
//	dc.Filename <- "2022-08-09 09.37.csv"
//	dc.Filename <- "obama-shrug.gif"
//	dc.RunArchive <- true
//
//	// then
//	err = arch.archive()
//	<-dc.RunWatch
//	close(dc.Quit)
//	wg.Done()
//	Expect(err).ShouldNot(HaveOccurred())
//}
