package loopback

import (
	"encoding/binary"

	"github.com/sirupsen/logrus"
)

// HCI opcode = OCF | (OGF << 10).
//
//	OGF 0x01 = Link Control
//	OGF 0x03 = Controller & Baseband
//	OGF 0x04 = Informational Parameters
//	OGF 0x08 = LE Controller
const (
	opcodeReset                = 0x0c03 // OGF=3, OCF=0x003
	opcodeSetEventMask         = 0x0c01 // OGF=3, OCF=0x001
	opcodeWriteLEHostSupport   = 0x0c6d // OGF=3, OCF=0x06d
	opcodeWriteFlowControlMode = 0x0c67 // OGF=3, OCF=0x067

	opcodeReadLocalVersion       = 0x1001 // OGF=4
	opcodeReadLocalCommands      = 0x1002
	opcodeReadLocalFeatures      = 0x1003
	opcodeReadBufferSize         = 0x1005
	opcodeReadBDADDR             = 0x1009

	opcodeLinkControlDisconnect = 0x0406 // OGF=1, OCF=0x006

	opcodeLESetEventMask                = 0x2001 // OGF=8
	opcodeLEReadBufferSize              = 0x2002
	opcodeLEReadLocalSupportedFeatures  = 0x2003
	opcodeLESetRandomAddress            = 0x2005
	opcodeLESetAdvertisingParameters    = 0x2006
	opcodeLEReadAdvChannelTxPower       = 0x2007
	opcodeLESetAdvertisingData          = 0x2008
	opcodeLESetScanResponseData         = 0x2009
	opcodeLESetAdvertisingEnable        = 0x200a
	opcodeLESetScanParameters           = 0x200b
	opcodeLESetScanEnable               = 0x200c
	opcodeLECreateConnection            = 0x200d
	opcodeLECreateConnectionCancel      = 0x200e
	opcodeLEReadWhiteListSize           = 0x200f
	opcodeLEClearWhiteList              = 0x2010
	opcodeLEAddDeviceToWhiteList        = 0x2011
	opcodeLERemoveDeviceFromWhiteList   = 0x2012
	opcodeLEConnectionUpdate            = 0x2013
	opcodeLEEnableEncryption            = 0x2019
	opcodeLELongTermKeyRequestReply     = 0x201a
	opcodeLELongTermKeyRequestNegReply  = 0x201b
	opcodeLERand                        = 0x2018
	opcodeLEReadBufferSizeV2            = 0x2060
	opcodeLERemoteConnParamReqReply     = 0x2020
	opcodeLERemoteConnParamReqNegReply  = 0x2021
)

// dispatchCommand routes an HCI command to its handler. Unknown
// commands respond with a generic success Command Complete (status=0)
// so the host's command queue does not stall.
func (e *Endpoint) dispatchCommand(opcode uint16, params []byte) error {
	switch opcode {
	case opcodeReset:
		e.sendCommandComplete(opcode, []byte{0x00})
	case opcodeSetEventMask, opcodeLESetEventMask, opcodeWriteLEHostSupport,
		opcodeWriteFlowControlMode:
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeReadLocalVersion:
		// Status, HCIVersion, HCIRevision, LMPVersion, ManufacturerName, LMPSubversion
		ret := []byte{
			0x00,       // status
			0x0B,       // HCI 5.2
			0x00, 0x00, // revision
			0x0B,       // LMP 5.2
			0x00, 0x00, // manufacturer
			0x00, 0x00, // LMP subversion
		}
		e.sendCommandComplete(opcode, ret)

	case opcodeReadLocalCommands:
		ret := make([]byte, 1+64)
		// Status=0; supported-commands bit-vector all-ones.
		for i := 1; i < len(ret); i++ {
			ret[i] = 0xFF
		}
		e.sendCommandComplete(opcode, ret)

	case opcodeReadLocalFeatures:
		// Status + 8 LMP features; LE supported (bit 38 = byte 4 bit 6).
		ret := make([]byte, 1+8)
		ret[1+4] = 0x40
		e.sendCommandComplete(opcode, ret)

	case opcodeLEReadLocalSupportedFeatures:
		// Status + 8-byte LE features bitmap. We claim Encryption (bit 0)
		// and LE Secure Connections (bit 28 — high byte).
		ret := make([]byte, 1+8)
		ret[1] = 0x01
		ret[1+3] = 0x10 // bit 28 = byte 3 bit 4 (LL Privacy / SC differ across versions; safe-ish default)
		e.sendCommandComplete(opcode, ret)

	case opcodeReadBufferSize:
		// Status, ACLDataPacketLength(2), SyncDataPacketLength(1),
		// TotalNumACLDataPackets(2), TotalNumSyncDataPackets(2)
		ret := []byte{
			0x00,
			0x00, 0x01, // 256-byte ACL
			0x00,
			0x10, 0x00, // 16 ACL slots
			0x00, 0x00,
		}
		e.sendCommandComplete(opcode, ret)

	case opcodeLEReadBufferSize:
		// Status, LEACLDataPacketLength(2), TotalNumLEACLDataPackets(1)
		ret := []byte{0x00, 0x00, 0x01, 0x10}
		e.sendCommandComplete(opcode, ret)

	case opcodeLEReadBufferSizeV2:
		// Some controllers don't support v2; return command-not-supported
		// (HCI error 0x01). The stack falls back to v1.
		e.sendCommandComplete(opcode, []byte{0x01})

	case opcodeReadBDADDR:
		ret := append([]byte{0x00}, e.bdAddr[:]...)
		e.sendCommandComplete(opcode, ret)

	case opcodeLESetRandomAddress:
		if len(params) >= 6 {
			e.world.mu.Lock()
			copy(e.randomAddr[:], params[:6])
			e.world.mu.Unlock()
		}
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLEReadAdvChannelTxPower:
		e.sendCommandComplete(opcode, []byte{0x00, 0x00})

	case opcodeLESetAdvertisingParameters:
		e.world.mu.Lock()
		_ = e.advParams.decode(params)
		e.world.mu.Unlock()
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLESetAdvertisingData:
		if len(params) >= 1 {
			n := int(params[0])
			if n > 31 {
				n = 31
			}
			if 1+n <= len(params) {
				e.world.mu.Lock()
				e.advData = append(e.advData[:0], params[1:1+n]...)
				e.world.mu.Unlock()
			}
		}
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLESetScanResponseData:
		if len(params) >= 1 {
			n := int(params[0])
			if n > 31 {
				n = 31
			}
			if 1+n <= len(params) {
				e.world.mu.Lock()
				e.scanRspData = append(e.scanRspData[:0], params[1:1+n]...)
				e.world.mu.Unlock()
			}
		}
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLESetAdvertisingEnable:
		if len(params) >= 1 {
			e.world.mu.Lock()
			e.advEnabled = params[0] != 0
			e.world.mu.Unlock()
		}
		e.sendCommandComplete(opcode, []byte{0x00})
		// Trigger an advertising-report poll on the peer if scanning.
		go e.world.checkAdvScan()

	case opcodeLESetScanParameters:
		e.world.mu.Lock()
		_ = e.scanParams.decode(params)
		e.world.mu.Unlock()
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLESetScanEnable:
		if len(params) >= 1 {
			e.world.mu.Lock()
			e.scanEnabled = params[0] != 0
			e.world.mu.Unlock()
		}
		e.sendCommandComplete(opcode, []byte{0x00})
		go e.world.checkAdvScan()

	case opcodeLEClearWhiteList:
		e.world.mu.Lock()
		e.whitelist = make(map[[7]byte]struct{})
		e.world.mu.Unlock()
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLEAddDeviceToWhiteList:
		if len(params) >= 7 {
			var key [7]byte
			copy(key[:], params[:7])
			e.world.mu.Lock()
			e.whitelist[key] = struct{}{}
			e.world.mu.Unlock()
		}
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLERemoveDeviceFromWhiteList:
		if len(params) >= 7 {
			var key [7]byte
			copy(key[:], params[:7])
			e.world.mu.Lock()
			delete(e.whitelist, key)
			e.world.mu.Unlock()
		}
		e.sendCommandComplete(opcode, []byte{0x00})

	case opcodeLEReadWhiteListSize:
		e.sendCommandComplete(opcode, []byte{0x00, 0x10}) // 16 entries

	case opcodeLECreateConnection:
		var c leCreateConnection
		if err := c.decode(params); err != nil {
			e.sendCommandStatus(opcode, 0x12) // invalid params
			return nil
		}
		e.world.mu.Lock()
		e.pendingCreate = &c
		e.world.mu.Unlock()
		// CreateConnection returns Status, with the actual completion
		// arriving via the LE Connection Complete event.
		e.sendCommandStatus(opcode, 0x00)
		go e.world.tryCompleteConnection()

	case opcodeLECreateConnectionCancel:
		e.world.mu.Lock()
		had := e.pendingCreate != nil
		e.pendingCreate = nil
		e.world.mu.Unlock()
		e.sendCommandComplete(opcode, []byte{0x00})
		if had {
			// Spec: a Cancel after a Create that wasn't completed yet
			// must produce a Connection Complete event with status 0x02
			// (Unknown Connection Identifier).
			payload := make([]byte, 18)
			// SubeventCode handled by sendLEMeta; payload is everything
			// after the subevent code. ConnectionComplete payload:
			// Status, ConnectionHandle(2), Role(1), PeerAddrType(1),
			// PeerAddr(6), Interval(2), Latency(2), Timeout(2), MasterClockAccuracy(1)
			payload[0] = 0x02 // status
			e.sendLEMeta(0x01, payload)
		}

	case opcodeLinkControlDisconnect:
		if len(params) >= 3 {
			handle := binary.LittleEndian.Uint16(params[0:2])
			reason := params[2]
			e.sendCommandStatus(opcode, 0x00)
			go e.world.disconnect(e, handle, reason)
		} else {
			e.sendCommandStatus(opcode, 0x12)
		}

	case opcodeLEConnectionUpdate:
		// Vol 4 Part E §7.8.18: Handle, IntervalMin, IntervalMax,
		// Latency, SupervisionTimeout, MinCELength, MaxCELength.
		if len(params) < 14 {
			e.sendCommandStatus(opcode, 0x12) // invalid params
			return nil
		}
		handle := binary.LittleEndian.Uint16(params[0:2])
		intervalMin := binary.LittleEndian.Uint16(params[2:4])
		intervalMax := binary.LittleEndian.Uint16(params[4:6])
		latency := binary.LittleEndian.Uint16(params[6:8])
		timeout := binary.LittleEndian.Uint16(params[8:10])
		e.sendCommandStatus(opcode, 0x00)
		// Run synchronously so multiple LEConnectionUpdate commands
		// from different endpoints don't race in the World — the
		// dispatcher is single-threaded per endpoint, and we serialise
		// at the World level via world.mu inside applyConnectionUpdate.
		e.world.connectionUpdateRequested(e, handle, intervalMin, intervalMax, latency, timeout)

	case opcodeLERemoteConnParamReqReply:
		// Vol 4 Part E §7.8.31: Handle, IntervalMin, IntervalMax,
		// Latency, Timeout, MinCELength, MaxCELength.
		if len(params) < 14 {
			e.sendCommandComplete(opcode, []byte{0x12, 0, 0})
			return nil
		}
		handle := binary.LittleEndian.Uint16(params[0:2])
		intervalMin := binary.LittleEndian.Uint16(params[2:4])
		intervalMax := binary.LittleEndian.Uint16(params[4:6])
		latency := binary.LittleEndian.Uint16(params[6:8])
		timeout := binary.LittleEndian.Uint16(params[8:10])
		// Return: Status (1) + ConnectionHandle (2)
		ret := make([]byte, 3)
		binary.LittleEndian.PutUint16(ret[1:3], handle)
		e.sendCommandComplete(opcode, ret)
		e.world.applyConnectionUpdate(handle, intervalMin, intervalMax, latency, timeout)

	case opcodeLERemoteConnParamReqNegReply:
		// Reject — no LE Connection Update Complete is emitted.
		if len(params) < 3 {
			e.sendCommandComplete(opcode, []byte{0x12, 0, 0})
			return nil
		}
		handle := binary.LittleEndian.Uint16(params[0:2])
		ret := make([]byte, 3)
		binary.LittleEndian.PutUint16(ret[1:3], handle)
		e.sendCommandComplete(opcode, ret)

	case opcodeLEEnableEncryption:
		// Vol 4 Part E §7.8.24 — central commands the controller to
		// start encryption. Parameters: ConnectionHandle (2),
		// RandomNumber (8), EncryptedDiversifier (2), LongTermKey (16).
		if len(params) < 28 {
			e.sendCommandStatus(opcode, 0x12)
			return nil
		}
		handle := binary.LittleEndian.Uint16(params[0:2])
		var rand [8]byte
		copy(rand[:], params[2:10])
		ediv := binary.LittleEndian.Uint16(params[10:12])
		var ltk [16]byte
		copy(ltk[:], params[12:28])
		e.sendCommandStatus(opcode, 0x00)
		e.world.encryptionRequested(e, handle, ediv, rand, ltk)

	case opcodeLELongTermKeyRequestReply:
		// Vol 4 Part E §7.8.25: ConnectionHandle (2) + LTK (16).
		// Returns Status (1) + ConnectionHandle (2).
		if len(params) < 18 {
			e.sendCommandComplete(opcode, []byte{0x12, 0, 0})
			return nil
		}
		handle := binary.LittleEndian.Uint16(params[0:2])
		var ltk [16]byte
		copy(ltk[:], params[2:18])
		ret := make([]byte, 3)
		binary.LittleEndian.PutUint16(ret[1:3], handle)
		e.sendCommandComplete(opcode, ret)
		e.world.encryptionLTKReplied(handle, ltk)

	case opcodeLELongTermKeyRequestNegReply:
		// Vol 4 Part E §7.8.26: ConnectionHandle (2).
		if len(params) < 2 {
			e.sendCommandComplete(opcode, []byte{0x12, 0, 0})
			return nil
		}
		handle := binary.LittleEndian.Uint16(params[0:2])
		ret := make([]byte, 3)
		binary.LittleEndian.PutUint16(ret[1:3], handle)
		e.sendCommandComplete(opcode, ret)
		e.world.encryptionRejected(handle)

	case opcodeLERand:
		// Status + 8 random bytes (deterministic to keep tests reproducible).
		ret := make([]byte, 1+8)
		ret[1] = 0xAB
		ret[2] = 0xCD
		ret[3] = 0xEF
		ret[4] = 0x01
		ret[5] = 0x23
		ret[6] = 0x45
		ret[7] = 0x67
		ret[8] = 0x89
		e.sendCommandComplete(opcode, ret)

	default:
		// Unknown command — answer with success so the host doesn't stall.
		// At debug level note that we faked it.
		if e.world.logger != nil {
			e.world.logger.WithFields(logrus.Fields{
				"0endpoint": e.name,
				"1opcode":   opcode,
			}).Debug("loopback: unhandled command, responding with generic success")
		}
		e.sendCommandComplete(opcode, []byte{0x00})
	}
	return nil
}

// --- Command parameter decoders ---

type leSetAdvertisingParameters struct {
	IntervalMin uint16
	IntervalMax uint16
	AdvType     uint8
	OwnAddrType uint8
	PeerType    uint8
	PeerAddr    [6]byte
	ChannelMap  uint8
	FilterPol   uint8
}

func (s *leSetAdvertisingParameters) decode(p []byte) error {
	if len(p) < 15 {
		return errShortParams
	}
	s.IntervalMin = binary.LittleEndian.Uint16(p[0:2])
	s.IntervalMax = binary.LittleEndian.Uint16(p[2:4])
	s.AdvType = p[4]
	s.OwnAddrType = p[5]
	s.PeerType = p[6]
	copy(s.PeerAddr[:], p[7:13])
	s.ChannelMap = p[13]
	s.FilterPol = p[14]
	return nil
}

type leSetScanParameters struct {
	ScanType    uint8
	Interval    uint16
	Window      uint16
	OwnAddrType uint8
	FilterPol   uint8
}

func (s *leSetScanParameters) decode(p []byte) error {
	if len(p) < 7 {
		return errShortParams
	}
	s.ScanType = p[0]
	s.Interval = binary.LittleEndian.Uint16(p[1:3])
	s.Window = binary.LittleEndian.Uint16(p[3:5])
	s.OwnAddrType = p[5]
	s.FilterPol = p[6]
	return nil
}

type leCreateConnection struct {
	ScanInterval        uint16
	ScanWindow          uint16
	InitiatorFilterPol  uint8
	PeerAddrType        uint8
	PeerAddr            [6]byte
	OwnAddrType         uint8
	IntervalMin         uint16
	IntervalMax         uint16
	Latency             uint16
	SupervisionTimeout  uint16
	MinCELength         uint16
	MaxCELength         uint16
}

func (s *leCreateConnection) decode(p []byte) error {
	if len(p) < 25 {
		return errShortParams
	}
	s.ScanInterval = binary.LittleEndian.Uint16(p[0:2])
	s.ScanWindow = binary.LittleEndian.Uint16(p[2:4])
	s.InitiatorFilterPol = p[4]
	s.PeerAddrType = p[5]
	copy(s.PeerAddr[:], p[6:12])
	s.OwnAddrType = p[12]
	s.IntervalMin = binary.LittleEndian.Uint16(p[13:15])
	s.IntervalMax = binary.LittleEndian.Uint16(p[15:17])
	s.Latency = binary.LittleEndian.Uint16(p[17:19])
	s.SupervisionTimeout = binary.LittleEndian.Uint16(p[19:21])
	s.MinCELength = binary.LittleEndian.Uint16(p[21:23])
	s.MaxCELength = binary.LittleEndian.Uint16(p[23:25])
	return nil
}
