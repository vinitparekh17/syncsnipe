package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	s "github.com/vinitparekh17/syncsnipe/internal/sync"
)

var logger service.Logger

func NewServiceCmd(q *database.Queries) *cobra.Command {
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the SyncSnipe service",
		Long: `SyncSnipe service allows you to run the application as a background service.
Use the available subcommands to install, start, stop, or uninstall the service.`,
	}

	prg := &program{
		dbTx: q,
	}
	svcConfig := &service.Config{
		Name:        "SyncSnipeService",
		DisplayName: "SyncSnipe Service",
		Description: "A service to manage file synchronization tasks.",
	}

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	serviceCmd.AddCommand(
		prg.newStartServiceCmd(svc),
		newStopServiceCmd(svc),
		newInstallServiceCmd(svc),
		newUninstallServiceCmd(svc),
		newStatusServiceCmd(svc),
	)

	return serviceCmd
}

func (p *program) newStartServiceCmd(svc service.Service) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the SyncSnipe service",
		Long:  "Start the SyncSnipe service to run in the background and manage file synchronization tasks.",
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := svc.Start(); err != nil {
				return fmt.Errorf("failed to start service: %w", err)
			}

			fmt.Println("Service started successfully.")
			return nil
		},
	}

	return startCmd
}

func newStopServiceCmd(svc service.Service) *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the SyncSnipe service",
		Long:  "Stop the SyncSnipe service that is currently running in the background.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Stop(); err != nil {
				return fmt.Errorf("failed to stop service: %w", err)
			}

			fmt.Println("Service stopped successfully.")
			return nil
		},
	}

	return stopCmd
}

func newInstallServiceCmd(svc service.Service) *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install the SyncSnipe service",
		Long:  "Install the SyncSnipe service to run in the background and manage file synchronization tasks.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Install(); err != nil {
				return fmt.Errorf("failed to install service: %w", err)
			}

			fmt.Println("Service installed successfully.")
			return nil
		},
	}

	return installCmd
}

func newUninstallServiceCmd(svc service.Service) *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall the SyncSnipe service",
		Long:  "Uninstall the SyncSnipe service that is running in the background.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Uninstall(); err != nil {
				return fmt.Errorf("failed to uninstall service: %w", err)
			}

			fmt.Println("Service uninstalled successfully.")
			return nil
		},
	}

	return uninstallCmd
}

func newStatusServiceCmd(svc service.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check SyncSnipe service status",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := svc.Status()
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}
			switch status {
			case service.StatusRunning:
				fmt.Println("Service is running.")
				// TODO: Query DB for sync stats
			case service.StatusStopped:
				fmt.Println("Service is stopped.")
			default:
				fmt.Println("Service status unknown.")
			}
			return nil
		},
	}
}

type program struct {
	mu      sync.Mutex
	logFile *os.File
	logger  service.Logger
	dbTx    *database.Queries
	cancel  context.CancelFunc // Add a cancel function to the struct
	wg      sync.WaitGroup
}

func (p *program) run(ctx context.Context) {
	defer p.wg.Done()

	if err := p.setupLogging(); err != nil {
		p.logger.Errorf("Failed to setup logging: %v", err)
		return
	}

	watcher, err := s.NewSyncWatcher(p.dbTx)
	if err != nil {
		p.logger.Errorf("Failed to create watcher: %v", err)
		return
	}

	// Start the watcher in a goroutine
	var watcherWg sync.WaitGroup
	watcherWg.Add(1)
	go func() {
		defer watcherWg.Done()
		watcher.Start(ctx)
	}()

	colorlog.Info("SyncSnipe service started. Waiting for stop signal.")

	<-ctx.Done()

	// Wait for the watcher to finish its cleanup
	watcherWg.Wait()
	colorlog.Info("Watcher stopped gracefully.")
}

func (p *program) setupLogging() error {
	logPath := getLogPath()
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open log file: %w", err)
	}

	p.logFile = logFile
	// Log to both service logger and file
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	return nil
}

func (p *program) Start(s service.Service) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel != nil {
		return fmt.Errorf("service already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	p.wg.Add(1)
	go p.run(ctx)
	return nil
}

func (p *program) Stop(s service.Service) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Graceful shutdown completed
	case <-time.After(30 * time.Second):
		p.logger.Warning("Graceful shutdown timeout, forcing stop")
	}
	if p.logFile != nil {
		p.logFile.Close()
	}
	colorlog.Info("Service stopped successfully.")
	return nil
}

// getLogPath returns a log path appropriate for the OS
func getLogPath() string {
	appName := "syncsnipe"

	// os.UserConfigDir() is generally preferred for config/log data
	// It finds paths like:
	// - Windows: %AppData%\syncsnipe
	// - macOS:   ~/Library/Application Support/syncsnipe
	// - Linux:   ~/.config/syncsnipe
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback for headless systems or weird environments
		return filepath.Join(os.TempDir(), appName, "service.log")
	}

	appLogDir := filepath.Join(configDir, appName)
	os.MkdirAll(appLogDir, 0755)
	return filepath.Join(appLogDir, "service.log")
}
