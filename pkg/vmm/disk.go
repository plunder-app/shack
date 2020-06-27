package vmm

import (
	"fmt"
	"os"
	"os/exec"
)

// CreateDisk - will create a disk for a qemu instance
func CreateDisk(uuid, size string) error {
	imagePath := fmt.Sprintf("%s.qcow2", uuid)

	// Check file stats
	_, err := os.Stat(imagePath)
	// If it doesn't exist then create it
	if os.IsNotExist(err) {
		if _, err := exec.Command("qemu-img", "create", "-f", "qcow2", imagePath, size).CombinedOutput(); err != nil {
			return err
		}
	}

	return nil
}

// DeleteDisk - will create a disk for a qemu instance
func DeleteDisk(uuid string) error {
	imagePath := fmt.Sprintf("%s.qcow2", uuid)

	// Check file stats
	_, err := os.Stat(imagePath)
	// If it doesn't exist then create it
	if os.IsNotExist(err) {
		return fmt.Errorf("The VM Disk [%s] does not exist", imagePath)
	}

	return os.Remove(imagePath)
}
