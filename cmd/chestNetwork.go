package cmd

import (
	"fmt"

	"github.com/plunder-app/chest/pkg/network"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configPath string

func init() {

	chestNetwork.PersistentFlags().StringVarP(&configPath, "config", "c", "chest.yaml", "The path to the chest environment configuration")

	chestNetwork.AddCommand(chestNetworkCreate)
	chestNetwork.AddCommand(chestNetworkCheck)
	chestNetwork.AddCommand(chestNetworkDelete)
	chestNetwork.AddCommand(chestNetworkNat)
}

var chestNetworkCheck = &cobra.Command{
	Use:   "check",
	Short: "check the bridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chest Networking configuration\n")

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

var chestNetworkCreate = &cobra.Command{
	Use:   "create",
	Short: "Create the bridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chest Networking configuration\n")

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

var chestNetworkDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete the bridge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chest Networking configuration\n")
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

var chestNetworkNat = &cobra.Command{
	Use:   "nat",
	Short: "Enable Nat",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chest Networking configuration\n")
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
