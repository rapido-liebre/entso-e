package main

import (
	"entso-e_reports/pkg/api"
	"entso-e_reports/pkg/common/config"
	"entso-e_reports/pkg/db_connector"
	"entso-e_reports/pkg/parser"
	"entso-e_reports/pkg/processor"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/takama/daemon"
)

const (

	// name of the service
	name        = "entsoe"
	description = "Generate Entso-e Reports Service"

	// port which daemon should be listen
	port = ":9977"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage(cfgPathPtr string) (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Do something, call your goroutines, etc
	runLocally(cfgPathPtr)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf[:numbytes])
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	// parse args
	cfgPathPtr := flag.String("config", "./", "Specific path to config file")
	noDeamonPtr := flag.Bool("no-deamon", false, "Run locally or as a deamon")
	flag.Parse()
	log.Printf("noDeamonPtr: %v\n", *noDeamonPtr)
	log.Printf("App name: %s\n", config.GetAppName())

	if *noDeamonPtr {
		runLocally(*cfgPathPtr)
	} else {
		srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
		if err != nil {
			errlog.Println("Error: ", err)
			os.Exit(1)
		}
		service := &Service{srv}
		status, err := service.Manage(*cfgPathPtr)
		if err != nil {
			errlog.Println(status, "\nError: ", err)
			os.Exit(1)
		}
		fmt.Println(status)

		//service, err := daemon.New("name", "description", daemon.SystemDaemon)
		//if err != nil {
		//	log.Fatal("Error: ", err)
		//}
		//status, err := service.Install()
		//if err != nil {
		//	log.Fatal(status, "\nError: ", err)
		//}
		//fmt.Println(status)

		//runLocally(*cfgPathPtr)
	}

	log.Println("Main says Bye bye")
}

func runLocally(cfgPathPtr string) {
	// load and parse config
	cfg, err := config.GetConfig(cfgPathPtr)
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
	dbConnector := db_connector.NewService(cfg, &channels)
	//if *runGeneratorPtr {
	//	go generator.RunGenerator(cfg, &channels)
	//}

	app := fiber.New()
	// Default config enables CORS middleware for fiber https://docs.gofiber.io/api/middleware/cors
	app.Use(cors.New())

	// initialize API
	api.RegisterRoutes(app, &channels)

	var wg sync.WaitGroup
	wg.Add(3)
	go parser.Run(&wg)
	go processor.Run(&wg)
	go dbConnector.Run(&wg)

	//start listening on port
	go func() {
		err = app.Listen(cfg.Params.Port)
		if err != nil {
			log.Printf("Cannot start listening the port %s. \n[%+v]", cfg.Params.Port, err)
			close(channels.APIChannels.Quit)
		}
	}()

	//start the program
	channels.Run()

	wg.Wait()
}
