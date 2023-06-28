package btsnoop

import (
	"encoding/binary"
	"io"
	"os"
	"sync"
	"time"

	hciconst "github.com/BertoldVdb/go-ble/hci/const"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
)

type logger struct {
	sync.Mutex
	hciinterface.HCIInterface
	out io.WriteCloser
}

func WrapFile(dev hciinterface.HCIInterface, path string) (hciinterface.HCIInterface, error) {
	file, err := os.Create(path)
	if err != nil {
		return dev, err
	}

	w, err := Wrap(dev, file)
	if err != nil {
		file.Close()
	}

	return w, err
}

func Wrap(dev hciinterface.HCIInterface, out io.WriteCloser) (hciinterface.HCIInterface, error) {
	/* Write header */
	if _, err := out.Write([]byte("btsnoop\x00")); err != nil {
		return dev, err
	} else if _, err := out.Write(binary.BigEndian.AppendUint32(nil, 1)); err != nil {
		return dev, err
	} else if _, err := out.Write(binary.BigEndian.AppendUint32(nil, 1002)); err != nil {
		return dev, err
	}

	return &logger{
		HCIInterface: dev,
		out:          out,
	}, nil
}

var refTime = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

func (l *logger) logPacket(isTransmit bool, ts time.Time, data []byte) {
	if len(data) < 1 {
		return
	}

	entry := make([]byte, 24+len(data))
	copy(entry[24:], data)

	flags := uint32(0)
	if !isTransmit {
		flags |= 1
	}
	if data[0] == hciconst.MsgTypeCommand || data[0] == hciconst.MsgTypeEvent { /* Command or Event */
		flags |= 2
	}

	interval := ts.Sub(refTime).Microseconds() + 0x00E03AB44A676000

	binary.BigEndian.PutUint32(entry, uint32(len(data)))
	binary.BigEndian.PutUint32(entry[4:], uint32(len(data)))
	binary.BigEndian.PutUint32(entry[8:], flags)
	binary.BigEndian.PutUint64(entry[16:], uint64(interval))

	l.Lock()
	l.out.Write(entry)
	l.Unlock()
}

func (l *logger) SendPacket(pkt hciinterface.HCITxPacket) error {
	l.logPacket(true, time.Now(), pkt.Data)
	return l.HCIInterface.SendPacket(pkt)
}

func (l *logger) SetRecvHandler(cb hciinterface.HCIRxHandler) error {
	return l.HCIInterface.SetRecvHandler(func(pkt hciinterface.HCIRxPacket) error {
		l.logPacket(false, pkt.RxTime, pkt.Data)
		return cb(pkt)
	})
}

func (l *logger) Close() error {
	l.out.Close()
	return l.HCIInterface.Close()
}
