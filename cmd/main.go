package main

import (
	"entso-e_reports/pkg/api"
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/parser"
	"entso-e_reports/pkg/processor"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"sync"
)

func main() {
	// parse args
	cfgPathPtr := flag.String("config", "./", "Specific path to config file")
	runGeneratorPtr := flag.Bool("run_generator", false, "Run generator")
	flag.Parse()
	log.Printf("runGeneratorPtr: %v\n", *runGeneratorPtr)

	// load and parse config
	cfg, err := config.GetConfig(*cfgPathPtr)
	if err != nil {
		log.Fatalln("Failed at loading init config", err)
		return
	}

	// validate config fields
	if err = cfg.Validate(); err != nil {
		log.Fatalln("Failed at validating init config", err)
		return
	}

	// initialize channels
	channels := config.GetChannels()

	// initialize services
	parser := parser.NewService(cfg, &channels)
	processor := processor.NewService(cfg, &channels)
	//if *runGeneratorPtr {
	//	go generator.RunGenerator(cfg, &channels)
	//}

	app := fiber.New()
	// Default config enables CORS middleware for fiber https://docs.gofiber.io/api/middleware/cors
	app.Use(cors.New())

	// initialize API
	api.RegisterRoutes(app, &channels)

	var wg sync.WaitGroup
	wg.Add(2)
	go parser.Run(&wg)
	go processor.Run(&wg)

	//start listening on port
	go func() {
		err = app.Listen(cfg.Params.Port)
		if err != nil {
			log.Printf("Cannot start listening the port %s. \n[%+v]", cfg.Params.Port, err)
			close(channels.APIChannels.Quit)
		}
	}()

	wg.Wait()

	log.Println("Main says Bye bye")
}
