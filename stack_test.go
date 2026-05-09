package ble

import (
	"context"
	"io"
	"testing"

	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	"github.com/sirupsen/logrus"
)

// nopHCIInterface is a never-deliver HCI interface used to verify that
// stack.New survives partial / sparse BluetoothStackConfig values
// without nil-panicking inside a sub-component.
type nopHCIInterface struct{ closed chan struct{} }

func newNopHCI() *nopHCIInterface                                      { return &nopHCIInterface{closed: make(chan struct{})} }
func (n *nopHCIInterface) Run() error                                  { <-n.closed; return nil }
func (n *nopHCIInterface) Close() error                                { close(n.closed); return nil }
func (n *nopHCIInterface) SetRecvHandler(hciinterface.HCIRxHandler) error { return nil }
func (n *nopHCIInterface) SendPacket(hciinterface.HCITxPacket) error   { return nil }

func silentLogger() *logrus.Entry {
	l := logrus.New()
	l.Out = io.Discard
	l.Level = logrus.PanicLevel
	return logrus.NewEntry(l)
}

// New must default every sub-config the caller leaves nil. A
// partially-populated config must not nil-panic inside HCI controller,
// scanner, advertiser, connecter or SMP construction.
func TestNewWithPartialConfig(t *testing.T) {
	cases := []struct {
		name string
		cfg  *BluetoothStackConfig
	}{
		{"all nil sub-configs", &BluetoothStackConfig{}},
		{"only HCI set", &BluetoothStackConfig{HCIControllerConfig: nil}},
		{"only scanner-use flag", &BluetoothStackConfig{BLEScannerUse: true}},
		{"connecter without config", &BluetoothStackConfig{BLEConnecterUse: true}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic: %v", r)
				}
			}()
			dev := newNopHCI()
			defer dev.Close()
			s := New(silentLogger(), c.cfg, dev)
			if s == nil {
				t.Fatal("New returned nil")
			}
			if s.Controller == nil {
				t.Error("Controller is nil")
			}
			if s.BLEScanner == nil {
				t.Error("BLEScanner is nil")
			}
			if s.BLEAdvertiser == nil {
				t.Error("BLEAdvertiser is nil")
			}
			if s.BLEConnecter == nil {
				t.Error("BLEConnecter is nil")
			}
		})
	}
}

// New(_, nil, _) must not panic.
func TestNewWithNilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: %v", r)
		}
	}()
	dev := newNopHCI()
	defer dev.Close()
	if s := New(silentLogger(), nil, dev); s == nil {
		t.Fatal("New returned nil for nil config")
	}
}

// Cancel of context just to keep it referenced.
var _ = context.Background
