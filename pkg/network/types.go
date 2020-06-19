package network

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/vishvananda/netlink"
)

// Environment defines the configuration of the chest environment
type Environment struct {
	// Bridge configuration
	BridgeName    string `json:"bridgeName"`
	BridgeAddress string `json:"bridgeAddress"`

	// Used during runtime
	BridgeLink netlink.Link `json:"-"`

	// VM Nic configuration
	NicPrefix    string `json:"nicPrefix"`
	NicMacPrefix string `json:"nicMacPrefix"`
}

// OpenFile will open an file and parse the contents
func OpenFile(path string) (*Environment, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("Unable to find file [%s]", path)
	}

	var e Environment

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &e)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	return &e, nil
}

// ExampleConfig will return a config output
func ExampleConfig() string {
	cfg := Environment{}
	cfg.BridgeAddress = "192.168.1.1/24"
	cfg.BridgeName = "plunder"
	cfg.NicPrefix = "plunderVM"
	cfg.NicMacPrefix = "c0:ff:ee:"

	b, _ := yaml.Marshal(cfg)

	return string(b)
}
