package cmd

import (
	"fmt"
	"os"

	"github.com/plunder-app/shack/pkg/network"

	"github.com/spf13/cobra"
)

// Release - this struct contains the release information populated when building shack
var Release struct {
	Version string
	Build   string
}

func init() {

	// Main function commands
	shackCmd.AddCommand(shackExample)
	shackCmd.AddCommand(shackNetwork)
	shackCmd.AddCommand(shackVM)
	shackCmd.AddCommand(shackVersion)
}

//shackCmd is the parent command
var shackCmd = &cobra.Command{
	Use:   "shack",
	Short: "This is a tool for building a deployment environment",
}

// Execute - starts the command parsing process
func Execute() {
	if err := shackCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//// Sub commands

var shackVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the shack enviroment manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shack Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}

var shackNetwork = &cobra.Command{
	Use:   "network",
	Short: "Create the networking",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var shackExample = &cobra.Command{
	Use:   "example",
	Short: "Print example configuratiopn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(network.ExampleConfig())
	},
}
