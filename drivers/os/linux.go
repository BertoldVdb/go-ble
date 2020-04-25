// +build linux

package oshci

import (
	"bytes"
	"encoding/binary"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"

	hciinterface "github.com/BertoldVdb/go-ble/drivers/interface"
	"golang.org/x/sys/unix"
)

const (
	cmdHCIDevUp   = 0x400448c9
	cmdHCIDevDown = 0x400448ca
	cmdHCIDevList = 0x800448d2
	cmdHCIDevInfo = 0x800448d3

	optHCIDataDir   = 1
	optHCIFilter    = 2
	optHCITimestamp = 3

	solHCI = 0

	eventCommandComplete = 0xE
	eventCommandStatus   = 0xE
	eventLeMeta          = 0x3E

	cmsgHCIDir       = 0x0001
	cmsgHCITimestamp = 0x0002
)

// HCILinux is a hciinterface.HCIInterface using the linux kernel as backend.
type HCILinux struct {
	hciinterface.HCIInterface
	sync.Mutex

	rxHandler hciinterface.HCIRxHandler
	closed    bool
	sock      int

	closePipe [2]int
}

func bluetoothSocket() (int, error) {
	return unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW, unix.BTPROTO_HCI)
}

func ioctl(fd int, method uintptr, value uintptr) (uintptr, error) {
	r, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), method, value)

	if err != 0 {
		return 0, os.NewSyscallError("SYS_IOCTL", err)
	}

	return r, nil
}

func testBit(input uint64, bit int) bool {
	value := uint64(1) << uint64(bit)
	return (input & value) > 0
}

func checkDevice(sock int, deviceID uint16) []byte {
	type deviceInfoStats struct {
		ErrRx  uint32
		ErrTx  uint32
		CmdTx  uint32
		EvtRx  uint32
		ACLTx  uint32
		ACLRx  uint32
		SCOTx  uint32
		SCORx  uint32
		ByteRx uint32
		ByteTx uint32
	}

	type devInfoStruct struct {
		id   uint16
		name [8]byte

		bdAddr [6]byte

		flags uint32
		type8 uint8

		features [8]byte

		pktType    uint32
		linkPolicy uint32
		linkMode   uint32

		aclMTU  uint16
		aclPkts uint16
		scoMTU  uint16
		scoPkts uint16

		stats deviceInfoStats
	}

	param := devInfoStruct{
		id: deviceID,
	}

	_, err := ioctl(sock, cmdHCIDevInfo, uintptr(unsafe.Pointer(&param)))
	if err != nil {
		return nil
	}

	var nameLen int
	for nameLen = 0; nameLen < len(param.name) && param.name[nameLen] != 0; nameLen++ {
	}

	return param.name[:nameLen]
}

type foundDevice struct {
	name     string
	deviceID int
	channel  uint16
}

func listDevicesInternal(sock int) ([]foundDevice, error) {
	type devReqStruct struct {
		devID  uint16
		devOpt uint32
	}

	type devListStruct struct {
		count uint16
		list  [32]devReqStruct
	}

	param := devListStruct{
		count: 32,
	}

	_, err := ioctl(sock, cmdHCIDevList, uintptr(unsafe.Pointer(&param)))
	if err != nil {
		return nil, err
	}

	var devices []foundDevice

	for i := uint16(0); i < param.count; i++ {
		if name := checkDevice(sock, param.list[i].devID); name != nil {
			devices = append(devices, foundDevice{
				name:     string(name),
				deviceID: int(param.list[i].devID),
				channel:  unix.HCI_CHANNEL_RAW,
			}, foundDevice{
				name:     string(name) + "u",
				deviceID: int(param.list[i].devID),
				channel:  unix.HCI_CHANNEL_USER,
			})
		}
	}

	return devices, nil
}

// ListDevices returns an array of strings containing all found devices.
func ListDevices() ([]string, error) {
	sock, err := bluetoothSocket()
	if err != nil {
		return nil, err
	}
	defer unix.Close(sock)

	devices, err := listDevicesInternal(sock)
	if err != nil {
		return nil, err
	}

	var strings = make([]string, len(devices))
	for i := range devices {
		strings[i] = devices[i].name
	}

	return strings, nil
}

// SetHCIFilter can be used to filter out certain packet types and events.
// By default, all data is passed through.
func (d *HCILinux) SetHCIFilter(types []uint8, events []uint8) error {
	type filterType struct {
		TypeMask  uint32
		EventMask [2]uint32
		Opcode    uint16
	}

	hciFilterValue := filterType{}

	if types != nil {
		for _, m := range types {
			hciFilterValue.TypeMask |= 1 << m
		}
	} else {
		hciFilterValue.TypeMask = 0xFFFFFFFF

	}

	if events != nil {
		for _, m := range events {
			if m >= 32 {
				hciFilterValue.EventMask[1] |= 1 << (m - 32)
			} else {
				hciFilterValue.EventMask[0] |= 1 << m
			}
		}
	} else {
		a := ^uint32(0)
		hciFilterValue.EventMask[0] = a
		hciFilterValue.EventMask[1] = a
	}

	/* Dirtier code is not possible */
	const bytes = int(unsafe.Sizeof(filterType{}))
	filterFakeString := string((*(*[bytes]byte)(unsafe.Pointer(&hciFilterValue)))[:])
	defer runtime.KeepAlive(filterFakeString)

	return unix.SetsockoptString(d.sock, solHCI, optHCIFilter, filterFakeString)
}

// Open creates the the hciinterface.HCIInterface bound to a specified device
func Open(deviceName string) (hciinterface.HCIInterface, error) {
	sock, err := bluetoothSocket()
	if err != nil {
		return nil, err
	}

	closeErr := func(err error) (*HCILinux, error) {
		unix.Close(sock)
		return nil, err
	}

	devices, err := listDevicesInternal(sock)
	if err != nil {
		return closeErr(err)
	}

	var device *foundDevice
	for i, m := range devices {
		if m.name == deviceName {
			device = &devices[i]
		}
	}
	if device == nil {
		return closeErr(hciinterface.ErrorDeviceNotFound)
	}

	err = unix.SetsockoptInt(sock, solHCI, optHCIDataDir, 1)
	if err != nil {
		return closeErr(err)
	}

	err = unix.SetsockoptInt(sock, solHCI, optHCITimestamp, 1)
	if err != nil {
		return closeErr(err)
	}

	/* Turn off the device to use it with a user channel, and turn it on for a raw channel */
	if device.channel == unix.HCI_CHANNEL_USER {
		err = unix.IoctlSetInt(sock, cmdHCIDevDown, device.deviceID)
		if err != nil {
			return closeErr(err)
		}
	} else {
		/* Can fail, is harmless */
		unix.IoctlSetInt(sock, cmdHCIDevUp, device.deviceID)
	}

	/* Bind to the device */
	hciAddr := unix.SockaddrHCI{Dev: uint16(device.deviceID), Channel: device.channel}
	err = unix.Bind(sock, &hciAddr)
	if err != nil {
		return closeErr(err)
	}

	hci := &HCILinux{
		sock:      sock,
		closePipe: [2]int{-1, -1},
	}

	/* Message filtering is not possible for user channels due to bypassing the code that handles it */
	if device.channel == unix.HCI_CHANNEL_RAW {
		err = hci.SetHCIFilter(nil, nil)
		if err != nil {
			return closeErr(err)
		}
	}

	return hci, nil
}

// Run is the worker function. It needs to be running before packets can be sent or received.
func (d *HCILinux) Run() error {
	d.Lock()
	if d.closed {
		return ErrorClosed
		d.Unlock()
	}

	err := unix.Pipe(d.closePipe[:])
	if err != nil {
		unix.Close(d.sock)
		d.closed = true
		return err
	}
	d.Unlock()

	defer func() {
		d.Lock()
		d.closed = true
		unix.Close(d.sock)
		unix.Close(d.closePipe[0])
		unix.Close(d.closePipe[1])
		d.Unlock()
	}()

	pfd := []unix.PollFd{{
		Fd:     int32(d.sock),
		Events: unix.POLLIN,
	}, {
		Fd:     int32(d.closePipe[0]),
		Events: unix.POLLIN,
	}}

	bufData := make([]byte, 2048)
	bufOOB := make([]byte, 256)

	for {
		d.Lock()
		if d.closed {
			d.Unlock()
			return ErrorClosed
		}
		d.Unlock()

		_, err := unix.Poll(pfd, -1)
		if err != nil {
			return err
		}

		if pfd[0].Revents&(unix.POLLHUP|unix.POLLERR) > 0 {
			return ErrorHUP
		}

		if pfd[0].Revents&unix.POLLIN > 0 {
			n, oobn, _, _, err := unix.Recvmsg(d.sock, bufData, bufOOB, 0)
			if err != nil {
				return err
			}
			if n == 0 {
				continue
			}

			msgs, err := unix.ParseSocketControlMessage(bufOOB[:oobn])
			if err != nil {
				return err
			}

			pkt := hciinterface.HCIRxPacket{
				Data:     bufData[:n],
				Received: true,
			}

			for _, msg := range msgs {
				switch msg.Header.Type {
				case cmsgHCIDir:
					if len(msg.Data) > 0 {
						pkt.Received = msg.Data[0] == 1
					}

				case cmsgHCITimestamp:
					var timeval unix.Timeval
					if binary.Read(bytes.NewBuffer(msg.Data), binary.LittleEndian, &timeval) == nil {
						pkt.RxTime = time.Unix(timeval.Unix())
						pkt.TimeFromHardware = true
					}
				}
			}

			if !pkt.TimeFromHardware {
				/* Oh well, most applications don't need it */
				pkt.RxTime = time.Now()
			}

			d.Lock()
			rxHandler := d.rxHandler
			d.Unlock()

			if rxHandler != nil {
				err = rxHandler(pkt)
				if err != nil {
					return err
				}
			}
		}
	}
}

// Close closes the the interface. It can be called at any time and multiple times as well.
// It will terminate Run, if it was running.
func (d *HCILinux) Close() error {
	d.Lock()
	defer d.Unlock()

	if d.closed {
		return nil
	}
	d.closed = true

	if d.closePipe[1] >= 0 {
		_, err := unix.Write(d.closePipe[1], []byte{'c'})
		return err
	}

	return unix.Close(d.sock)
}

// SendPacket sends a HCI packet to the device.
func (d *HCILinux) SendPacket(pkt hciinterface.HCITxPacket) error {
	d.Lock()
	defer d.Unlock()

	if d.closed {
		return ErrorClosed
	}
	_, err := unix.Write(d.sock, pkt.Data)
	return err
}

// SetRecvHandler configures the receive handler callback function. It will be called
// when a HCI packet is received.
func (d *HCILinux) SetRecvHandler(handler hciinterface.HCIRxHandler) error {
	d.Lock()
	defer d.Unlock()

	d.rxHandler = handler
	return nil
}
