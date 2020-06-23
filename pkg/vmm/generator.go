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
func GenVMMac(prefix, s string) string {
	for i := 2; i < len(s); i += 3 {
		s = s[:i] + ":" + s[i:]
	}
	return fmt.Sprintf("%s%s", prefix, s)
}
