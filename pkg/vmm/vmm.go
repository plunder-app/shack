package vmm

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"

	"context"

	"github.com/intel/govmm/qemu"
)

// Machine data
//  - UUID xx
//  - MACAddress de:ad:be:ef:x:x
//  - State directory ~/qmp
//  - qmp-socket /tmp/qmp-xx

// VMM - Virtual Machine Manager

//Start w
func Start(mac, uuid, nicPrefix string, foreground, vnc bool) error {
	var vncPort int
	params := make([]string, 0, 32)

	// Rootfs
	params = append(params, "-drive", fmt.Sprintf("file=%s.qcow2,if=virtio,aio=threads,format=qcow2", uuid))

	// Network
	net := fmt.Sprintf("tap,model=virtio-net-pci,mac=%s,ifname=%s-%s", mac, nicPrefix, uuid)
	params = append(params, "-nic", net)

	// kvm
	params = append(params, "-enable-kvm", "-cpu", "host")

	// resources
	params = append(params, "-m", "1024", "-display", "none")

	if foreground {
		params = append(params, "-curses")
		return foreGroundRunner(params)
	}

	if vnc {
		vncPort = vncPortGenerator()
		params = append(params, "-vnc", fmt.Sprintf(":%d", vncPort))
	}

	params = append(params, "-daemonize", "-qmp", fmt.Sprintf("unix:/tmp/qmp-%s,server,nowait", uuid))

	// LaunchCustomQemu should return as soon as the instance has launched as we
	// are using the --daemonize flag.  It will set up a unix domain socket
	// in /tmp/qmp-uuid that we can use to manage the instance.
	details, err := qemu.LaunchCustomQemu(context.Background(), "", params, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("%v : %v", details, err)
	}

	fmt.Printf("Network Device:\t%s-%s\nVM MAC Address:\t%s\nVM UUID:\t%s\n", nicPrefix, uuid, mac, uuid)

	if vnc {
		// Add the base VNC port (5900)
		fmt.Printf("VNC Port:\t%d\n", (vncPort + 5900))
	}
	return nil
}

// Stop will allow us to stop a VM
func Stop(uuid string) error {
	// This channel will be closed when the instance dies.
	disconnectedCh := make(chan struct{})

	// Set up our options.  We don't want any logging or to receive any events.
	cfg := qemu.QMPConfig{}
	// Start monitoring the qemu instance.  This functon will block until we have
	// connect to the QMP socket and received the welcome message.
	q, _, err := qemu.QMPStart(context.Background(), fmt.Sprintf("/tmp/qmp-%s", uuid), cfg, disconnectedCh)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// This has to be the first command executed in a QMP session.
	err = q.ExecuteQMPCapabilities(context.Background())
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// Let's try to shutdown the VM.  If it hasn't shutdown in 10 seconds we'll
	// send a quit message.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = q.ExecuteSystemPowerdown(ctx)
	cancel()
	if err != nil {
		err = q.ExecuteQuit(context.Background())
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	q.Shutdown()
	err = os.Remove(fmt.Sprintf("/tmp/qmp-%s", uuid))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func foreGroundRunner(params []string) error {

	path := "qemu-system-x86_64"

	/* #nosec */
	cmd := exec.Command(path, params...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Qemu error [%v]", err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Shell error [%v]", err)
	}
	return nil
}

func vncPortGenerator() int {
	rand.Seed(time.Now().UnixNano())
	port := 1 + rand.Intn(10000-1)

	return port
}
