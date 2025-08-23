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
	"github.com/vinitparekh17/syncsnipe/internal/database"
	s "github.com/vinitparekh17/syncsnipe/internal/sync"
)

const (
	serviceName        = "SyncSnipeService"
	serviceDisplayName = "SyncSnipe Service"
	serviceDescription = "A service to manage file synchronization tasks."
	shutdownTimeout    = 30 * time.Second
	appName            = "syncsnipe"
	logFileName        = "service.log"
	logFilePerms       = 0644
	logDirPerms        = 0755
)

// ServiceManager handles service operations
type ServiceManager struct {
	svc    service.Service
	logger service.Logger
	dbTx   *database.Queries
}

// NewServiceCmd creates the service management command
func NewServiceCmd(q *database.Queries) *cobra.Command {
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the SyncSnipe service",
		Long: `SyncSnipe service allows you to run the application as a background service.
Use the available subcommands to install, start, stop, or uninstall the service.`,
	}

	manager, err := NewServiceManager(q)
	if err != nil {
		log.Fatal(err)
	}

	serviceCmd.AddCommand(
		manager.newRunServiceCmd(), // Add run command for foreground execution
		manager.newStartServiceCmd(),
		manager.newStopServiceCmd(),
		manager.newInstallServiceCmd(),
		manager.newUninstallServiceCmd(),
		manager.newStatusServiceCmd(),
	)

	return serviceCmd
}

// NewServiceManager creates a new service manager
func NewServiceManager(q *database.Queries) (*ServiceManager, error) {
	// Get the current executable path for proper service configuration
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	prg := &program{
		dbTx: q,
	}

	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
		Executable:  execPath,
		Arguments:   []string{"service", "run"}, // Important: tell service how to run
	}

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	logger, err := svc.Logger(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Set the logger in the program
	prg.logger = logger

	return &ServiceManager{
		svc:    svc,
		logger: logger,
		dbTx:   q,
	}, nil
}

// newRunServiceCmd creates the run command for foreground service execution
func (sm *ServiceManager) newRunServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the SyncSnipe service (used internally)",
		Long:  "Run the SyncSnipe service. This is typically used internally by the service manager.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// This is called by the service manager, not directly by users
			return sm.svc.Run()
		},
	}
}

func (sm *ServiceManager) newStartServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the SyncSnipe service",
		Long:  "Start the SyncSnipe service to run in the background and manage file synchronization tasks.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if service is installed first
			status, err := sm.svc.Status()
			if err != nil {
				fmt.Println("Service may not be installed. Try running 'install' first.")
				return fmt.Errorf("failed to get service status: %w", err)
			}

			if status == service.StatusRunning {
				fmt.Println("Service is already running.")
				return nil
			}

			if err := sm.svc.Start(); err != nil {
				return fmt.Errorf("failed to start service: %w", err)
			}

			fmt.Println("Service started successfully.")

			// Give it a moment to start and check status
			time.Sleep(2 * time.Second)
			if status, err := sm.svc.Status(); err == nil {
				switch status {
				case service.StatusRunning:
					fmt.Println("Service is now running.")
				case service.StatusStopped:
					fmt.Println("Warning: Service appears to have stopped immediately. Check logs.")
				}
			}

			return nil
		},
	}
}

func (sm *ServiceManager) newStopServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the SyncSnipe service",
		Long:  "Stop the SyncSnipe service that is currently running in the background.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := sm.svc.Stop(); err != nil {
				return fmt.Errorf("failed to stop service: %w", err)
			}
			fmt.Println("Service stopped successfully.")
			return nil
		},
	}
}

func (sm *ServiceManager) newInstallServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install the SyncSnipe service",
		Long:  "Install the SyncSnipe service to run in the background and manage file synchronization tasks.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := sm.svc.Install(); err != nil {
				return fmt.Errorf("failed to install service: %w", err)
			}
			fmt.Println("Service installed successfully.")
			fmt.Println("You can now use 'start' to start the service.")
			return nil
		},
	}
}

func (sm *ServiceManager) newUninstallServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall the SyncSnipe service",
		Long:  "Uninstall the SyncSnipe service that is running in the background.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Stop the service first if it's running
			if status, err := sm.svc.Status(); err == nil && status == service.StatusRunning {
				fmt.Println("Stopping service before uninstalling...")
				sm.svc.Stop()
				time.Sleep(2 * time.Second)
			}

			if err := sm.svc.Uninstall(); err != nil {
				return fmt.Errorf("failed to uninstall service: %w", err)
			}
			fmt.Println("Service uninstalled successfully.")
			return nil
		},
	}
}

func (sm *ServiceManager) newStatusServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check SyncSnipe service status",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := sm.svc.Status()
			if err != nil {
				fmt.Printf("Failed to get service status: %v\n", err)
				fmt.Println("Service may not be installed. Try running 'install' first.")
				return nil
			}

			switch status {
			case service.StatusRunning:
				fmt.Println("✓ Service is running.")
				if stats, err := sm.getSyncStats(); err == nil {
					fmt.Printf("Sync statistics: %s\n", stats)
				}
			case service.StatusStopped:
				fmt.Println("✗ Service is stopped.")
			default:
				fmt.Println("? Service status unknown.")
			}

			// Show log path for debugging
			if logPath, err := getLogPath(); err == nil {
				fmt.Printf("Log file: %s\n", logPath)
			}

			return nil
		},
	}
}

// getSyncStats retrieves synchronization statistics from the database
func (sm *ServiceManager) getSyncStats() (string, error) {
	// TODO: Implement actual database query for sync statistics
	return "No statistics available", nil
}

// program implements the service.Interface
type program struct {
	mu      sync.RWMutex
	logFile *os.File
	logger  service.Logger
	dbTx    *database.Queries
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	running bool
}

// Start implements service.Interface.Start
func (p *program) Start(s service.Service) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("service is already running")
	}

	// Setup logging early
	if err := p.setupLogging(); err != nil {
		if p.logger != nil {
			p.logger.Errorf("Failed to setup logging: %v", err)
		}
		return fmt.Errorf("failed to setup logging: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.running = true

	p.wg.Add(1)
	go p.run(ctx)

	// Give the service a moment to initialize
	time.Sleep(100 * time.Millisecond)

	if p.logger != nil {
		p.logger.Info("Service Start() completed successfully")
	}

	return nil
}

// Stop implements service.Interface.Stop
func (p *program) Stop(s service.Service) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return nil // Already stopped
	}

	if p.logger != nil {
		p.logger.Info("Service Stop() initiated")
	}

	if p.cancel != nil {
		p.cancel()
	}

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if p.logger != nil {
			p.logger.Info("Service stopped gracefully")
		}
	case <-time.After(shutdownTimeout):
		if p.logger != nil {
			p.logger.Warning("Graceful shutdown timeout, forcing stop")
		}
	}

	p.cleanup()
	p.running = false
	return nil
}

// run is the main service execution loop
func (p *program) run(ctx context.Context) {
	defer p.wg.Done()

	if p.logger != nil {
		p.logger.Info("Service main loop starting")
	}

	watcher, err := p.createWatcher()
	if err != nil {
		if p.logger != nil {
			p.logger.Errorf("Failed to create watcher: %v", err)
		}
		return
	}

	if p.logger != nil {
		p.logger.Info("SyncSnipe service started successfully")
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		watcher.Start(ctx)
	}()

	<-ctx.Done()

	if p.logger != nil {
		p.logger.Info("Service shutdown initiated")
	}
}

func (p *program) createWatcher() (*s.SyncWatcher, error) {
	watcher, err := s.NewSyncWatcher(p.dbTx)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync watcher: %w", err)
	}
	return watcher, nil
}

func (p *program) setupLogging() error {
	logPath, err := getLogPath()
	if err != nil {
		return fmt.Errorf("failed to get log path: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, logFilePerms)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", logPath, err)
	}

	p.logFile = logFile

	// Configure multi-writer for both file and stdout (when running as service, stdout goes to service logs)
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)

	if p.logger != nil {
		p.logger.Infof("Logging initialized: %s", logPath)
	}

	return nil
}

func (p *program) cleanup() {
	if p.logFile != nil {
		if err := p.logFile.Close(); err != nil {
			if p.logger != nil {
				p.logger.Errorf("Failed to close log file: %v", err)
			}
		}
		p.logFile = nil
	}
}

func getLogPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to temp directory
		tempLogDir := filepath.Join(os.TempDir(), appName)
		if err := os.MkdirAll(tempLogDir, logDirPerms); err != nil {
			return "", fmt.Errorf("failed to create temp log directory: %w", err)
		}
		return filepath.Join(tempLogDir, logFileName), nil
	}

	appLogDir := filepath.Join(configDir, appName)
	if err := os.MkdirAll(appLogDir, logDirPerms); err != nil {
		return "", fmt.Errorf("failed to create log directory %s: %w", appLogDir, err)
	}

	return filepath.Join(appLogDir, logFileName), nil
}
