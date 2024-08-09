package main

import (
	"GoShort/database"
	"GoShort/internal/router"
	"GoShort/logging"
	"GoShort/settings"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

var (
	currentLogger = logging.NewLogger("goshort-server")
)

var rootCmd = &cobra.Command{
	Use:   "goshort-server",
	Short: "GoShort! Server",
	Long:  `GoShort! is a pure CLI URL shortener and URL redirection service.`,
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the database models/schemas",
	Run: func(cmd *cobra.Command, args []string) {
		settings.LoadSettings()
		currentLogger.Info("Starting model/schema migration...")
		database.Connect()
		database.MigrateModels()
		currentLogger.Success("Model migration/schema completed!")
	},
}

var startCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve GoShort! 🚀",
	Run: func(cmd *cobra.Command, args []string) {
		currentLogger.Info("GoShort! HTTP Server is starting!")
		settings.LoadSettings()
		database.Connect()
		goShortRouter := router.SetupRouter()

		currentLogger.Info("GoShort! serving on [::1]:8000!")
		goShortHttpServer := &http.Server{
			Handler:      goShortRouter,
			Addr:         "[::1]:8000",
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		}
		log.Fatal(goShortHttpServer.ListenAndServe())
	},
}

func main() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}
