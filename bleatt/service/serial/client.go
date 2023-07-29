package serviceserial

import (
	"context"
	"io"
	"sync"

	"github.com/BertoldVdb/go-ble/bleatt"
)

func (s *SerialConfig) ClientFactory(conn io.ReadWriteCloser) func(ctx context.Context, dev *bleatt.GattDevice) {
	readUUID := s.ReadUUID
	if readUUID.IsZero() {
		readUUID = s.ServiceUUID.CreateVariantAlt(1)
	}

	writeUUID := s.WriteUUID
	if writeUUID.IsZero() {
		writeUUID = s.ServiceUUID.CreateVariantAlt(2)
	}

	return func(ctx context.Context, dev *bleatt.GattDevice) {
		defer conn.Close()
		if dev == nil {
			return
		}

		go func() {
			<-ctx.Done()
			conn.Close()
		}()

		structure := dev.ClientGetStructure(ctx)
		if structure == nil {
			return
		}

		serviceJbd := structure.GetService(s.ServiceUUID)
		if serviceJbd == nil {
			return
		}

		charRd := serviceJbd.GetCharacteristic(readUUID)
		charWr := serviceJbd.GetCharacteristic(writeUUID)
		if charRd == nil || charWr == nil {
			return
		}

		var txMtx sync.Mutex

		if err := charRd.Subscribe(ctx, func(value []byte) {
			txMtx.Lock()
			conn.Write(value)
			txMtx.Unlock()
		}); err != nil {
			return
		}

		var rxBuf [1024]byte
		for {
			n, err := conn.Read(rxBuf[:])
			if err != nil {
				return
			}
			in := rxBuf[:n]

			index := 0
			for index < len(in) {
				bytes, err := charWr.SetValue(ctx, in[index:])
				if bytes < 0 || err != nil {
					return
				}
				if bytes == 0 {
					break
				}
				index += bytes
			}
		}
	}
}
