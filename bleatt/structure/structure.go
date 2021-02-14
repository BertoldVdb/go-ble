package attstructure

import (
	"context"
	"errors"
	"fmt"
	"sync"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

type ClientReadHandler func(ctx context.Context, handle uint16, buf []byte) ([]byte, error)
type ClientWriteHandler func(ctx context.Context, handle uint16, buf []byte, withRsp bool) (int, error)
type ClientNotifyHandler func(value []byte)

type Structure struct {
	isClient    bool
	clientRead  ClientReadHandler
	clientWrite ClientWriteHandler

	clientNotifyMutex sync.Mutex
	clientNotifyMap   map[uint16](ClientNotifyHandler)

	services []*Service
	exported *ExportedStructure
}

type Service struct {
	isPrimary       bool
	parent          *Structure
	uuid            bleutil.UUID
	characteristics []*Characteristic
}

type Characteristic struct {
	parent       *Service
	uuid         bleutil.UUID
	flags        CharacteristicFlag
	initialValue []byte

	valueConfig ValueConfig

	ValueHandle *GATTHandle
}

type ValueConfig struct {
	ValueBeforeReadCb func(h *GATTHandle, offset int) error
	ValueAfterReadCb  func(h *GATTHandle, offset int, bytes int) error
	ValueWriteCb      func(h *GATTHandle) error
	LengthFixed       bool
	LengthMax         uint16
}

func NewStructure() *Structure {
	return &Structure{}
}

func (s *Structure) AddPrimaryService(uuid bleutil.UUID) *Service {
	p := &Service{
		parent:    s,
		uuid:      uuid,
		isPrimary: true,
	}
	s.services = append(s.services, p)
	return p
}

func (s *Structure) GetServices() []*Service {
	return s.services
}

func (s *Structure) GetService(uuid bleutil.UUID) *Service {
	for _, m := range s.services {
		if m.uuid == uuid {
			return m
		}
	}
	return nil
}

func (p *Service) AddCharacteristic(uuid bleutil.UUID, flags CharacteristicFlag, valueConfig ValueConfig) *Characteristic {
	c := &Characteristic{
		parent: p,
		uuid:   uuid,
		flags:  flags,

		valueConfig: valueConfig,
	}

	p.characteristics = append(p.characteristics, c)

	return c
}

func (p *Service) AddCharacteristicReadOnly(uuid bleutil.UUID, value []byte) *Characteristic {
	c := p.AddCharacteristic(uuid, CharacteristicRead, ValueConfig{})
	c.initialValue = value
	return c
}

func (p *Service) GetCharacteristics() []*Characteristic {
	return p.characteristics
}

func (p *Service) GetCharacteristic(uuid bleutil.UUID) *Characteristic {
	for _, m := range p.characteristics {
		if m.uuid == uuid {
			return m
		}
	}
	return nil
}

type HandleInfo struct {
	Handle         uint16
	GroupEndHandle uint16
	UUIDWidth      int
	UUID           bleutil.UUID

	Flags CharacteristicFlag
}

type GATTHandle struct {
	Info HandleInfo

	Value       []byte
	ValueConfig ValueConfig

	CCCHandle *GATTHandle
}

func (c *Characteristic) SetValue(ctx context.Context, new []byte) (int, error) {
	if c.parent.parent.isClient {
		useAck := c.flags&CharacteristicWriteAck > 0

		return c.parent.parent.clientWrite(ctx, c.ValueHandle.Info.Handle, new, useAck)
	}

	e := c.parent.parent.exported
	e.Lock()

	if c.valueConfig.LengthFixed {
		copy(c.ValueHandle.Value, new)
	} else {
		c.ValueHandle.Value = append(c.ValueHandle.Value[:0], new...)
	}

	e.Unlock()

	return e.HandleSet(c, new)
}

func (c *Characteristic) GetValue(ctx context.Context, buf []byte) ([]byte, error) {
	if c.parent.parent.isClient {
		return c.parent.parent.clientRead(ctx, c.ValueHandle.Info.Handle, buf)
	}

	e := c.parent.parent.exported
	e.Lock()
	defer e.Unlock()

	return append(buf[:0], c.ValueHandle.Value...), nil
}

func (c *Characteristic) GetFlags() CharacteristicFlag {
	return c.flags
}

func (c *Characteristic) Subscribe(ctx context.Context, handler ClientNotifyHandler) error {
	if !c.parent.parent.isClient {
		return errors.New("Invalid mode")
	}

	hasIndicate := c.flags&CharacteristicIndicate > 0
	hasNotify := c.flags&CharacteristicNotify > 0

	if (!hasIndicate && !hasNotify) || c.ValueHandle.CCCHandle == nil {
		return errors.New("This characteristic does not support subscribing")
	}

	new := []byte{0}
	handle := c.ValueHandle.Info.Handle

	c.parent.parent.clientNotifyMutex.Lock()
	if handler != nil {
		if hasIndicate {
			new[0] = 1 << 1
		} else if hasNotify {
			new[0] = 1 << 0
		}
		c.parent.parent.clientNotifyMap[handle] = handler
	} else {
		delete(c.parent.parent.clientNotifyMap, handle)
	}
	defer c.parent.parent.clientNotifyMutex.Unlock()

	_, err := c.parent.parent.clientWrite(ctx, c.ValueHandle.CCCHandle.Info.Handle, new, true)
	return err
}

func (s Structure) String() string {
	result := ""
	ln := func(format string, a ...interface{}) {
		newLine := fmt.Sprintln()
		result += fmt.Sprintf(format+newLine, a...)
	}

	for _, m := range s.services {
		ln("%v:", m.uuid)
		for _, k := range m.characteristics {
			ln(" %v (%02x):", k.uuid, k.flags)
			if k.ValueHandle != nil {
				ln("  ValueHandle: %04x", k.ValueHandle.Info.Handle)
				if k.ValueHandle.CCCHandle != nil {
					ln("   CCCHandle: %04x", k.ValueHandle.CCCHandle.Info.Handle)
				}
			}
		}
	}

	return result
}

func (s *Structure) InjectNotify(handle uint16, data []byte) {
	if s.clientNotifyMap == nil {
		return
	}

	s.clientNotifyMutex.Lock()
	f, ok := s.clientNotifyMap[handle]
	s.clientNotifyMutex.Unlock()
	if ok && f != nil {
		f(data)
	}
}
