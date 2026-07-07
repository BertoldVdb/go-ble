/*
teleport copies a remote BLE device's beacon to a local transmitter and
forwards connections to it.

Two BluetoothStack instances run inside the same process:

  - "remote" stack: typically pointed at an HCI controller that is in
    radio range of the target device (often via ble2net over TCP). It
    scans for the target's advertisement and acts as central when a
    proxy connection is needed.
  - "local" stack: the controller close to the user. It clones the
    target's beacon out the air verbatim (advertising data + scan
    response) and accepts an incoming peripheral connection.

Once a peer connects to the local advertiser, the remote stack opens a
central connection to the target and the program shuttles traffic
between the two endpoints.

By default the program forwards full L2CAP frames opaquely — ATT, SMP,
fixed-channel traffic — which works for plaintext GATT. Link-layer
encryption is per-link and *not* tunneled, so a peer that insists on a
paired/encrypted link won't go through.

Pass -pin <NNNNNN> to enable independent legacy pairing on each leg:
SMP is terminated locally on both sides using the supplied static
passcode (the peer pairs with us, we pair with the target), and ATT
traffic flows through both encrypted links. Bonding is enabled, so the
keys are persisted under the user's config dir and re-pairing is
skipped on reconnect.

The btsnoop wrap lets either or both HCI streams be captured, which is
the primary motivation: a verbatim trace of the proxied connection from
the user's side, the device's side, or both.

Typical invocation:

	teleport \
	    -device hci0,net -netaddr remote-host:3000 \
	    -target AA:BB:CC:DD:EE:FF \
	    -btsnoop-local local.btsnoop

Pass -peer <mac>,<mac>... to restrict who is allowed to connect to the
cloned advertiser; by default any peer is accepted.
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/BertoldVdb/go-ble"
	"github.com/BertoldVdb/go-ble/bleadvertiser"
	"github.com/BertoldVdb/go-ble/bleconnecter"
	"github.com/BertoldVdb/go-ble/blescanner"
	"github.com/BertoldVdb/go-ble/blesmp"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	"github.com/BertoldVdb/go-ble/hci/drivers/btsnoop"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	blel2cap "github.com/BertoldVdb/go-ble/l2cap"
	bleutil "github.com/BertoldVdb/go-ble/util"
	bleutilparam "github.com/BertoldVdb/go-ble/util/param"
	"github.com/BertoldVdb/go-misc/logrusconfig"
	"github.com/BertoldVdb/go-misc/multirun"
	"github.com/sirupsen/logrus"
)

type teleporter struct {
	logger *logrus.Entry

	// targetMac is the user-supplied MAC. The address type is learned
	// from scan reports and stored in targetAddrType because the same
	// MAC bytes can show up as Public or Random and the connect-side
	// HCI commands need the right type.
	targetMac    bleutil.MacAddr
	allowedPeers []bleutil.BLEAddr

	// pin is the static legacy passcode to use on both legs. -1 disables
	// pairing entirely and the proxy operates as a transparent L2CAP
	// forwarder.
	pin int32

	localStack  *ble.BluetoothStack
	remoteStack *ble.BluetoothStack

	advSlot *bleadvertiser.LegacyAdvertisingSlot

	advMu          sync.Mutex
	targetAddrType bleutil.MacAddrType
	haveAddrType   bool
	beaconPkt      []byte
	scanRspPkt     []byte
	haveBeacon     bool
	haveScanRsp    bool
	advUpdated     chan struct{}

	scanCBHandle blescanner.CallbackHandle

	ctx    context.Context
	cancel context.CancelFunc
}

// captureAdvReport buffers raw advertising and scan-response payloads
// for the target device. Match is on the MAC bytes only — the address
// type (Public vs Random) is learned from the report so that a target
// using a static-random address still gets matched even though the
// user only supplied the MAC on the command line.
func (t *teleporter) captureAdvReport(r *blescanner.BLEAdvertisingReport) bool {
	if r.Addr.MacAddr != t.targetMac {
		return false
	}

	t.advMu.Lock()
	defer t.advMu.Unlock()

	if !t.haveAddrType || t.targetAddrType != r.Addr.MacAddrType {
		t.targetAddrType = r.Addr.MacAddrType
		t.haveAddrType = true
	}

	changed := false
	switch r.PktType {
	case blescanner.EventTypeInd, blescanner.EventTypeScanInd, blescanner.EventTypeNonConnInd:
		if !bytesEqual(t.beaconPkt, r.Data) {
			t.beaconPkt = append(t.beaconPkt[:0], r.Data...)
			changed = true
		}
		if !t.haveBeacon {
			t.haveBeacon = true
			changed = true
		}
	case blescanner.EventTypeScanRsp:
		if !bytesEqual(t.scanRspPkt, r.Data) {
			t.scanRspPkt = append(t.scanRspPkt[:0], r.Data...)
			changed = true
		}
		if !t.haveScanRsp {
			t.haveScanRsp = true
			changed = true
		}
	default:
		return false
	}

	if changed {
		select {
		case t.advUpdated <- struct{}{}:
		default:
		}
	}
	return false
}

// snapshotTargetAddr returns the target's address with the type
// learned from scanning. If no advertisement has been seen yet it
// falls back to Public, but the main loop only calls this after
// waitForAdvData so that path isn't normally hit.
func (t *teleporter) snapshotTargetAddr() bleutil.BLEAddr {
	t.advMu.Lock()
	defer t.advMu.Unlock()
	addrType := bleutil.MacAddrPublic
	if t.haveAddrType {
		addrType = t.targetAddrType
	}
	return bleutil.BLEAddr{MacAddr: t.targetMac, MacAddrType: addrType}
}

// snapshotAdv returns the cached beacon and scan-response payloads.
// The third return is false until at least the connectable
// advertisement payload has been seen.
func (t *teleporter) snapshotAdv() ([]byte, []byte, bool) {
	t.advMu.Lock()
	defer t.advMu.Unlock()

	if !t.haveBeacon {
		return nil, nil, false
	}
	beacon := append([]byte(nil), t.beaconPkt...)
	var scanRsp []byte
	if t.haveScanRsp {
		scanRsp = append([]byte(nil), t.scanRspPkt...)
	}
	return beacon, scanRsp, true
}

// applyAdvData writes the cached beacon and scan-response bytes into
// the advertiser's base slot, leaving Active/Type alone — those are
// flipped by BLEConnecter.Connect (peripheral mode) when it calls
// LegacyAdvertisingSetConnection on the same slot.
func (t *teleporter) applyAdvData() error {
	beacon, scanRsp, ok := t.snapshotAdv()
	if !ok {
		return errors.New("no advertising data captured yet")
	}

	data, _ := t.advSlot.GetData()
	data.BeaconPacket = beacon
	data.ScanPacket = scanRsp
	if data.IntervalMin == 0 {
		data.IntervalMin = 0x20
		data.IntervalMax = 0x40
	}
	_, err := t.advSlot.ReplaceData(true, data)
	return err
}

// pump copies L2CAP-frame buffers from src to dst until either side
// errors. WriteBuffer takes ownership of the buffer (releases it on
// error or after sending), so there's no manual buffer accounting.
func (t *teleporter) pump(ctx context.Context, name string, src, dst hciconnmgr.BufferConn) {
	for {
		buf, err := src.ReadBuffer(ctx)
		if err != nil {
			t.logger.WithError(err).WithField("dir", name).Debug("Pump terminated")
			return
		}
		if err := dst.WriteBuffer(buf); err != nil {
			t.logger.WithError(err).WithField("dir", name).Debug("Pump write failed")
			return
		}
	}
}

// runProxy bridges an established peripheral connection (peer↔local) to
// a freshly-created central connection (local↔target). The forwarding
// strategy depends on whether a static PIN was given:
//
//   - pin < 0: forward whole L2CAP frames (transparent — works only
//     for plaintext links).
//   - pin >= 0: terminate SMP locally on each leg using the static
//     passcode and forward only ATT through both encrypted links.
func (t *teleporter) runProxy(localConn *bleconnecter.BLEConnection) {
	defer localConn.Close()

	target := t.snapshotTargetAddr()
	t.logger.WithFields(logrus.Fields{
		"peer":   localConn.RemoteAddr(),
		"target": target,
	}).Info("Peer connected; opening upstream connection")

	connCtx, connCancel := context.WithTimeout(t.ctx, 10*time.Second)
	remoteConn, _, err := t.remoteStack.BLEConnecter.Connect(connCtx, true, []bleutil.BLEAddr{target},
		bleconnecter.BLEConnectionParametersRequested{
			ConnectionIntervalMin: 6,
			ConnectionIntervalMax: 15,
			ConnectionLatency:     0,
			SupervisionTimeout:    100,
		})
	connCancel()
	if err != nil {
		t.logger.WithError(err).Warn("Failed to connect to target; closing peer")
		return
	}
	defer remoteConn.Close()

	bridgeCtx, bridgeCancel := context.WithCancel(t.ctx)
	defer bridgeCancel()

	if t.pin < 0 {
		t.bridgeRaw(bridgeCtx, bridgeCancel, localConn, remoteConn)
	} else {
		t.bridgePaired(bridgeCtx, bridgeCancel, localConn, remoteConn)
	}
	t.logger.Info("Bridge torn down")
}

// bridgeRaw forwards complete L2CAP frames in both directions. ATT,
// SMP, signalling — everything is passed through without
// interpretation.
func (t *teleporter) bridgeRaw(ctx context.Context, cancel context.CancelFunc, localConn, remoteConn hciconnmgr.BufferConn) {
	t.logger.Info("Both legs up; bridging L2CAP frames")

	var wg sync.WaitGroup
	wg.Add(2)
	run := func(name string, src, dst hciconnmgr.BufferConn) {
		defer wg.Done()
		defer cancel()
		t.pump(ctx, name, src, dst)
	}
	go run("peer→target", localConn, remoteConn)
	go run("target→peer", remoteConn, localConn)
	wg.Wait()
}

// bridgePaired terminates SMP locally on each leg with the static
// passcode and pumps only ATT (CID 4) frames between the two encrypted
// links. The remote (central) side initiates pairing right away
// (SecureOnConnect); the local (peripheral) side responds when the
// peer asks.
func (t *teleporter) bridgePaired(ctx context.Context, cancel context.CancelFunc, localConn, remoteConn *bleconnecter.BLEConnection) {
	t.logger.Info("Both legs up; pairing each leg before bridging ATT")

	// The two legs have asymmetric pairing roles, so they need
	// asymmetric IO capabilities:
	//
	//   - Remote leg (we are central toward the real target): the
	//     target advertises DisplayOnly (its PIN is on a sticker). We
	//     must input that PIN. Advertise KeyboardOnly — only
	//     InputNumeric set — so getLegacyAlgorithmType picks
	//     "initiator inputs, responder displays."
	//
	//   - Local leg (the phone connects to us, we are the cloned
	//     "device"): the phone should prompt the user to enter the
	//     PIN it has read off the real device's sticker. Advertise
	//     DisplayOnly — only DisplayNumeric set — so the negotiation
	//     picks "responder displays, initiator inputs." We don't
	//     actually need a display: StaticPasscode makes
	//     CryptoGeneratePassKey hand back the operator-supplied PIN,
	//     and DisplayNumeric just no-ops because nobody on this side
	//     reads it back.
	pin := uint32(t.pin)
	remoteSMP := &blesmp.SMPConnConfig{
		StaticPasscode:  t.pin,
		InputNumeric:    func(*blesmp.SMPConn) (uint32, error) { return pin, nil },
		AuthReq:         0x05, // bonding + MITM, legacy (no SC)
		MinKeySize:      7,    // accept short legacy keys; many devices still use them
		SecureOnConnect: true,
	}
	localSMP := &blesmp.SMPConnConfig{
		StaticPasscode:  t.pin,
		DisplayNumeric:  func(*blesmp.SMPConn, uint32) error { return nil },
		AuthReq:         0x05,
		MinKeySize:      7,
		SecureOnConnect: false,
	}

	type attBridge struct {
		mu        sync.Mutex
		local     hciconnmgr.BufferConn
		remote    hciconnmgr.BufferConn
		ready     chan struct{}
		readyOnce sync.Once
	}
	bridge := &attBridge{ready: make(chan struct{})}
	maybeSignalReady := func() {
		if bridge.local != nil && bridge.remote != nil {
			bridge.readyOnce.Do(func() { close(bridge.ready) })
		}
	}

	localL2 := blel2cap.New(localConn, nil, func(psm blel2cap.PSMType, accept blel2cap.L2CAPConnAccepter) {
		switch psm {
		case blel2cap.PSMTypeSecurityManager:
			sc := t.localStack.SMP.AddConn(accept(), localSMP)
			if sc != nil {
				/* The library's auto-GoSecure on AddConn is gated to
				   central. Real Victron-style peripherals send an SMP
				   Security Request as soon as the SMP channel opens so the
				   phone is prompted to encrypt (or pair fresh). Mimic that
				   here, otherwise the phone happily uses an unauthenticated
				   link and never asks for the static PIN. */
				go sc.GoSecure(ctx, true)
			}
		case blel2cap.PSMTypeATT:
			bridge.mu.Lock()
			bridge.local = accept()
			maybeSignalReady()
			bridge.mu.Unlock()
		}
	})

	remoteL2 := blel2cap.New(remoteConn, nil, func(psm blel2cap.PSMType, accept blel2cap.L2CAPConnAccepter) {
		switch psm {
		case blel2cap.PSMTypeSecurityManager:
			t.remoteStack.SMP.AddConn(accept(), remoteSMP)
		case blel2cap.PSMTypeATT:
			bridge.mu.Lock()
			bridge.remote = accept()
			maybeSignalReady()
			bridge.mu.Unlock()
		}
	})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := localL2.Run(); err != nil {
			t.logger.WithError(err).Debug("local L2CAP terminated")
		}
	}()
	go func() {
		defer wg.Done()
		defer cancel()
		if err := remoteL2.Run(); err != nil {
			t.logger.WithError(err).Debug("remote L2CAP terminated")
		}
	}()

	select {
	case <-bridge.ready:
	case <-ctx.Done():
		localL2.Close()
		remoteL2.Close()
		wg.Wait()
		return
	}

	t.logger.Info("ATT channels open on both legs; pumping ATT frames")

	var pumpWg sync.WaitGroup
	pumpWg.Add(2)
	runPump := func(name string, src, dst hciconnmgr.BufferConn) {
		defer pumpWg.Done()
		defer cancel()
		t.pump(ctx, name, src, dst)
	}
	go runPump("peer→target ATT", bridge.local, bridge.remote)
	go runPump("target→peer ATT", bridge.remote, bridge.local)

	pumpWg.Wait()
	localL2.Close()
	remoteL2.Close()
	wg.Wait()
}

func (t *teleporter) Run(ready func()) error {
	defer t.cancel()

	t.scanCBHandle = t.remoteStack.BLEScanner.RegisterAdvertisingReportCallback(t.captureAdvReport)
	defer t.remoteStack.BLEScanner.UnregisterAdvertisingReportCallback(t.scanCBHandle)

	// Use the advertiser's BASE slot so BLEConnecter's
	// LegacyAdvertisingSetConnection (which only flips Active/Type and
	// preserves the rest) keeps our cloned bytes when it makes the
	// advertisement connectable.
	t.advSlot = t.localStack.BLEAdvertiser.LegacyAdvertisingGetBaseSlot()

	ready()

	for {
		t.logger.WithField("target", t.targetMac).Info("Waiting for advertising data from target")
		if err := t.waitForAdvData(); err != nil {
			return err
		}
		if err := t.applyAdvData(); err != nil {
			t.logger.WithError(err).Warn("Failed to write cloned data; retrying")
			if !t.sleepCtx(time.Second) {
				return nil
			}
			continue
		}

		t.logger.Info("Cloned beacon ready; advertising for peer connection")
		// Connect blocks in peripheral mode until a matching peer
		// connects. Empty allowedPeers means "any peer".
		conn, _, err := t.localStack.BLEConnecter.Connect(t.ctx, false, t.allowedPeers,
			bleconnecter.BLEConnectionParametersRequested{
				ConnectionIntervalMin: 6,
				ConnectionIntervalMax: 24,
				ConnectionLatency:     0,
				SupervisionTimeout:    200,
			})
		if err != nil {
			if t.ctx.Err() != nil {
				return nil
			}
			t.logger.WithError(err).Warn("Accept failed; retrying")
			if !t.sleepCtx(time.Second) {
				return nil
			}
			continue
		}

		t.runProxy(conn)
	}
}

func (t *teleporter) Close() error {
	t.cancel()
	return nil
}

// waitForAdvData blocks until the scanner has captured at least the
// connectable advertisement, then waits a short additional grace period
// for the scan response so we usually clone both.
func (t *teleporter) waitForAdvData() error {
	for {
		_, _, ok := t.snapshotAdv()
		if ok {
			break
		}
		select {
		case <-t.advUpdated:
		case <-t.ctx.Done():
			return t.ctx.Err()
		}
	}

	graceTimer := time.NewTimer(2 * time.Second)
	defer graceTimer.Stop()
	for {
		t.advMu.Lock()
		haveBoth := t.haveBeacon && t.haveScanRsp
		t.advMu.Unlock()
		if haveBoth {
			return nil
		}
		select {
		case <-t.advUpdated:
		case <-graceTimer.C:
			return nil
		case <-t.ctx.Done():
			return t.ctx.Err()
		}
	}
}

func (t *teleporter) sleepCtx(d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-t.ctx.Done():
		return false
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func parseAddrList(s string) ([]bleutil.BLEAddr, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	out := make([]bleutil.BLEAddr, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		addr, err := bleutil.MacAddrFromString(p)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %w", p, err)
		}
		out = append(out, bleutil.BLEAddr{MacAddr: addr, MacAddrType: 0})
	}
	return out, nil
}

// teleportKeysPath returns a stack-specific path for the SMP key
// store. The two stacks must not share a file: each one writes the
// whole map back on every save, so concurrent writes would clobber
// each other. role is "local" or "remote".
func teleportKeysPath(role string) string {
	dir, err := os.UserConfigDir()
	if err != nil || dir == "" {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "go-ble", "teleport-"+role+"-smp.gob")
}

func openHCI(name, snoopPath string) (hciinterface.HCIInterface, error) {
	dev, err := hcidrivers.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening %q: %w", name, err)
	}
	if snoopPath != "" {
		dev, err = btsnoop.WrapFile(dev, snoopPath)
		if err != nil {
			return nil, fmt.Errorf("wrapping %q with btsnoop %q: %w", name, snoopPath, err)
		}
	}
	return dev, nil
}

func main() {
	targetStr := flag.String("target", "", "MAC of the target BLE device to teleport")
	peerStr := flag.String("peer", "", "Optional comma-separated MAC whitelist for clients connecting to the cloned advertiser. Empty = accept any peer.")
	pin := flag.Int("pin", -1, "Static legacy passcode (0..999999) used for SMP pairing on both legs. -1 disables pairing.")
	snoopLocal := flag.String("btsnoop-local", "", "Write btsnoop file for the local controller")
	snoopRemote := flag.String("btsnoop-remote", "", "Write btsnoop file for the remote controller")

	bleutilparam.Init()
	logrusconfig.InitParam()
	flag.Parse()

	logger := logrusconfig.GetLogger(0)

	if *targetStr == "" {
		logger.Fatalln("-target is required")
	}
	targetMac, err := bleutil.MacAddrFromString(*targetStr)
	if err != nil {
		logger.Fatalf("Bad -target: %v", err)
	}
	allowed, err := parseAddrList(*peerStr)
	if err != nil {
		logger.Fatalf("Bad -peer: %v", err)
	}
	if *pin > 999999 {
		logger.Fatalln("-pin must be in 0..999999")
	}

	localDev, err := bleutilparam.GetDeviceNameMulti(0)
	if err != nil {
		return
	}
	remoteDev, err := bleutilparam.GetDeviceNameMulti(1)
	if err != nil {
		logger.Fatalln("Two devices are required: -device <local>,<remote> (the second is typically `net` with -netaddr <ble2net-host>:<port>)")
	}

	localHCI, err := openHCI(localDev, *snoopLocal)
	if err != nil {
		logger.Fatalln(err)
	}
	remoteHCI, err := openHCI(remoteDev, *snoopRemote)
	if err != nil {
		logger.Fatalln(err)
	}

	// Local stack: advertiser + connecter. AlwaysAdvertising=false keeps
	// the base slot Active=false until BLEConnecter.Connect flips it.
	// We write the cloned beacon bytes into that same base slot, so
	// SetConnection's Active/Type toggle preserves the cloned content.
	localCfg := ble.DefaultConfig()
	localCfg.BLEScannerUse = false
	localCfg.BLEAdvertiserUse = true
	localCfg.BLEConnecterUse = true
	localCfg.BLEAdvertiserConfig = bleadvertiser.DefaultConfig()
	localCfg.BLEAdvertiserConfig.AlwaysAdvertising = false
	localCfg.HCIControllerConfig.PrivacyAdvertise = false
	localCfg.SMPConfig = blesmp.DefaultConfig()
//	localCfg.SMPConfig.StoredKeysPath = teleportKeysPath("local")

	localStack := ble.New(bleutil.LogWithPrefix(logger, "local"), localCfg, localHCI)
	if localStack == nil {
		logger.Fatalln("Could not create local stack")
	}

	// Remote stack: scanner + connecter, no advertiser.
	remoteCfg := ble.DefaultConfig()
	remoteCfg.BLEScannerUse = true
	remoteCfg.BLEAdvertiserUse = false
	remoteCfg.BLEConnecterUse = true
	remoteCfg.BLEScannerConfig = &blescanner.BLEScannerConfig{
		ScanCycleActiveDuty: 1, // always actively scanning so we capture scan responses
		LEScanInterval:      0x40,
		LEScanWindow:        0x30,
	}
	remoteCfg.HCIControllerConfig.PrivacyConnect = false
	remoteCfg.HCIControllerConfig.PrivacyScan = false
	remoteCfg.SMPConfig = blesmp.DefaultConfig()
//	remoteCfg.SMPConfig.StoredKeysPath = teleportKeysPath("remote")

	remoteStack := ble.New(bleutil.LogWithPrefix(logger, "remote"), remoteCfg, remoteHCI)
	if remoteStack == nil {
		logger.Fatalln("Could not create remote stack")
	}

	ctx, cancel := context.WithCancel(context.Background())
	t := &teleporter{
		logger:       bleutil.LogWithPrefix(logger, "tele"),
		targetMac:    targetMac,
		allowedPeers: allowed,
		pin:          int32(*pin),
		localStack:   localStack,
		remoteStack:  remoteStack,
		advUpdated:   make(chan struct{}, 1),
		ctx:          ctx,
		cancel:       cancel,
	}

	m := multirun.MultiRun{}
	m.HandleSIGTERM()
	m.RegisterRunnableReady(localStack)
	m.RegisterRunnableReady(remoteStack)
	m.RegisterRunnableReady(t)

	logger.WithFields(logrus.Fields{
		"target": t.targetMac,
		"peers":  t.allowedPeers,
		"local":  localDev,
		"remote": remoteDev,
	}).Info("Starting teleport")
	logger.Fatalln(m.Run(func() {
		logger.Info("Ready!")
	}))
}
