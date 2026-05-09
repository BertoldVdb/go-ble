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
	out    io.WriteCloser
	failed bool
}

// WrapFile opens path with mode 0600 (snoop logs may contain link keys and
// other sensitive material — never world-readable) and wraps dev so all
// HCI traffic is recorded in btsnoop format.
func WrapFile(dev hciinterface.HCIInterface, path string) (hciinterface.HCIInterface, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return dev, err
	}

	w, err := Wrap(dev, file)
	if err != nil {
		file.Close()
		return dev, err
	}

	return w, nil
}

func Wrap(dev hciinterface.HCIInterface, out io.WriteCloser) (hciinterface.HCIInterface, error) {
	/* Write header. Use a single contiguous buffer + writeAll so a partial
	   write or write error is detected and we don't ship a half-formed file
	   header (the format is record-aligned: any short write below the file
	   header desynchronizes the entire log). */
	header := make([]byte, 0, 16)
	header = append(header, "btsnoop\x00"...)
	header = binary.BigEndian.AppendUint32(header, 1)    // version
	header = binary.BigEndian.AppendUint32(header, 1002) // datalink type: HCI UART (H4)

	if err := writeAll(out, header); err != nil {
		return dev, err
	}

	return &logger{
		HCIInterface: dev,
		out:          out,
	}, nil
}

// writeAll keeps writing until the buffer is fully consumed or a non-EAGAIN
// error is returned. io.Writer permits short writes; the btsnoop format
// can't tolerate them.
func writeAll(w io.Writer, buf []byte) error {
	for len(buf) > 0 {
		n, err := w.Write(buf)
		if n > 0 {
			buf = buf[n:]
		}
		if err != nil {
			if n == 0 {
				return err
			}
			// short write with error — keep retrying until n==0 returns
			continue
		}
	}
	return nil
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
	if !l.failed {
		if err := writeAll(l.out, entry); err != nil {
			// Stop writing once the file is broken — partial records render
			// every following record unparseable.
			l.failed = true
		}
	}
	l.Unlock()
}

func (l *logger) SendPacket(pkt hciinterface.HCITxPacket) error {
	l.logPacket(true, time.Now(), pkt.Data)
	return l.HCIInterface.SendPacket(pkt)
}

func (l *logger) SetRecvHandler(cb hciinterface.HCIRxHandler) error {
	return l.HCIInterface.SetRecvHandler(func(pkt hciinterface.HCIRxPacket) error {
		l.logPacket(false, pkt.RxTime, pkt.Data)
		if cb == nil {
			return nil
		}
		return cb(pkt)
	})
}

func (l *logger) Close() error {
	l.out.Close()
	return l.HCIInterface.Close()
}
