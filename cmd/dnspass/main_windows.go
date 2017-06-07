// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/freman/dnspass"
	"github.com/freman/eventloghook"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var windowsCommands = []cli.Command{
	cli.Command{
		Name:   "install",
		Usage:  "Create a service in the local service control manager",
		Action: installService,
	},
	cli.Command{
		Name:   "remove",
		Usage:  "Remove an existing service from the local service control manager",
		Action: removeService,
	},
	cli.Command{
		Name:   "start",
		Usage:  "Start the service in the local service control manager",
		Action: startService,
	},
	cli.Command{
		Name:   "stop",
		Usage:  "Stop the service in the local service control manager",
		Action: stopService,
	},
	cli.Command{
		Name:   "run",
		Usage:  "Run the service interactively",
		Action: runService,
	},
}

var (
	runningStatus = svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop,
	}
	stoppingStatus = svc.Status{
		State:   svc.StopPending,
		Accepts: svc.AcceptStop,
	}
)

func installService(c *cli.Context) error {
	exePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to find the dnspass executable: %s\n", err)
	}

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("Failed to connect to service control manager: %s\n", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(svcName)
	if err == nil {
		s.Close()
		return fmt.Errorf("Service %s already exists", svcName)
	}

	mgrConfig := mgr.Config{
		DisplayName:  svcName,
		Description:  svcDesc,
		Dependencies: []string{"Netman", "Dnscache", "Eventlog"},
		StartType:    mgr.StartAutomatic,
	}

	s, err = m.CreateService(svcName, exePath, mgrConfig)
	if err != nil {
		return fmt.Errorf("Failed to install service: %s\n", err)
	}
	defer s.Close()

	err = eventlog.InstallAsEventCreate(svcName, eventlog.Error|eventlog.Warning|eventlog.Info)

	if err != nil {
		s.Delete()
		return fmt.Errorf("Failed to install event log registry entries: %s", err)
	}
	return nil
}

func removeService(c *cli.Context) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("Failed to connect to service control manager: %s\n", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(svcName)
	if err != nil {
		return fmt.Errorf("Service %s is not installed", svcName)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		return fmt.Errorf("Failed to delete service: %s\n", err)
	}

	err = eventlog.Remove(svcName)
	if err != nil {
		return fmt.Errorf("Failed to remove event log registry entries: %s\n", err)
	}
	return nil
}

func stopService(c *cli.Context) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("Failed to connect to service control manager: %s\n", err)
	}
	defer m.Disconnect()

	service, err := m.OpenService(svcName)
	if err != nil {
		return fmt.Errorf("Failed to open service: %s\n", err)
	}
	defer service.Close()

	_, err = service.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("Failed to stop service: %s\n", err)
	}

	fmt.Printf("Service stopped successfully")
	return nil
}

func startService(c *cli.Context) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("Failed to connect to service control manager: %s\n", err)
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err != nil {
		return fmt.Errorf("Failed to open service: %s\n", err)
	}
	defer s.Close()
	err = s.Start()
	if err != nil {
		return fmt.Errorf("Failed to start service: %s\n", err)
	}
	return nil
}

type service struct {
	context     *cli.Context
	interactive bool
	log         logrus.FieldLogger
}

func runService(c *cli.Context) error {
	interactive, err := svc.IsAnInteractiveSession()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to detect interactive session: %s", err), 1)
	}

	run := svc.Run
	if interactive {
		run = debug.Run
	}

	elog, err := eventlog.Open(svcName)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Unable to open the event log: %s", err), 1)
	}
	defer elog.Close()

	log := logrus.New()
	if c.Bool("log-debug") {
		log.Level = logrus.DebugLevel
	}
	log.Hooks.Add(eventloghook.NewHook(elog))

	s := service{
		context:     c,
		interactive: interactive,
		log:         log,
	}

	if err := run(svcName, &s); err != nil {
		err = cli.NewExitError(err.Error(), 1)
	}
	return err
}

func (service *service) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	configFilename := service.context.String("config-filename")
	if configFilename == "config.toml" {
		exePath, err := filepath.Abs(os.Args[0])
		if err != nil {
			service.log.Fatalf("Failed to find the dnspass executable: %s", err)
		}
		configFilename = filepath.Join(filepath.Dir(exePath), configFilename)
	}

	server := &dnspass.Server{
		ConfigFile: configFilename,
		Log:        service.log,
	}

	err := server.Run()
	if err != nil {
		service.log.Fatalf("Service failed to run: %s", err)
	}

	s <- runningStatus
	for {
		request := <-r
		switch request.Cmd {
		case svc.Interrogate:
			s <- runningStatus

		case svc.Stop:
			s <- stoppingStatus
			service.log.Info("Shutting down")

			time.AfterFunc(gracefulTimeout, func() {
				service.log.Fatal("Failed to stop quickly; stopping forcefully")
				os.Exit(1)
			})

			server.Shutdown()

			return false, 0

		default:
			service.log.Errorf("Received unsupported service command from service control manager: %d", request.Cmd)
		}
	}
}

func main() {
	interactive, err := svc.IsAnInteractiveSession()
	if err != nil {
		fmt.Errorf("Failed to detect interactive session: %s", err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = svcName
	app.Usage = svcName + " - " + svcDesc
	app.Version = dnspass.Version
	app.Flags = defaultFlags

	if !interactive {
		app.Action = runService
	} else {
		app.Commands = windowsCommands
	}

	app.Run(os.Args)
}
