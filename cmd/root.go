package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "go-api",
	Short: "A Go API boilerplate with Fiber framework",
	Long: `A complete Go API boilerplate built with Fiber web framework.

This CLI provides commands for:
- Starting the HTTP API server
- Managing database migrations
- Running database seeders

Use the available subcommands to manage your API application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
