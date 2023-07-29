package serviceserial

import (
	"context"
	"errors"
	"io"
	"sync"

	attperipheral "github.com/BertoldVdb/go-ble/bleatt/helpers/peripheral"
	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

type SerialConfig struct {
	ServiceUUID bleutil.UUID
	ReadUUID    bleutil.UUID
	WriteUUID   bleutil.UUID
	Connect     func() (io.ReadWriteCloser, error)
	Secure      bool
}

func DefaultConfig() *SerialConfig {
	return &SerialConfig{
		ServiceUUID: bleutil.UUIDFromStringPanic("6e400001-b5a3-f393-e0a9-e50e24dcca9e"),
		Connect: func() (io.ReadWriteCloser, error) {
			return nil, errors.New("Connect method not specified")
		},
	}
}

type SerialNordic struct {
	config *SerialConfig

	connMutex sync.Mutex
	conn      io.ReadWriteCloser

	dataTx *attstructure.Characteristic
}

func (s *SerialNordic) CreateStructure(structure *attstructure.Structure) error {
	secure := attstructure.CharacteristicNeedsEncryption
	if !s.config.Secure {
		secure = 0
	}

	readUUID := s.config.ReadUUID
	if readUUID.IsZero() {
		readUUID = s.config.ServiceUUID.CreateVariantAlt(1)
	}

	writeUUID := s.config.WriteUUID
	if writeUUID.IsZero() {
		writeUUID = s.config.ServiceUUID.CreateVariantAlt(2)
	}

	pspp := structure.AddPrimaryService(s.config.ServiceUUID)
	pspp.AddCharacteristic(readUUID, attstructure.CharacteristicWriteAck|attstructure.CharacteristicWriteNoAck|secure, attstructure.ValueConfig{
		ValueWriteCb: func(h *attstructure.GATTHandle) error {
			s.connMutex.Lock()
			defer s.connMutex.Unlock()

			_, err := s.conn.Write(h.Value)
			return err
		},
	})
	s.dataTx = pspp.AddCharacteristic(writeUUID, attstructure.CharacteristicRead|attstructure.CharacteristicNotify|secure, attstructure.ValueConfig{})
	return nil
}

func (s *SerialNordic) Disconnected() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *SerialNordic) Connected(conn hciconnmgr.BufferConn) error {
	var err error
	s.conn, err = s.config.Connect()
	if err != nil {
		return err
	}

	var rxBuf [1024]byte
	ctx := context.Background()
	go func() {
		defer conn.Close()
		defer s.conn.Close()

		for {
			n, err := s.conn.Read(rxBuf[:])
			if err != nil {
				return
			}
			in := rxBuf[:n]

			index := 0
			for index < len(in) {
				bytes, err := s.dataTx.SetValue(ctx, in[index:])
				if bytes < 0 || err != nil {
					return
				}
				if bytes == 0 {
					break
				}
				index += bytes
			}
		}
	}()

	return nil
}

func CreateService(config *SerialConfig) func() attperipheral.PeripheralImplementation {
	return func() attperipheral.PeripheralImplementation {
		return &SerialNordic{
			config: config,
		}
	}
}
