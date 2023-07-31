package serviceserial

import (
	"context"
	"io"
	"sync"

	"github.com/BertoldVdb/go-ble/bleatt"
	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
)

func (s *SerialConfig) ClientFactory(conn io.ReadWriteCloser, cb func(conn io.ReadWriteCloser)) func(ctx context.Context, dev *bleatt.GattDevice) {
	return func(ctx context.Context, dev *bleatt.GattDevice) {
		defer conn.Close()
		if dev == nil {
			return
		}

		if cb != nil {
			defer cb(nil)
		}

		go func() {
			<-ctx.Done()
			conn.Close()
		}()

		structure := dev.ClientGetStructure(ctx)
		if structure == nil {
			return
		}

		service := structure.GetService(s.ServiceUUID)
		if service == nil {
			return
		}

		var charRd, charWr *attstructure.Characteristic
		if !s.ReadUUID.IsZero() {
			charRd = service.GetCharacteristic(s.ReadUUID)
		}
		if !s.WriteUUID.IsZero() {
			charWr = service.GetCharacteristic(s.WriteUUID)
		}

		/* If no characteristic UUID given, search for usable ones */
		for _, m := range service.GetCharacteristics() {
			flags := m.GetFlags()
			if charWr == nil && ((flags&attstructure.CharacteristicWriteNoAck > 0) || (flags&attstructure.CharacteristicWriteAck > 0)) {
				charWr = m
			}
			if charRd == nil && ((flags&attstructure.CharacteristicIndicate > 0) || (flags&attstructure.CharacteristicNotify > 0)) {
				charRd = m
			}
		}

		if charRd == nil || charWr == nil {
			return
		}

		if cb != nil {
			cb(conn)
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
