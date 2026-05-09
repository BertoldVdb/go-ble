package loopback

import (
	"encoding/binary"
	"errors"
)

var errShortParams = errors.New("loopback: HCI command parameters truncated")

// checkAdvScan checks whether one side is advertising and the other
// scanning, and emits an LE Advertising Report event to the scanner.
// Called whenever advertising or scanning state changes.
func (w *World) checkAdvScan() {
	w.maybeReport(w.a, w.b)
	w.maybeReport(w.b, w.a)
}

func (w *World) maybeReport(adv, scan *Endpoint) {
	w.mu.Lock()
	if !adv.advEnabled || !scan.scanEnabled {
		w.mu.Unlock()
		return
	}
	addr := adv.peerSourceAddr()
	addrType := adv.peerSourceAddrType()
	advType := adv.advParams.AdvType
	data := append([]byte(nil), adv.advData...)
	scanRsp := append([]byte(nil), adv.scanRspData...)
	scanIsActive := scan.scanParams.ScanType == 1
	w.mu.Unlock()

	// LE Advertising Report payload (subevent 0x02):
	//   NumReports (1)
	//   EventType[NumReports] (each 1 byte)
	//   AddressType[NumReports] (each 1 byte)
	//   Address[NumReports] (each 6 bytes)
	//   DataLength[NumReports] (each 1 byte)
	//   Data[NumReports] (variable)
	//   RSSI[NumReports] (each 1 byte)
	emit := func(eventType uint8, payload []byte) {
		body := []byte{1}
		body = append(body, eventType)
		body = append(body, addrType)
		body = append(body, addr[:]...)
		body = append(body, byte(len(payload)))
		body = append(body, payload...)
		body = append(body, 0xCE) // RSSI = -50
		scan.sendLEMeta(0x02, body)
	}

	emit(advType, data)

	// In active scanning the scanner emits a SCAN_REQ and the
	// advertiser responds with SCAN_RSP. The host sees the SCAN_RSP as
	// a separate Advertising Report with EventType=4 (SCAN_RSP).
	if scanIsActive && len(scanRsp) > 0 {
		emit(4, scanRsp)
	}
}

// peerSourceAddr returns the address this endpoint *appears as* on the
// air when advertising. If the host has set a random address and asked
// for OwnAddressType=1 in the advertising parameters, we use the
// random address; otherwise the public BD_ADDR.
func (e *Endpoint) peerSourceAddr() [6]byte {
	if e.advParams.OwnAddrType == 1 {
		return e.randomAddr
	}
	return e.bdAddr
}

func (e *Endpoint) peerSourceAddrType() uint8 {
	return e.advParams.OwnAddrType
}

// tryCompleteConnection looks for a CreateConnection-pending side and
// a matching advertising peer; if found, fires LE Connection Complete
// on both sides.
func (w *World) tryCompleteConnection() {
	w.mu.Lock()
	defer w.mu.Unlock()

	tryPair := func(initiator, advertiser *Endpoint) bool {
		c := initiator.pendingCreate
		if c == nil {
			return false
		}
		if !advertiser.advEnabled {
			return false
		}
		// Match by peer address (or InitiatorFilterPolicy=1 → use whitelist;
		// we already populated whitelist with the desired peer, so accept
		// the advertiser if its source addr is in initiator.whitelist).
		match := false
		var key [7]byte
		advAddr := advertiser.peerSourceAddr()
		key[0] = advertiser.peerSourceAddrType()
		copy(key[1:], advAddr[:])
		if c.InitiatorFilterPol == 1 {
			if _, ok := initiator.whitelist[key]; ok {
				match = true
			}
		} else {
			if c.PeerAddrType == advertiser.peerSourceAddrType() &&
				c.PeerAddr == advertiser.peerSourceAddr() {
				match = true
			}
		}
		if !match {
			return false
		}

		handle := w.allocHandle()
		// Fire LE Connection Complete on both sides.
		// Subevent 0x01 payload:
		//   Status(1), ConnectionHandle(2), Role(1), PeerAddrType(1),
		//   PeerAddr(6), Interval(2), Latency(2), Timeout(2),
		//   MasterClockAccuracy(1)
		buildPayload := func(role uint8, peerType uint8, peerAddr [6]byte) []byte {
			b := make([]byte, 19)
			b[0] = 0x00 // status
			binary.LittleEndian.PutUint16(b[1:3], handle)
			b[3] = role
			b[4] = peerType
			copy(b[5:11], peerAddr[:])
			binary.LittleEndian.PutUint16(b[11:13], c.IntervalMin)
			binary.LittleEndian.PutUint16(b[13:15], c.Latency)
			binary.LittleEndian.PutUint16(b[15:17], c.SupervisionTimeout)
			b[17] = 0x00
			return b[:18]
		}

		// Initiator perspective: role=0 (central), peer=advertiser.
		initiator.sendLEMeta(0x01, buildPayload(0, advertiser.peerSourceAddrType(), advertiser.peerSourceAddr()))
		// Advertiser perspective: role=1 (peripheral), peer=initiator.
		advertiser.sendLEMeta(0x01, buildPayload(1, initiator.advParams.OwnAddrType, initiator.peerSourceAddr()))

		// Stop advertising on the peer (link is now established).
		advertiser.advEnabled = false

		// Bookkeep the link.
		w.connections[handle] = &linkState{
			handleA:   handle,
			handleB:   handle,
			roleA:     0,
			addrA:     initiator.peerSourceAddr(),
			addrTypeA: initiator.peerSourceAddrType(),
			addrB:     advertiser.peerSourceAddr(),
			addrTypeB: advertiser.peerSourceAddrType(),
			interval:  c.IntervalMax,
			latency:   c.Latency,
			timeout:   c.SupervisionTimeout,
			central:   initiator,
		}
		initiator.txOutstanding[handle] = 0
		advertiser.txOutstanding[handle] = 0
		initiator.pendingCreate = nil
		return true
	}

	tryPair(w.a, w.b)
	tryPair(w.b, w.a)
}

// connectionUpdateRequested is invoked when an endpoint issued
// HCI LEConnectionUpdate. The behaviour depends on the role:
//
//   - Central side: the controller would just program the LL with the
//     new params, so we directly emit LEConnectionUpdateComplete on
//     both sides with the new (negotiated) values.
//   - Peripheral side: a real controller would translate this into an
//     LL_CONNECTION_PARAM_REQ which the central's controller surfaces
//     as an LE Remote Connection Parameter Request event. We mirror
//     that here. The central's host then answers with Reply or
//     NegativeReply, both of which we already handle.
func (w *World) connectionUpdateRequested(by *Endpoint, handle, intervalMin, intervalMax, latency, timeout uint16) {
	w.mu.Lock()
	link, ok := w.connections[handle]
	if !ok {
		w.mu.Unlock()
		return
	}
	central := link.central
	w.mu.Unlock()

	// Pick the actual interval — anywhere in [min,max] is allowed; we
	// take the maximum so it's deterministic and easy to assert.
	interval := intervalMax

	if by == central {
		// Central programs the LL directly.
		w.applyConnectionUpdate(handle, interval, interval, latency, timeout)
		return
	}
	// Peripheral asks the central for permission via the standard
	// LL_CONNECTION_PARAM_REQ → LE Remote Connection Parameter Request
	// event flow.
	central.sendRemoteConnectionParameterRequest(handle, intervalMin, intervalMax, latency, timeout)
}

// applyConnectionUpdate finalises a parameter change: stores the new
// values on the link and emits LEConnectionUpdateComplete on both
// endpoints. Called from the central path (direct update) or from
// LERemoteConnParamReqReply (after the central accepted a peripheral's
// request).
func (w *World) applyConnectionUpdate(handle, intervalMin, intervalMax, latency, timeout uint16) {
	w.mu.Lock()
	link, ok := w.connections[handle]
	if !ok {
		w.mu.Unlock()
		return
	}
	interval := intervalMax
	link.interval = interval
	link.latency = latency
	link.timeout = timeout
	a, b := w.a, w.b
	w.mu.Unlock()
	_ = intervalMin

	a.sendConnectionUpdateComplete(handle, interval, latency, timeout)
	b.sendConnectionUpdateComplete(handle, interval, latency, timeout)
}

// encryptionRequested handles an LEEnableEncryption HCI command. Only
// the central is allowed to start encryption; if the peripheral
// somehow issues this, we silently drop. Records the LTK the central
// supplied (so we can verify it matches the peripheral's reply) and
// emits LELongTermKeyRequest on the peripheral.
func (w *World) encryptionRequested(by *Endpoint, handle uint16, ediv uint16, rand [8]byte, ltk [16]byte) {
	w.mu.Lock()
	link, ok := w.connections[handle]
	if !ok {
		w.mu.Unlock()
		return
	}
	if by != link.central {
		w.mu.Unlock()
		// Peripheral can't start encryption directly — surface an error
		// status on the encryption-change event so the host doesn't hang
		// waiting forever.
		by.sendEncryptionChange(handle, 0x12 /* command disallowed */, 0)
		return
	}
	link.pendingLTK = ltk
	peripheral := link.central.peer()
	w.mu.Unlock()

	peripheral.sendLELongTermKeyRequest(handle, ediv, rand)
}

// encryptionLTKReplied handles the peripheral's LELongTermKeyRequestReply.
// If the LTK matches what the central supplied, encryption succeeds
// and EncryptionChange (status=0, enabled=1) fires on both sides.
// Otherwise EncryptionChange fires with HCI error 0x06 (PIN or Key
// Missing) on the central and the peripheral sees no change.
func (w *World) encryptionLTKReplied(handle uint16, ltk [16]byte) {
	w.mu.Lock()
	link, ok := w.connections[handle]
	if !ok {
		w.mu.Unlock()
		return
	}
	central := link.central
	peripheral := central.peer()
	expected := link.pendingLTK
	link.pendingLTK = [16]byte{}
	match := expected == ltk
	if match {
		link.encrypted = true
	}
	w.mu.Unlock()

	if match {
		central.sendEncryptionChange(handle, 0x00, 0x01)
		peripheral.sendEncryptionChange(handle, 0x00, 0x01)
	} else {
		// Per the spec a key mismatch surfaces as HCI Error 0x06 on the
		// central; the peripheral typically also gets the failure
		// notification so its host can mark the LTK invalid.
		central.sendEncryptionChange(handle, 0x06, 0x00)
		peripheral.sendEncryptionChange(handle, 0x06, 0x00)
	}
}

// encryptionRejected handles the peripheral's NegativeReply (no LTK
// available). Encryption fails with HCI error 0x06.
func (w *World) encryptionRejected(handle uint16) {
	w.mu.Lock()
	link, ok := w.connections[handle]
	if !ok {
		w.mu.Unlock()
		return
	}
	central := link.central
	peripheral := central.peer()
	link.pendingLTK = [16]byte{}
	w.mu.Unlock()
	central.sendEncryptionChange(handle, 0x06, 0x00)
	peripheral.sendEncryptionChange(handle, 0x06, 0x00)
}

// disconnect tears down the link with the given handle and emits
// DisconnectionComplete on both endpoints.
func (w *World) disconnect(initiator *Endpoint, handle uint16, reason uint8) {
	w.mu.Lock()
	_, ok := w.connections[handle]
	if ok {
		delete(w.connections, handle)
	}
	w.mu.Unlock()
	if !ok {
		return
	}
	initiator.sendDisconnectionComplete(handle, reason)
	initiator.peer().sendDisconnectionComplete(handle, reason)
}

// handleACL routes an outgoing ACL packet from one host to the other.
// The wire format is:
//
//	handle_lo, handle_hi (12-bit handle + PB(2)+BC(2) flags),
//	length_lo, length_hi (LE uint16),
//	payload...
func (e *Endpoint) handleACL(data []byte) error {
	if len(data) < 4 {
		return nil // silently drop
	}
	hdr := binary.LittleEndian.Uint16(data[0:2])
	handle := hdr & 0xFFF
	flagPB := (hdr >> 12) & 0x3
	flagBC := (hdr >> 14) & 0x3
	plen := int(binary.LittleEndian.Uint16(data[2:4]))
	if 4+plen > len(data) {
		return nil
	}

	w := e.world
	w.mu.Lock()
	_, ok := w.connections[handle]
	w.mu.Unlock()
	if !ok {
		return nil
	}

	// Forward to peer with PB=2 (start, host→controller becomes
	// "complete L2CAP / first fragment from controller→host" on the
	// other side, which the receiving stack treats as a fresh
	// reassembly via flagPB==2). For simplicity always forward as a
	// complete (start) fragment; the loopback never fragments.
	_ = flagPB
	_ = flagBC

	// Peer receives the packet as type=ACL with PB=2, BC=0.
	peer := e.peer()
	out := make([]byte, 5+plen)
	out[0] = 0x02 // HCI ACL
	outHdr := uint16(handle&0xFFF) | (2 << 12)
	binary.LittleEndian.PutUint16(out[1:3], outHdr)
	binary.LittleEndian.PutUint16(out[3:5], uint16(plen))
	copy(out[5:], data[4:4+plen])
	peer.deliver(out)

	// Send NumberOfCompletedPackets back to the sender so its TX
	// accounting decrements.
	go e.sendNumCompletedPackets(handle, 1)
	return nil
}
