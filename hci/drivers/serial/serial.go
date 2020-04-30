package hcidriverserial

import (
	"errors"
	"io"
	"sync"
	"time"

	hciconst "github.com/BertoldVdb/go-ble/hci/const"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
)

var portOpenFuncs = make(map[string]HCISerialOpenFunc)

var (
	// ErrorSynchronisationLost is returned by the Run function when bad data is received,
	// and this cannot be recovered from (H4 mode)
	ErrorSynchronisationLost = errors.New("Synchronization lost")
)

// HCISerialOpenFunc is the signature of the function used to open a device name
type HCISerialOpenFunc func(deviceName string) (io.ReadWriteCloser, error)

type HCISerial struct {
	hciinterface.HCIInterface
	sync.Mutex

	port      io.ReadWriteCloser
	rxHandler hciinterface.HCIRxHandler
}

// ListDevices returns an array of strings containing all found devices.
func ListDevices() ([]string, error) {
	keys := make([]string, len(portOpenFuncs))

	i := 0
	for k := range portOpenFuncs {
		keys[i] = k
		i++
	}

	return keys, nil
}

// OpenPort returns the HCIInterface with a given io.ReadWriteCloser.
func OpenPort(port io.ReadWriteCloser) (hciinterface.HCIInterface, error) {
	return &HCISerial{
		port: port,
	}, nil
}

// Open returns the HCIInterface with a given device name.
func Open(deviceName string) (hciinterface.HCIInterface, error) {
	f, ok := portOpenFuncs[deviceName]

	if !ok {
		return nil, hciinterface.ErrorDeviceNotFound
	}

	port, err := f(deviceName)
	if err != nil {
		return nil, err
	}

	return OpenPort(port)
}

func getHCIPacketLength(buffer []byte) int {
	if len(buffer) < 1 {
		return 0
	}

	switch buffer[0] {
	case hciconst.MsgTypeCommand:
		fallthrough
	case hciconst.MsgTypeSCO:
		if len(buffer) >= 4 {
			return 4 + int(buffer[3])
		}
	case hciconst.MsgTypeACL:
		if len(buffer) >= 5 {
			return 5 + int(buffer[3]) | (int(buffer[4]) << 8)
		}
	case hciconst.MsgTypeEvent:
		if len(buffer) >= 3 {
			return 3 + int(buffer[2])
		}
	case hciconst.MsgTypeISO:
		if len(buffer) >= 5 {
			return 5 + int(buffer[3]) | (int(buffer[4]&0xF3) << 8)
		}
	default:
		return -1
	}

	return 0
}

// Run is the worker function. It needs to be running before packets can be sent or received.
func (d *HCISerial) Run() error {
	defer d.Close()

	rxBuf := make([]byte, 8192)
	rxBufIndex := 0

	for {
		n, err := d.port.Read(rxBuf[rxBufIndex:])
		if n == 0 || err != nil {
			return err
		}

		rxBufIndex += n
		if rxBufIndex >= len(rxBuf) {
			return ErrorSynchronisationLost
		}

		readIndex := 0
		for {
			workBuf := rxBuf[readIndex:rxBufIndex]

			pktLen := getHCIPacketLength(workBuf)
			if pktLen < 0 || pktLen >= len(rxBuf) {
				return ErrorSynchronisationLost
			}

			if pktLen > 0 && pktLen <= len(workBuf) {
				d.Lock()
				rxHandler := d.rxHandler
				d.Unlock()
				if rxHandler != nil {
					rxHandler(hciinterface.HCIRxPacket{
						Received:         true,
						Data:             workBuf[:pktLen],
						RxTime:           time.Now(),
						TimeFromHardware: false,
					})
				}

				readIndex += pktLen

			} else {
				if readIndex != 0 {
					copy(rxBuf, rxBuf[readIndex:rxBufIndex])
					rxBufIndex -= readIndex
				}
				break
			}
		}
	}
}

// Close closes the the interface. It can be called at any time and multiple times as well.
// It will terminate Run, if it was running.
func (d *HCISerial) Close() error {
	return d.port.Close()
}

// SendPacket sends a HCI packet to the device.
func (d *HCISerial) SendPacket(pkt hciinterface.HCITxPacket) error {
	_, err := d.port.Write(pkt.Data)
	return err
}

// SetRecvHandler configures the receive handler callback function. It will be called
// when a HCI packet is received.
func (d *HCISerial) SetRecvHandler(handler hciinterface.HCIRxHandler) error {
	d.Lock()
	defer d.Unlock()

	d.rxHandler = handler
	return nil
}

// RegisterDevice registers an open function for a device name.
func RegisterDevice(deviceName string, f HCISerialOpenFunc) {
	portOpenFuncs[deviceName] = f
}
