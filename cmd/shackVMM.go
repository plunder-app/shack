package cmd

import (
	"fmt"

	"github.com/plunder-app/shack/pkg/network"
	"github.com/plunder-app/shack/pkg/vmm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var vmUUID string
var foreground, vnc, disk bool

func init() {

	shackVM.PersistentFlags().StringVar(&vmUUID, "id", "000000", "The UUID for a virtual machine")
	shackVMStart.Flags().BoolVarP(&foreground, "foreground", "f", false, "The UUID for a virtual machine")
	shackVMStart.Flags().BoolVarP(&vnc, "vnc", "v", false, "Enable VNC")
	shackVMStop.Flags().BoolVarP(&disk, "disk", "d", false, "Delete Disk")

	// Add subcommands
	shackVM.AddCommand(shackVMStart)
	shackVM.AddCommand(shackVMStop)
}

var shackVM = &cobra.Command{
	Use:   "vm",
	Short: "Create the networking",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shack VM configuration\n")
		cmd.Help()
	},
}

var shackVMStart = &cobra.Command{
	Use:   "start",
	Short: "Start a virtual Machine",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shack VM configuration\n")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}
		if vmUUID == "000000" {
			// Generate VM UUID
			b, err := vmm.GenVMUUID()
			if err != nil {
				log.Fatal(err)
			}
			vmUUID = fmt.Sprintf("%02x%02x%02x", b[0], b[1], b[2])
		}

		vmInterface := fmt.Sprintf("%s-%s", cfg.NicPrefix, vmUUID)
		if len(vmInterface) > 15 {
			log.Fatalf("The interface name [%s] is too long for the interface standard, shorten the nicPrefix", vmInterface)
		}
		// Create Tap Device (and add to bridge)
		cfg.CreateTap(vmInterface)

		// Generate MAC address using UUID and Mac prefix
		mac := vmm.GenVMMac(cfg.NicMacPrefix, vmUUID)

		// Create Disk
		vmm.CreateDisk(vmUUID, "4G")

		// Start Virtual Machine
		vmm.Start(mac, vmUUID, cfg.NicPrefix, foreground, vnc)

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

var shackVMStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop a virtual Machine",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		// Stop Virtual Machine
		err = vmm.Stop(vmUUID)
		if err != nil {
			log.Fatal(err)
		}

		// Remove Networking configuration
		err = cfg.DeleteTap(fmt.Sprintf("%s-%s", cfg.NicPrefix, vmUUID))
		if err != nil {
			log.Fatal(err)
		}
		// Delete disk
		if disk {
			err = vmm.DeleteDisk(vmUUID)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
