package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AtahanPoyraz/TermDrive/server/cmd"
	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/config/database"
	"github.com/AtahanPoyraz/TermDrive/server/config/manager"
	"github.com/AtahanPoyraz/TermDrive/server/config/middleware"
	"github.com/AtahanPoyraz/TermDrive/server/util/fsutil"
	"github.com/gorilla/mux"
)

var (
	logger        *log.Logger
	configuration *config.Configuration
	router        *mux.Router
	server        *http.Server

	sigChan = make(chan os.Signal, 1)
	errChan = make(chan error, 10)

	swg sync.WaitGroup
	err error
)

// init initializes the server configuration, loads necessary configurations,
// sets up the router and middleware, processes command-line arguments,
// and prepares the server and dependencies for execution. It also handles the
// initialization of required directories and the database configuration.
func init() {
	logger = log.New(os.Stdout, "Server API -> ", log.LstdFlags)

	configuration, err = config.LoadConfig("server/config.yaml")
	if err != nil {
		logger.Fatalf("Error loading configuration: %v\n", err)
	}

	database, err := database.SetConfiguration(configuration)
	if err != nil {
		logger.Fatalf("Error initializing database: %v\n", err)
	}

	router = middleware.SetConfiguration(configuration, mux.NewRouter().StrictSlash(true).SkipClean(false))

	dependencies := config.Dependencies{
		AppContext:    config.NewAppContext(configuration, logger),
		SystemContext: config.NewSystemContext(router, database),
	}

	if err := fsutil.Create(configuration.TermDrive.StoragePath); err != nil {
		logger.Fatalf("Error creating storage directory: %v\n", err)
	}

	if err := cmd.ArgumentProcessing(&dependencies); err != nil {
		logger.Fatalf("Error processing arguments: %v.\n", err)
	}

	server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", configuration.Server.Host, configuration.Server.Port),
		Handler:           router,
		ReadTimeout:       time.Duration(configuration.Server.Timeouts.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(configuration.Server.Timeouts.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(configuration.Server.Timeouts.IdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(configuration.Server.Timeouts.ReadHeaderTimeout) * time.Second,
	}

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	authManager := manager.NewAuthManager(&dependencies)
	userManager := manager.NewUserManager(&dependencies)
	storageManager := manager.NewStorageManager(&dependencies)

	authManager.RegisterRoutes()
	userManager.RegisterRoutes()
	storageManager.RegisterRoutes()
}

// main starts the server, listens for termination signals, and gracefully shuts down the server upon receiving one.
// It uses a WaitGroup to ensure that the server runs concurrently and handles shutdown properly.
func main() {
	swg.Add(1)
	go func() {
		defer swg.Done()
		logger.Printf("Server is starting at %s:%d\n", configuration.Server.Host, configuration.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("error starting server: %v", err)
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		select {
		case sig := <-sigChan:
			logger.Printf("Received termination signal: %v. Initiating graceful shutdown...\n", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				errChan <- fmt.Errorf("error shutting down server gracefully: %v", err)
			} else {
				logger.Println("Server shut down gracefully")
			}

		case err := <-errChan:
			logger.Fatalf("Critical error occurred: %v\n", err)
		}
	}()

	swg.Wait()
}
