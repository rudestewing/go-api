package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "",
	Short: "A Go API boilerplate with Fiber framework",	Long: `A complete Go API boilerplate built with Fiber web framework.

This CLI provides commands for:
- Starting the HTTP API server (serve)
- Managing database migrations (migrate)
- Running database seeders (seed)

Examples:
  serve                     # Start the server
  migrate up                # Run migrations
  seed create posts         # Create seeder

Use the available subcommands to manage your API application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Commands are automatically registered via their respective init() functions
	// in migrate.go and seed.go files
}
