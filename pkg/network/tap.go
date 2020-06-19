package network

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

//CreateTap  will create a tap device for qemu
func (e *Environment) CreateTap(tapName string) error {
	// Find bridge
	if e.BridgeLink == nil {
		bridge, err := netlink.LinkByName(e.BridgeName)
		if err != nil {
			return err
		}
		e.BridgeLink = bridge
	}

	// Create TAP
	tap := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: tapName}, Mode: netlink.TUNTAP_MODE_TAP}
	err := netlink.LinkAdd(tap)
	if err != nil {
		return fmt.Errorf("Could not add %s: %v", tap.Name, err)
	}

	// Add Tap to bridge
	err = netlink.LinkSetMaster(tap, e.BridgeLink)
	if err != nil {
		return fmt.Errorf("Could not add %s to bridge %s : %v", tap.Name, e.BridgeName, err)
	}
	return nil
}

//DeleteTap  will remove a tap device from qemu and the bridge
func (e *Environment) DeleteTap(tapName string) error {

	// Find Tap
	tapLink, err := netlink.LinkByName(tapName)
	if err != nil {
		return err
	}

	// Remove Tap device
	err = netlink.LinkDel(tapLink)
	if err != nil {
		return fmt.Errorf("Could not delete %s: %v", e.BridgeName, err)
	}
	return nil

}
