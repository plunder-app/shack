package cmd

import (
	"fmt"

	"github.com/plunder-app/shack/pkg/network"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configPath string

func init() {

	shackNetwork.PersistentFlags().StringVarP(&configPath, "config", "c", "shack.yaml", "The path to the shack environment configuration")

	shackNetwork.AddCommand(shackNetworkCreate)
	shackNetwork.AddCommand(shackNetworkCheck)
	shackNetwork.AddCommand(shackNetworkDelete)
	shackNetwork.AddCommand(shackNetworkNat)
}

var shackNetworkCheck = &cobra.Command{
	Use:   "check",
	Short: "check the bridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shack Networking configuration\n")

		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		err = cfg.CheckBridge()
		if err != nil {
			log.Warn(err)
		}

	},
}

var shackNetworkCreate = &cobra.Command{
	Use:   "create",
	Short: "Create the bridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shack Networking configuration\n")

		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		err = cfg.CreateBridge()
		if err != nil {
			log.Warn(err)
		}

		err = cfg.AddBridgeAddress()
		if err != nil {
			log.Warn(err)
		}

		err = cfg.BridgeUp()
		if err != nil {
			log.Warn(err)
		}
	},
}

var shackNetworkDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete the bridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shack Networking configuration\n")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		err = cfg.DeleteBridge()
		if err != nil {
			log.Warn(err)
		}

	},
}

var shackNetworkNat = &cobra.Command{
	Use:   "nat",
	Short: "Enable Nat",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shack Networking configuration\n")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		err = cfg.EnableNat()
		if err != nil {
			log.Warn(err)
		}

	},
}
