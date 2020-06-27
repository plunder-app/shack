package cmd

import (
	"fmt"
	"os"

	"github.com/plunder-app/chest/pkg/network"

	"github.com/spf13/cobra"
)

// Release - this struct contains the release information populated when building chest
var Release struct {
	Version string
	Build   string
}

func init() {

	// Main function commands
	chestCmd.AddCommand(chestExample)
	chestCmd.AddCommand(chestNetwork)
	chestCmd.AddCommand(chestVM)
	chestCmd.AddCommand(chestVersion)
}

//chestCmd is the parent command
var chestCmd = &cobra.Command{
	Use:   "chest",
	Short: "This is a tool for building a deployment environment",
}

// Execute - starts the command parsing process
func Execute() {
	if err := chestCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//// Sub commands

var chestVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the Chest enviroment manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chest Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}

var chestNetwork = &cobra.Command{
	Use:   "network",
	Short: "Create the networking",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var chestExample = &cobra.Command{
	Use:   "example",
	Short: "Print example configuratiopn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(network.ExampleConfig())
	},
}
