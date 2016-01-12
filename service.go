package speedtestd

import (
	// speedtest "github.com/kogai/speedtest-go"
	"github.com/robfig/cron"
	"github.com/takama/daemon"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	daemon.Daemon
	// port string
}

var (
	port                   = ":9977"
	cronJob     *cron.Cron = cron.New()
	jobSchedule string     = "0 * * * * *"
)

func (service *Service) StartJob() (string, error) {
	// speedtester := speedtest.New()
	// speedtester.SetLogFilePath("./logger.log")
	// speedtester.ToggleLogDest()
	f, _ := os.OpenFile("./logger.log", os.O_RDWR, 0777)
	log.Println("can write")
	f.WriteString("can write")

	cronJob.AddFunc(jobSchedule, func() {
		log.Println("can write with cron")
		f.WriteString("can write with cron")
		// var serverIds []int
		// speedtester.FetchServers()
		// speedtester.ShowResult(serverIds)
	})

	cronJob.Start()
	return service.Start()
}

func (service *Service) StopJob() (string, error) {
	cronJob.Stop()
	return service.Stop()
}

func (service *Service) Manage() (string, error) {
	usage := "Usage: myservice install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.StartJob()
		case "stop":
			return service.StopJob()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

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
			log.Println("Got signal:", killSignal)
			log.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	return usage, nil
}

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
		client.Write(buf)
	}
}
