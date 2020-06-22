package cmd

import (
	"fmt"
	"os"

	"github.com/plunder-app/chest/pkg/network"
	"github.com/plunder-app/chest/pkg/vmm"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var vmID string
var foreground bool

// Release - this struct contains the release information populated when building chest
var Release struct {
	Version string
	Build   string
}

func init() {

	chestVMStop.Flags().StringVar(&vmID, "id", "", "The UUID for a virtual machine")
	chestVMStart.Flags().BoolVarP(&foreground, "foreground", "f", false, "The UUID for a virtual machine")

	// Add subcommands
	chestVM.AddCommand(chestVMStart)
	chestVM.AddCommand(chestVMStop)

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

var chestVM = &cobra.Command{
	Use:   "vm",
	Short: "Create the networking",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chest VM configuration\n")
		cmd.Help()
	},
}

var chestVMStart = &cobra.Command{
	Use:   "start",
	Short: "Start a virtual Machine",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chest VM configuration\n")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}
		// Generate VM UUID
		b, err := vmm.GenVMUUID()
		if err != nil {
			log.Fatal(err)
		}
		uuid := fmt.Sprintf("%02x%02x%02x", b[0], b[1], b[2])
		vmInterface := fmt.Sprintf("%s-%s", cfg.NicPrefix, uuid)
		if len(vmInterface) > 15 {
			log.Fatalf("The interface name [%s] is too long for the interface standard, shorten the nicPrefix", vmInterface)
		}
		// Create Tap Device (and add to bridge)
		cfg.CreateTap(vmInterface)

		// Generate MAC address using UUID and Mac prefix
		mac := vmm.GenVMMac(cfg.NicMacPrefix, b)

		// Create Disk
		vmm.CreateDisk(uuid, "4G")
		// Start Virtual Machine
		vmm.Start(mac, uuid, cfg.NicPrefix, foreground)

		// If this is ran in the foreground then we will want to tidy up the created interface
		if foreground {
			log.Infof("Deleting interface [%s]", vmInterface)
			err = cfg.DeleteTap(vmInterface)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var chestVMStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop a virtual Machine",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}
		// Stop Virtual Machine
		vmm.Stop(vmID)
		// Remove Networking configuration
		err = cfg.DeleteTap(fmt.Sprintf("%s-%s", cfg.NicPrefix, vmID))
		if err != nil {
			log.Fatal(err)
		}
		// TODO - Delete disk?
	},
}

var chestExample = &cobra.Command{
	Use:   "example",
	Short: "Print example configuratiopn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(network.ExampleConfig())
	},
}
