// Package loopback provides an in-process pair of HCIInterface
// endpoints connected through a virtual controller. It lets two
// instances of the BLE stack drive each other end-to-end (advertise →
// scan → connect → ATT/SMP exchange → disconnect) without any real
// hardware or kernel involvement, which is invaluable for integration
// testing.
//
// The loopback is NOT a complete controller — it implements just
// enough of the HCI command set and event flow for the stack's
// happy-path bring-up + GAP/ACL traffic. Unimplemented commands are
// answered with a generic success CommandComplete so the host's
// command queue never stalls; these are logged at debug level so
// missing functionality is visible during development.
//
// The shared World object owns connection-handle allocation, the
// (peer-visible) advertising state of each endpoint, and routes ACL
// data between the two sides.
package loopback

import (
	"encoding/binary"
	"errors"
	"sync"
	"time"

	hciconst "github.com/BertoldVdb/go-ble/hci/const"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	"github.com/sirupsen/logrus"
)

// World is the shared state of the two endpoints. It tracks each
// side's advertising/scan state, allocates connection handles, and
// routes ACL traffic by handle.
type World struct {
	mu sync.Mutex

	logger *logrus.Entry

	a, b *Endpoint

	nextHandle uint16

	// Each handle maps to the endpoint that *received* the LE
	// Connection Complete and uses the handle to send/receive ACL.
	// Both endpoints share a handle (full-duplex link) but each side
	// owns it locally.
	connections map[uint16]*linkState
}

type linkState struct {
	handleA, handleB uint16 // handle as seen by side A / side B (we use the same)
	roleA            uint8  // 0=central, 1=peripheral
	addrA, addrB     [6]byte
	addrTypeA        uint8
	addrTypeB        uint8

	// Current connection parameters. Updated by Connection Update or
	// remote-param-request flows.
	interval uint16
	latency  uint16
	timeout  uint16

	// Endpoint that plays the central role for this link. Required to
	// route an LE Remote Connection Parameter Request to the right side
	// when the peripheral wants to update params.
	central *Endpoint

	// Encryption state.
	encrypted  bool
	pendingLTK [16]byte // LTK provided by central in LEEnableEncryption
}

// NewWorld returns the shared world plus the two endpoints that link
// to it. Each endpoint satisfies hciinterface.HCIInterface.
func NewWorld(logger *logrus.Entry) (*World, *Endpoint, *Endpoint) {
	w := &World{
		logger:      logger,
		nextHandle:  0x40,
		connections: make(map[uint16]*linkState),
	}
	a := newEndpoint(w, "A", [6]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11})
	b := newEndpoint(w, "B", [6]byte{0x22, 0x22, 0x22, 0x22, 0x22, 0x22})
	w.a = a
	w.b = b
	return w, a, b
}

// allocHandle reserves a fresh 12-bit connection handle.
func (w *World) allocHandle() uint16 {
	h := w.nextHandle
	w.nextHandle++
	if w.nextHandle == 0xFFF {
		w.nextHandle = 0x40
	}
	return h
}

// Endpoint is one side of the loopback. It satisfies
// hciinterface.HCIInterface.
type Endpoint struct {
	world *World
	name  string

	mu        sync.Mutex
	deliverMu sync.Mutex
	closed    bool
	rxHandler hciinterface.HCIRxHandler

	// State affected by host commands.
	bdAddr        [6]byte
	randomAddr    [6]byte
	advParams     leSetAdvertisingParameters
	advData       []byte
	scanRspData   []byte
	advEnabled    bool
	scanParams    leSetScanParameters
	scanEnabled   bool
	whitelist     map[[7]byte]struct{} // type+addr
	pendingCreate *leCreateConnection

	// Per-link txOutstanding for NumberOfCompletedPackets accounting.
	txOutstanding map[uint16]int

	stopCh chan struct{}
}

func newEndpoint(w *World, name string, bdAddr [6]byte) *Endpoint {
	return &Endpoint{
		world:         w,
		name:          name,
		bdAddr:        bdAddr,
		whitelist:     make(map[[7]byte]struct{}),
		txOutstanding: make(map[uint16]int),
		stopCh:        make(chan struct{}),
	}
}

// peer returns the other endpoint.
func (e *Endpoint) peer() *Endpoint {
	if e == e.world.a {
		return e.world.b
	}
	return e.world.a
}

// Run blocks until Close is called. The loopback is event-driven (work
// happens in SendPacket and via background tickers in the World), so
// Run has nothing to do but wait.
func (e *Endpoint) Run() error {
	<-e.stopCh
	return nil
}

func (e *Endpoint) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return nil
	}
	e.closed = true
	close(e.stopCh)
	return nil
}

func (e *Endpoint) SetRecvHandler(h hciinterface.HCIRxHandler) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rxHandler = h
	return nil
}

// deliver hands an HCI packet to the host. Calls to the rx handler are
// serialised per endpoint so the events package's cached decoder
// objects are never overwritten concurrently.
func (e *Endpoint) deliver(data []byte) {
	e.deliverMu.Lock()
	defer e.deliverMu.Unlock()
	e.mu.Lock()
	h := e.rxHandler
	closed := e.closed
	e.mu.Unlock()
	if closed || h == nil {
		return
	}
	pkt := hciinterface.HCIRxPacket{
		Received: true,
		Data:     append([]byte(nil), data...),
		RxTime:   time.Now(),
	}
	_ = h(pkt)
}

// SendPacket is called by the host to send an HCI command or ACL data
// to the controller.
func (e *Endpoint) SendPacket(pkt hciinterface.HCITxPacket) error {
	if len(pkt.Data) < 1 {
		return errors.New("loopback: empty HCI packet")
	}
	switch pkt.Data[0] {
	case hciconst.MsgTypeCommand:
		return e.handleCommand(pkt.Data[1:])
	case hciconst.MsgTypeACL:
		return e.handleACL(pkt.Data[1:])
	default:
		// SCO / ISO / unknown — silently dropped.
		return nil
	}
}

// --- Command framing helpers ---

// handleCommand decodes the HCI command header and dispatches.
func (e *Endpoint) handleCommand(data []byte) error {
	if len(data) < 3 {
		return errors.New("loopback: short HCI command")
	}
	opcode := binary.LittleEndian.Uint16(data[0:2])
	plen := int(data[2])
	if 3+plen > len(data) {
		return errors.New("loopback: HCI command body shorter than declared")
	}
	params := data[3 : 3+plen]

	if e.world.logger != nil && e.world.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		e.world.logger.WithFields(logrus.Fields{
			"0endpoint": e.name,
			"1opcode":   opcode,
			"2params":   params,
		}).Trace("loopback: command")
	}

	return e.dispatchCommand(opcode, params)
}

// sendCommandComplete builds and delivers an HCI Command Complete event
// (eventCode = 0x0E) for the given opcode with the supplied return
// parameters. NumHCICommandPackets is hard-coded to 1 (we accept one
// outstanding command at a time, which matches what the cmdmgr issues).
func (e *Endpoint) sendCommandComplete(opcode uint16, ret []byte) {
	bodyLen := 1 + 2 + len(ret) // numCmdPackets + opcode + return params
	pkt := make([]byte, 3+bodyLen)
	pkt[0] = hciconst.MsgTypeEvent
	pkt[1] = 0x0E
	pkt[2] = byte(bodyLen)
	pkt[3] = 1 // NumHCICommandPackets — 1 slot available
	binary.LittleEndian.PutUint16(pkt[4:6], opcode)
	copy(pkt[6:], ret)
	e.deliver(pkt)
}

// sendCommandStatus builds and delivers an HCI Command Status event
// (eventCode = 0x0F). Used by commands that return Status before
// completion (LECreateConnection in particular).
func (e *Endpoint) sendCommandStatus(opcode uint16, status uint8) {
	pkt := []byte{
		hciconst.MsgTypeEvent,
		0x0F,
		4,        // body length
		status,   // Status
		1,        // NumHCICommandPackets
		byte(opcode), byte(opcode >> 8),
	}
	e.deliver(pkt)
}

// sendLEMeta builds and delivers an HCI LE Meta event (0x3E) with the
// given subevent code and payload.
func (e *Endpoint) sendLEMeta(subevent uint8, payload []byte) {
	body := append([]byte{subevent}, payload...)
	pkt := append([]byte{hciconst.MsgTypeEvent, 0x3E, byte(len(body))}, body...)
	e.deliver(pkt)
}

// sendDisconnectionComplete fires the standard 0x05 event.
func (e *Endpoint) sendDisconnectionComplete(handle uint16, reason uint8) {
	pkt := []byte{
		hciconst.MsgTypeEvent,
		0x05,
		4,
		0x00, // status
		byte(handle), byte(handle >> 8),
		reason,
	}
	e.deliver(pkt)
}

// sendConnectionUpdateComplete fires LE Meta event subevent 0x03 with
// the parameters the host should now consider in effect.
func (e *Endpoint) sendConnectionUpdateComplete(handle, interval, latency, timeout uint16) {
	body := make([]byte, 9)
	body[0] = 0x00 // status
	binary.LittleEndian.PutUint16(body[1:3], handle)
	binary.LittleEndian.PutUint16(body[3:5], interval)
	binary.LittleEndian.PutUint16(body[5:7], latency)
	binary.LittleEndian.PutUint16(body[7:9], timeout)
	e.sendLEMeta(0x03, body)
}

// sendRemoteConnectionParameterRequest fires LE Meta event subevent
// 0x06 to ask the host to accept new connection parameters proposed by
// the remote peer. The host replies with LERemoteConnectionParameterRequestReply
// or NegativeReply.
func (e *Endpoint) sendRemoteConnectionParameterRequest(handle, intervalMin, intervalMax, latency, timeout uint16) {
	body := make([]byte, 10)
	binary.LittleEndian.PutUint16(body[0:2], handle)
	binary.LittleEndian.PutUint16(body[2:4], intervalMin)
	binary.LittleEndian.PutUint16(body[4:6], intervalMax)
	binary.LittleEndian.PutUint16(body[6:8], latency)
	binary.LittleEndian.PutUint16(body[8:10], timeout)
	e.sendLEMeta(0x06, body)
}

// sendLELongTermKeyRequest fires LE Meta event subevent 0x05 to ask
// the peripheral host for the LTK matching the given EDIV/Rand.
func (e *Endpoint) sendLELongTermKeyRequest(handle uint16, ediv uint16, rand [8]byte) {
	body := make([]byte, 12)
	binary.LittleEndian.PutUint16(body[0:2], handle)
	copy(body[2:10], rand[:])
	binary.LittleEndian.PutUint16(body[10:12], ediv)
	e.sendLEMeta(0x05, body)
}

// sendEncryptionChange fires the standard 0x08 event with the encryption
// status (0 = enabled OK, non-zero = HCI error code) and the new state.
func (e *Endpoint) sendEncryptionChange(handle uint16, status uint8, enabled uint8) {
	pkt := []byte{
		hciconst.MsgTypeEvent,
		0x08,
		4,
		status,
		byte(handle), byte(handle >> 8),
		enabled,
	}
	e.deliver(pkt)
}

// sendNumCompletedPackets reports `count` packets completed for handle.
func (e *Endpoint) sendNumCompletedPackets(handle uint16, count uint16) {
	pkt := []byte{
		hciconst.MsgTypeEvent,
		0x13,
		5,
		1, // num handles
		byte(handle), byte(handle >> 8),
		byte(count), byte(count >> 8),
	}
	e.deliver(pkt)
}
