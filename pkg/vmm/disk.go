package vmm

import (
	"fmt"
	"os/exec"
)

// CreateDisk - will create a disk for a qemu instance
func CreateDisk(uuid, size string) error {
	imagePath := fmt.Sprintf("%s.qcow2", uuid)
	if _, err := exec.Command("qemu-img", "create", "-f", "qcow2", imagePath, size).CombinedOutput(); err != nil {
		return err
	}
	return nil
}
