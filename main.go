package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/g0disd3ad/rbt/internal/api"
	"github.com/g0disd3ad/rbt/internal/dictionary"
	"github.com/g0disd3ad/rbt/internal/storage"
)

func main() {
	rbtStorage := dictionary.NewRBTStorage()
	dict := dictionary.NewDictionary(rbtStorage)

	fmt.Println("~--- Initialising the dictionary ---~")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var pg *storage.PostgresStorage
	var dbErr error

	if dbPassword == "" {
		fmt.Println("Warning: DB_PASSWORD environment variable is missing. Database storage is disabled.")
	} else {
		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}
		if dbUser == "" {
			dbUser = "postgres"
		}
		if dbName == "" {
			dbName = "dict_db"
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
		pg, dbErr = storage.NewPostgresStorage(dsn)
		if dbErr != nil {
			fmt.Printf("Warning: Could not connect to DB: %v\n", dbErr)
		} else {
			fmt.Println("Successfully connected to PostgreSQL, loading data into tree...")
			if err := pg.LoadToTree(dict.Insert); err != nil {
				fmt.Printf("Error loading from database: %v\n", err)
			}
		}
	}

	apiServer := api.NewAPI(dict)

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Println("API Server is running on http://localhost:8080")
		serverErrors <- apiServer.StartServer(":8080")
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("Fatal API error: %v\n", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		fmt.Printf("Received signal: %v. Starting shutdown...\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := apiServer.Shutdown(ctx); err != nil {
			fmt.Printf("Shutdown failed: %v\n", err)
		}
	}

	if pg != nil {
		fmt.Println("Saving tree data to PostgreSQL before exit...")
		if err := pg.SaveFromTree(dict); err != nil {
			fmt.Printf("Error saving to database: %v\n", err)
		}
		if err := pg.Close(); err != nil {
			fmt.Printf("Error closing database: %v\n", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}

	fmt.Println("The program has finished working.")
}
