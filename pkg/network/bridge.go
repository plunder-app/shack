package network

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// CheckBridge will examine the bridge for its status
func (e *Environment) CheckBridge() error {

	if e.BridgeLink == nil {
		bridge, err := netlink.LinkByName(e.BridgeName)
		if err != nil {
			return err
		}
		e.BridgeLink = bridge
	}
	state := e.BridgeLink.Attrs()

	// Check Administrative state
	if net.FlagUp&state.Flags == 0 {
		return fmt.Errorf("Bridge exists, but is configured to a [down] state")
	}

	// Check the link state of bridge
	switch state.OperState {
	case netlink.OperDown:
		return fmt.Errorf("Bridge exists, but Link is physically in a [down] state")
	}

	return nil
}

//CreateBridge will create a new Layer 2 bridge, and configure it
func (e *Environment) CreateBridge() error {
	// Create the bridge
	mybridge := &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: e.BridgeName}}
	err := netlink.LinkAdd(mybridge)
	if err != nil {
		return fmt.Errorf("Could not add %s: %v", e.BridgeName, err)
	}
	return nil
}

// DeleteBridge will remove an existing bridge
func (e *Environment) DeleteBridge() error {
	if e.BridgeLink == nil {
		bridge, err := netlink.LinkByName(e.BridgeName)
		if err != nil {
			return err
		}
		e.BridgeLink = bridge
	}
	err := netlink.LinkDel(e.BridgeLink)
	if err != nil {
		return fmt.Errorf("Could not delete %s: %v", e.BridgeName, err)
	}

	// Remove and reference to this bridge
	e.BridgeLink = nil
	return nil
}

// AddBridgeAddress will add an address to an existing bridge
func (e *Environment) AddBridgeAddress() error {
	if e.BridgeLink == nil {
		bridge, err := netlink.LinkByName(e.BridgeName)
		if err != nil {
			return err
		}
		e.BridgeLink = bridge
	}

	addr, err := netlink.ParseAddr(e.BridgeAddress)
	if err != nil {
		return fmt.Errorf("Could not parse address %s: %v", e.BridgeAddress, err)
	}

	err = netlink.AddrAdd(e.BridgeLink, addr)
	if err != nil {
		return fmt.Errorf("Could not add address %s to bridge %s: %v", e.BridgeAddress, e.BridgeName, err)
	}

	return nil
}

// DelBridgeAddress will delete an address on an existing bridge
func (e *Environment) DelBridgeAddress() error {
	if e.BridgeLink == nil {
		bridge, err := netlink.LinkByName(e.BridgeName)
		if err != nil {
			return err
		}
		e.BridgeLink = bridge
	}

	addr, err := netlink.ParseAddr(e.BridgeAddress)
	if err != nil {
		return fmt.Errorf("Could not parse address %s: %v", e.BridgeAddress, err)
	}

	err = netlink.AddrDel(e.BridgeLink, addr)
	if err != nil {
		return fmt.Errorf("Could not del address %s on bridge %s: %v", e.BridgeAddress, e.BridgeName, err)
	}

	return nil
}

// BridgeUp sets the bridge to an enabled state
func (e *Environment) BridgeUp() error {
	if e.BridgeLink == nil {
		bridge, err := netlink.LinkByName(e.BridgeName)
		if err != nil {
			return err
		}
		e.BridgeLink = bridge
	}
	return netlink.LinkSetUp(e.BridgeLink)
}
