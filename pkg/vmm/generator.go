package vmm

import (
	"crypto/rand"
	"fmt"
)

// GenVMUUID creates the three byte UUID
func GenVMUUID() (buf []byte, err error) {

	buf = make([]byte, 3)

	_, err = rand.Read(buf)
	if err != nil {
		return
	}
	// TODO - should length be checked?

	return
}

// GenVMMac will create a mac address from the UUID and prefix
func GenVMMac(prefix string, buf []byte) string {
	return fmt.Sprintf("%s%02x:%02x:%02x", prefix, buf[0], buf[1], buf[2])
}
