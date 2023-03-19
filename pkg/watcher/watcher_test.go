package watcher

//func Test_watcher_watch(t *testing.T) {
//	defer ginkgo.GinkgoRecover()
//
//	//given
//	dc := config.GetChannels()
//	cfg, err := config.GetConfig()
//	if err != nil {
//		t.Errorf("Loading config error = %v, wantErr %v", err, false)
//	}
//	var filenames []string
//	interval := 10 * time.Second
//
//	var wg sync.WaitGroup
//	wg.Add(1)
//
//	watch := NewService(cfg, &dc, interval)
//	go watch.Run(&wg)
//
//	// when
//	for n := 0; n <= 1; n++ { //run it twice
//		dc.RunWatch <- true // start watching next files
//		time.Sleep(1 * time.Second)
//		for v := range dc.Filename {
//			log.Printf("Got filename '%s' from channel\n", v)
//			//time.Sleep(1 * time.Second)
//			filenames = append(filenames, v)
//		}
//		log.Printf("Received %d files", len(filenames))
//		<-dc.RunArchive // wait until watcher completes files and triggers an archiver
//	}
//	// then
//	close(dc.Quit)
//	wg.Done()
//}
