package loopback

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	ble "github.com/BertoldVdb/go-ble"
	"github.com/BertoldVdb/go-ble/bleadvertiser"
	"github.com/BertoldVdb/go-ble/bleatt"
	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	"github.com/BertoldVdb/go-ble/bleconnecter"
	"github.com/BertoldVdb/go-ble/blescanner"
	"github.com/BertoldVdb/go-ble/blesmp"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	blel2cap "github.com/BertoldVdb/go-ble/l2cap"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

// silentLogger returns a logrus.Entry that swallows everything; the
// loopback driver and the stack are quite chatty otherwise.
func silentLogger() *logrus.Entry {
	l := logrus.New()
	l.Out = io.Discard
	l.Level = logrus.PanicLevel
	return logrus.NewEntry(l)
}

// debugLogger is used during development to surface what the stack /
// loopback are doing. Tests should use silentLogger().
func debugLogger() *logrus.Entry {
	l := logrus.New()
	l.Level = logrus.TraceLevel
	return logrus.NewEntry(l)
}

var _ = debugLogger

// startStack wires a BluetoothStack to a loopback endpoint and runs it.
// Returns the stack and a stop function; the stack is fully ready
// (advertiser/scanner/connecter all configured) when this returns.
func startStack(t *testing.T, ep *Endpoint, scannerOn, advertiserOn bool) (*ble.BluetoothStack, func()) {
	t.Helper()

	cfg := ble.DefaultConfig()
	cfg.BLEScannerUse = scannerOn
	cfg.BLEAdvertiserUse = advertiserOn
	cfg.BLEConnecterUse = true
	cfg.BLEAdvertiserConfig = bleadvertiser.DefaultConfig()
	cfg.BLEAdvertiserConfig.AlwaysAdvertising = false
	cfg.BLEScannerConfig = &blescanner.BLEScannerConfig{
		StoreGAPMap:         true,
		ScanCycleDurationMs: 100,
		ScanCycleActiveDuty: 1, // continuous active scanning
	}
	// Disable LTK persistence path for tests by directing it nowhere.
	cfg.SMPConfig = blesmp.DefaultConfig()
	cfg.SMPConfig.StoredKeysPath = ""
	cfg.HCIControllerConfig.WatchdogTimeout = 0
	cfg.HCIControllerConfig.AwaitStartup = false
	cfg.HCIControllerConfig.PrivacyAdvertise = false
	cfg.HCIControllerConfig.PrivacyConnect = false
	cfg.HCIControllerConfig.PrivacyScan = false
	cfg.HCIControllerConfig.LERandomAddrBits = 32

	stack := ble.New(silentLogger(), cfg, ep)

	ready := make(chan struct{})
	go func() {
		stack.Run(func() { close(ready) })
	}()

	select {
	case <-ready:
	case <-time.After(2 * time.Second):
		t.Fatal("stack did not become ready within 2s")
	}

	return stack, func() {
		stack.Close()
	}
}

// Sanity test: two stacks bring up cleanly through the loopback.
func TestLoopbackBringup(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	stackA, stopA := startStack(t, a, true, false)
	defer stopA()
	stackB, stopB := startStack(t, b, false, false)
	defer stopB()

	if stackA == nil || stackB == nil {
		t.Fatal("stack init failed")
	}
	// Smoke: both controllers should have read their BD_ADDR.
	if stackA.Controller.Info.BdAddr == nil {
		t.Error("stack A: BD_ADDR not populated")
	}
	if stackB.Controller.Info.BdAddr == nil {
		t.Error("stack B: BD_ADDR not populated")
	}
}

// Two stacks: B advertises with a unique service UUID, A scans and
// discovers it.
func TestLoopbackScanFindsAdvertiser(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	advUUID := bleutil.UUIDFromStringPanic("18ff")

	cfgB := ble.DefaultConfig()
	cfgB.BLEScannerUse = false
	cfgB.BLEAdvertiserUse = true
	cfgB.BLEAdvertiserConfig = bleadvertiser.DefaultConfig()
	cfgB.BLEAdvertiserConfig.AlwaysAdvertising = true
	cfgB.BLEAdvertiserConfig.DeviceName = "B"
	cfgB.BLEAdvertiserConfig.DeviceService = advUUID
	cfgB.SMPConfig = blesmp.DefaultConfig()
	cfgB.SMPConfig.StoredKeysPath = ""
	cfgB.HCIControllerConfig.WatchdogTimeout = 0
	cfgB.HCIControllerConfig.AwaitStartup = false
	cfgB.HCIControllerConfig.PrivacyAdvertise = false
	cfgB.HCIControllerConfig.PrivacyConnect = false
	cfgB.HCIControllerConfig.PrivacyScan = false
	cfgB.HCIControllerConfig.LERandomAddrBits = 32

	stackB := ble.New(silentLogger(), cfgB, b)
	bReady := make(chan struct{})
	go func() { stackB.Run(func() { close(bReady) }) }()
	defer stackB.Close()
	<-bReady

	// Now A: scanner-only.
	stackA, stopA := startStack(t, a, true, false)
	defer stopA()

	// Wait up to 3s for the scanner to see B's advertisement (with the
	// distinguishing service UUID we set).
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		known := stackA.BLEScanner.KnownDevicesAddresses(nil)
		for _, addr := range known {
			dev := stackA.BLEScanner.GetDevice(addr)
			if dev == nil {
				continue
			}
			services := dev.GetServices(0, nil)
			for _, u := range services {
				if u == advUUID {
					return // success
				}
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("scanner did not discover advertiser within 3s")
}

// Full-stack loopback: A (central) and B (peripheral) both run their
// connecters; B has a pending Connect(isCentral=false), A initiates
// Connect(isCentral=true), they meet in the middle. Exercises every
// layer (HCI command + event encoding, ACL forwarding, L2CAP).
func TestLoopbackConnectAndDisconnect(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	// B is a peripheral; A is a central. Both have advertiser+connecter
	// running. The loopback assigns B's bdAddr to 22:22:22:22:22:22 and
	// A's to 11:11:11:11:11:11.
	stackB, stopB := startStack(t, b, false, true)
	defer stopB()
	stackA, stopA := startStack(t, a, true, false)
	defer stopA()

	aAddr := bleutil.BLEAddr{
		MacAddr:     0x111111111111,
		MacAddrType: bleutil.MacAddrPublic,
	}
	bAddr := bleutil.BLEAddr{
		MacAddr:     0x222222222222,
		MacAddrType: bleutil.MacAddrPublic,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type connResult struct {
		conn *bleconnecter.BLEConnection
		err  error
	}
	bDone := make(chan connResult, 1)
	go func() {
		c, _, err := stackB.BLEConnecter.Connect(ctx, false, []bleutil.BLEAddr{aAddr}, bleconnecter.BLEConnectionParametersRequested{})
		bDone <- connResult{c, err}
	}()

	// Give B a moment to enable connectable advertising.
	time.Sleep(100 * time.Millisecond)

	connA, _, err := stackA.BLEConnecter.Connect(ctx, true, []bleutil.BLEAddr{bAddr}, bleconnecter.BLEConnectionParametersRequested{})
	if err != nil {
		t.Fatalf("A Connect: %v", err)
	}
	if connA == nil {
		t.Fatal("A Connect returned nil")
	}

	select {
	case r := <-bDone:
		if r.err != nil {
			t.Fatalf("B Connect: %v", r.err)
		}
		if r.conn == nil {
			t.Fatal("B Connect returned nil")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("B Connect did not complete")
	}

	handle := connA.Connection.GetHandle()
	if got := stackB.Controller.ConnMgr.FindConnectionByHandle(handle); got == nil {
		t.Errorf("B ConnMgr missing handle %#x", handle)
	}

	if err := connA.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}

	deadline := time.Now().Add(1 * time.Second)
	for time.Now().Before(deadline) {
		if stackB.Controller.ConnMgr.FindConnectionByHandle(handle) == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if stackB.Controller.ConnMgr.FindConnectionByHandle(handle) != nil {
		t.Errorf("B ConnMgr still has handle %#x after A disconnected", handle)
	}

	_ = hciconnmgr.ErrorClosed
}

// Two stacks linked through the loopback exchange a GATT notification
// end-to-end:
//
//	1. B exposes a service/characteristic with the Notify flag.
//	2. A and B connect (central + peripheral).
//	3. Each side runs L2CAP over its BLEConnection so the ATT channel
//	   is visible.
//	4. A discovers B's structure, then Subscribes to the characteristic
//	   (which writes the CCCD).
//	5. B calls SetValue on its local characteristic, which goes through
//	   characteristicNotify → ATTHandleValueNTF → loopback ACL forward.
//	6. A's ClientNotifyHandler receives the new value.
//
// This exercises HCI command/event encoding, ACL forwarding, L2CAP
// signalling+data, ATT MTU exchange, ATT discovery, ATT Write Req/Rsp
// for the CCCD, and ATT Handle-Value-Notification — the full GATT
// notify path through every layer.
func TestLoopbackGATTNotification(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	stackA, stopA := startStack(t, a, false /* no scanner */, false)
	defer stopA()
	stackB, stopB := startStack(t, b, false, true /* advertiser */)
	defer stopB()

	aAddr := bleutil.BLEAddr{MacAddr: 0x111111111111, MacAddrType: bleutil.MacAddrPublic}
	bAddr := bleutil.BLEAddr{MacAddr: 0x222222222222, MacAddrType: bleutil.MacAddrPublic}

	// ---- Build B's GATT structure: one service with one notifiable +
	// readable characteristic. ----
	bStruct := attstructure.NewStructure()
	bSvc := bStruct.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	bChar := bSvc.AddCharacteristic(
		bleutil.UUIDFromStringPanic("2a3a"),
		attstructure.CharacteristicRead|attstructure.CharacteristicNotify,
		attstructure.ValueConfig{LengthMax: 32},
	)

	// Initial connect ------------------------------------------------
	// Generous deadline so the test holds up under -race -count=N
	// where the scheduler is heavily contested.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	type res struct {
		conn *bleconnecter.BLEConnection
		err  error
	}
	bDone := make(chan res, 1)
	go func() {
		c, _, err := stackB.BLEConnecter.Connect(ctx, false, []bleutil.BLEAddr{aAddr}, bleconnecter.BLEConnectionParametersRequested{})
		bDone <- res{c, err}
	}()
	time.Sleep(100 * time.Millisecond)

	connA, _, err := stackA.BLEConnecter.Connect(ctx, true, []bleutil.BLEAddr{bAddr}, bleconnecter.BLEConnectionParametersRequested{})
	if err != nil {
		t.Fatalf("A Connect: %v", err)
	}
	bRes := <-bDone
	if bRes.err != nil {
		t.Fatalf("B Connect: %v", bRes.err)
	}
	connB := bRes.conn

	// ---- Wire L2CAP + ATT on both sides ----
	// Central side: GattDevice will discover the remote structure.
	aGatt := bleatt.NewGattDeviceWithConn(connA, attstructure.NewStructure(), &bleatt.GattDeviceConfig{
		MTU:                     247,
		DeviceName:              "A",
		DiscoverRemoteOnConnect: true,
	})
	aL2 := blel2cap.New(connA, nil, func(psm blel2cap.PSMType, accept blel2cap.L2CAPConnAccepter) {
		switch psm {
		case blel2cap.PSMTypeATT:
			aGatt.AddConn(accept())
		}
	})
	go aL2.Run()
	defer aL2.Close()

	// Peripheral side: GattDevice exposes B's structure.
	bGatt := bleatt.NewGattDevice(bStruct, &bleatt.GattDeviceConfig{
		MTU:                     247,
		DeviceName:              "B",
		DiscoverRemoteOnConnect: false,
	})
	bL2 := blel2cap.New(connB, nil, func(psm blel2cap.PSMType, accept blel2cap.L2CAPConnAccepter) {
		switch psm {
		case blel2cap.PSMTypeATT:
			bGatt.AddConn(accept())
		}
	})
	go bL2.Run()
	defer bL2.Close()

	// ---- Wait for A to discover B's structure ----
	// Give the discovery a generous timeout so the test holds up even
	// under heavy parallel test runs (-race -count=N).
	discoverCtx, discoverCancel := context.WithTimeout(ctx, 10*time.Second)
	defer discoverCancel()
	remote := aGatt.ClientGetStructure(discoverCtx)
	if remote == nil {
		t.Fatal("A did not discover B's structure")
	}
	rSvc := remote.GetService(bleutil.UUIDFromStringPanic("180a"))
	if rSvc == nil {
		t.Fatal("A did not discover the 180a service")
	}
	rChar := rSvc.GetCharacteristic(bleutil.UUIDFromStringPanic("2a3a"))
	if rChar == nil {
		t.Fatal("A did not discover the 2a3a characteristic")
	}

	// ---- A subscribes to notifications ----
	got := make(chan []byte, 4)
	if err := rChar.Subscribe(ctx, func(value []byte) {
		// Snapshot the value (the buffer may be reused by the parser).
		select {
		case got <- append([]byte(nil), value...):
		default:
		}
	}); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	// ---- B fires a notification by writing its local value ----
	want := []byte("hello-loopback")
	if _, err := bChar.SetValue(ctx, want); err != nil {
		t.Fatalf("B SetValue: %v", err)
	}

	select {
	case rx := <-got:
		if string(rx) != string(want) {
			t.Errorf("notification payload: got %q want %q", rx, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("A did not receive notification within 2s")
	}

	// Send a second notification to confirm we keep delivering.
	want2 := []byte("again")
	if _, err := bChar.SetValue(ctx, want2); err != nil {
		t.Fatalf("B SetValue (2): %v", err)
	}
	select {
	case rx := <-got:
		if string(rx) != string(want2) {
			t.Errorf("second notification payload: got %q want %q", rx, want2)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("A did not receive second notification within 2s")
	}

	// Cleanup: A initiates disconnect.
	if err := connA.Close(); err != nil {
		t.Errorf("A Close: %v", err)
	}
}

// connectStacksWithParams brings up A as central + B as peripheral and
// returns both BLEConnections. The peripheral-side auto-UpdateParams
// (which bleconnecter.Connect always issues) uses bParams; pass values
// distinct from your test's user-driven UpdateParams so you can wait
// for the auto-update to settle before perturbing them.
func connectStacksWithParams(
	t *testing.T,
	ctx context.Context,
	stackA, stackB *ble.BluetoothStack,
	bParams bleconnecter.BLEConnectionParametersRequested,
) (connA, connB *bleconnecter.BLEConnection) {
	t.Helper()
	aAddr := bleutil.BLEAddr{MacAddr: 0x111111111111, MacAddrType: bleutil.MacAddrPublic}
	bAddr := bleutil.BLEAddr{MacAddr: 0x222222222222, MacAddrType: bleutil.MacAddrPublic}

	type res struct {
		conn *bleconnecter.BLEConnection
		err  error
	}
	bDone := make(chan res, 1)
	go func() {
		c, _, err := stackB.BLEConnecter.Connect(ctx, false, []bleutil.BLEAddr{aAddr}, bParams)
		bDone <- res{c, err}
	}()
	time.Sleep(100 * time.Millisecond)

	c, _, err := stackA.BLEConnecter.Connect(ctx, true, []bleutil.BLEAddr{bAddr}, bleconnecter.BLEConnectionParametersRequested{})
	if err != nil {
		t.Fatalf("A Connect: %v", err)
	}
	connA = c

	r := <-bDone
	if r.err != nil {
		t.Fatalf("B Connect: %v", r.err)
	}
	connB = r.conn
	return
}

// connectStacks is the simple variant — empty initial params on B.
func connectStacks(t *testing.T, ctx context.Context, stackA, stackB *ble.BluetoothStack) (connA, connB *bleconnecter.BLEConnection) {
	return connectStacksWithParams(t, ctx, stackA, stackB, bleconnecter.BLEConnectionParametersRequested{})
}

// waitForInterval polls until both endpoints' parametersActual report
// the given interval, or the deadline elapses. Tests use this to
// rendezvous on the peripheral's auto-update settling before doing
// their own UpdateParams.
func waitForInterval(t *testing.T, connA, connB *bleconnecter.BLEConnection, interval uint16, timeout time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		pa := connA.GetActualParameters()
		pb := connB.GetActualParameters()
		if pa.Interval == interval && pb.Interval == interval {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// Central-initiated UpdateParams: A calls UpdateParams, the loopback
// emits LE Connection Update Complete on both sides, the connecter's
// handler updates parametersActual on both sides accordingly.
func TestLoopbackConnectionUpdate_CentralInitiated(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	stackA, stopA := startStack(t, a, false, false)
	defer stopA()
	stackB, stopB := startStack(t, b, false, true)
	defer stopB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// bleconnecter.Connect on the peripheral side automatically issues
	// an UpdateParams after the link comes up. We pass a distinct interval
	// here so we can rendezvous on it settling before the test performs
	// its own UpdateParams (otherwise the auto-update may arrive after
	// our update and overwrite it under race ordering).
	bInitial := bleconnecter.BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x30,
		ConnectionIntervalMax: 0x30,
		ConnectionLatency:     0,
		SupervisionTimeout:    0x80,
	}
	connA, connB := connectStacksWithParams(t, ctx, stackA, stackB, bInitial)
	if !waitForInterval(t, connA, connB, 0x30, 2*time.Second) {
		t.Fatal("auto-update did not settle to the peripheral's requested interval")
	}

	// Issue an UpdateParams from the central side. The loopback assigns
	// `IntervalMax` as the new interval and emits LEConnectionUpdateComplete
	// on both sides.
	want := bleconnecter.BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x40,
		ConnectionIntervalMax: 0x40,
		ConnectionLatency:     0x05,
		SupervisionTimeout:    0x80,
	}
	if err := connA.UpdateParams(want); err != nil {
		t.Fatalf("UpdateParams: %v", err)
	}

	// Wait for both sides to observe the new parameters via the
	// LEConnectionUpdateComplete event.
	deadline := time.Now().Add(2 * time.Second)
	checkSide := func(label string, c *bleconnecter.BLEConnection) {
		for time.Now().Before(deadline) {
			actual := c.GetActualParameters()
			if actual.Interval == want.ConnectionIntervalMax &&
				actual.Latency == want.ConnectionLatency &&
				actual.Timeout == want.SupervisionTimeout {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
		actual := c.GetActualParameters()
		t.Errorf("%s did not see updated params: got %+v want interval=%d latency=%d timeout=%d",
			label, actual, want.ConnectionIntervalMax, want.ConnectionLatency, want.SupervisionTimeout)
	}
	checkSide("A", connA)
	checkSide("B", connB)
}

// Peripheral-initiated update: B calls UpdateParams; the loopback fires
// LE Remote Connection Parameter Request on A; A's connecter handler
// (with the default verify callback being nil → accept) replies with
// LERemoteConnectionParameterRequestReply, which the loopback turns
// into LEConnectionUpdateComplete on both sides.
func TestLoopbackConnectionUpdate_PeripheralInitiated(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	// A side: provide a verify callback that *accepts* any request, so
	// the connecter explicitly takes the accept path.
	cfgA := ble.DefaultConfig()
	cfgA.BLEScannerUse = false
	cfgA.BLEAdvertiserUse = false
	cfgA.BLEConnecterConfig = &bleconnecter.BLEConnecterConfig{
		BLEUpdateParametersVerify: func(c *bleconnecter.BLEConnection, intervalMin, intervalMax, latency, timeout uint16) bool {
			return true
		},
	}
	cfgA.SMPConfig = blesmp.DefaultConfig()
	cfgA.SMPConfig.StoredKeysPath = ""
	cfgA.HCIControllerConfig.WatchdogTimeout = 0
	cfgA.HCIControllerConfig.AwaitStartup = false
	cfgA.HCIControllerConfig.PrivacyAdvertise = false
	cfgA.HCIControllerConfig.PrivacyConnect = false
	cfgA.HCIControllerConfig.PrivacyScan = false
	cfgA.HCIControllerConfig.LERandomAddrBits = 32
	stackA := ble.New(silentLogger(), cfgA, a)
	aReady := make(chan struct{})
	go func() { stackA.Run(func() { close(aReady) }) }()
	defer stackA.Close()
	<-aReady

	stackB, stopB := startStack(t, b, false, true)
	defer stopB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bInitial := bleconnecter.BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x30,
		ConnectionIntervalMax: 0x30,
		ConnectionLatency:     0,
		SupervisionTimeout:    0x80,
	}
	connA, connB := connectStacksWithParams(t, ctx, stackA, stackB, bInitial)
	if !waitForInterval(t, connA, connB, 0x30, 2*time.Second) {
		t.Fatal("auto-update did not settle")
	}

	// B (peripheral) requests new parameters. In real BLE this would
	// turn into an LL_CONNECTION_PARAM_REQ, surfaced to the central as
	// LE Remote Connection Parameter Request. Our loopback mirrors that.
	want := bleconnecter.BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x60,
		ConnectionIntervalMax: 0x60,
		ConnectionLatency:     0x02,
		SupervisionTimeout:    0xA0,
	}
	if err := connB.UpdateParams(want); err != nil {
		t.Fatalf("B UpdateParams: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	checkSide := func(label string, c *bleconnecter.BLEConnection) {
		for time.Now().Before(deadline) {
			actual := c.GetActualParameters()
			if actual.Interval == want.ConnectionIntervalMax &&
				actual.Latency == want.ConnectionLatency &&
				actual.Timeout == want.SupervisionTimeout {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
		actual := c.GetActualParameters()
		t.Errorf("%s did not see updated params: got %+v want interval=%d latency=%d timeout=%d",
			label, actual, want.ConnectionIntervalMax, want.ConnectionLatency, want.SupervisionTimeout)
	}
	checkSide("A", connA)
	checkSide("B", connB)
}

// Peripheral-initiated update with the central rejecting: the central's
// verify callback returns false, the connecter sends NegativeReply, and
// no LEConnectionUpdateComplete is emitted — both sides retain the
// pre-existing parameters.
func TestLoopbackConnectionUpdate_PeripheralInitiated_Rejected(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	cfgA := ble.DefaultConfig()
	cfgA.BLEScannerUse = false
	cfgA.BLEAdvertiserUse = false
	cfgA.BLEConnecterConfig = &bleconnecter.BLEConnecterConfig{
		BLEUpdateParametersVerify: func(c *bleconnecter.BLEConnection, intervalMin, intervalMax, latency, timeout uint16) bool {
			return false // always reject
		},
	}
	cfgA.SMPConfig = blesmp.DefaultConfig()
	cfgA.SMPConfig.StoredKeysPath = ""
	cfgA.HCIControllerConfig.WatchdogTimeout = 0
	cfgA.HCIControllerConfig.AwaitStartup = false
	cfgA.HCIControllerConfig.PrivacyAdvertise = false
	cfgA.HCIControllerConfig.PrivacyConnect = false
	cfgA.HCIControllerConfig.PrivacyScan = false
	cfgA.HCIControllerConfig.LERandomAddrBits = 32
	stackA := ble.New(silentLogger(), cfgA, a)
	aReady := make(chan struct{})
	go func() { stackA.Run(func() { close(aReady) }) }()
	defer stackA.Close()
	<-aReady

	stackB, stopB := startStack(t, b, false, true)
	defer stopB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// In this test the central rejects, so the auto-update from B
	// will be NegativeReply'd. parametersActual stays at the
	// connection-complete values. Use empty initial params on B; the
	// final stable state is whatever the LEConnectionComplete event
	// carried (the loopback uses IntervalMax = 6 from makeValid).
	connA, connB := connectStacks(t, ctx, stackA, stackB)
	// Brief settle for the LEConnectionComplete events to land.
	time.Sleep(200 * time.Millisecond)

	originalA := connA.GetActualParameters()
	originalB := connB.GetActualParameters()

	rejected := bleconnecter.BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x60,
		ConnectionIntervalMax: 0x60,
		ConnectionLatency:     0x02,
		SupervisionTimeout:    0xA0,
	}
	_ = connB.UpdateParams(rejected) // may return without error; the rejection is silent

	// Give the loopback a chance to deliver any pending events.
	time.Sleep(200 * time.Millisecond)

	if got := connA.GetActualParameters(); got != originalA {
		t.Errorf("A side params changed despite rejection: %+v -> %+v", originalA, got)
	}
	if got := connB.GetActualParameters(); got != originalB {
		t.Errorf("B side params changed despite rejection: %+v -> %+v", originalB, got)
	}
}

// Encryption (HCI level): A initiates encryption with a known LTK.
// We register a stub LEEncryptionGetKey on B's connmgr so it returns
// the matching LTK; the loopback compares them, fires EncryptionChange
// on both sides, and both stacks observe the new state via the
// EncryptionChanged hook.
func TestLoopbackEncryption_LTKMatch(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	stackA, stopA := startStack(t, a, false, false)
	defer stopA()
	stackB, stopB := startStack(t, b, false, true)
	defer stopB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connA, _ := connectStacks(t, ctx, stackA, stackB)

	// The shared LTK both sides will use. SC pairing produces LTKs
	// like this; for the test we hard-code one.
	var ltk [16]byte
	for i := range ltk {
		ltk[i] = byte(0x10 + i)
	}

	// Override the SMP callbacks installed by stack.New(). On B we need
	// to return the matching LTK in response to the LongTermKeyRequest
	// event; on both sides we record the EncryptionChanged event so the
	// test can observe the outcome.
	type encChange struct {
		status  uint8
		enabled uint8
	}
	encA := make(chan encChange, 4)
	encB := make(chan encChange, 4)

	stackA.Controller.ConnMgr.SetEventsSMP(hciconnmgr.ConnectionMangerEventsSMP{
		EncryptionChanged: func(conn *hciconnmgr.Connection, ev *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
			select {
			case encA <- encChange{ev.Status, ev.EncryptionEnabled}:
			default:
			}
			return ev
		},
	})
	stackB.Controller.ConnMgr.SetEventsSMP(hciconnmgr.ConnectionMangerEventsSMP{
		LEEncryptionGetKey: func(conn *hciconnmgr.Connection, ev *hcievents.LELongTermKeyRequestEvent) ([]byte, *hcievents.LELongTermKeyRequestEvent) {
			return ltk[:], ev
		},
		EncryptionChanged: func(conn *hciconnmgr.Connection, ev *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
			select {
			case encB <- encChange{ev.Status, ev.EncryptionEnabled}:
			default:
			}
			return ev
		},
	})

	if err := connA.Encrypt(0, 0, ltk); err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	expect := func(t *testing.T, label string, ch chan encChange, status, enabled uint8) {
		t.Helper()
		select {
		case e := <-ch:
			if e.status != status || e.enabled != enabled {
				t.Errorf("%s EncryptionChanged: got {status=%#x enabled=%d} want {status=%#x enabled=%d}",
					label, e.status, e.enabled, status, enabled)
			}
		case <-time.After(2 * time.Second):
			t.Errorf("%s: no EncryptionChanged event received", label)
		}
	}
	expect(t, "A", encA, 0x00, 0x01)
	expect(t, "B", encB, 0x00, 0x01)
}

// Encryption (HCI level): LTK mismatch → both sides see status=0x06.
func TestLoopbackEncryption_LTKMismatch(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	stackA, stopA := startStack(t, a, false, false)
	defer stopA()
	stackB, stopB := startStack(t, b, false, true)
	defer stopB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connA, _ := connectStacks(t, ctx, stackA, stackB)

	var ltkA, ltkB [16]byte
	ltkA[0] = 0xAA
	ltkB[0] = 0xBB

	type encChange struct {
		status  uint8
		enabled uint8
	}
	encA := make(chan encChange, 4)
	encB := make(chan encChange, 4)

	stackA.Controller.ConnMgr.SetEventsSMP(hciconnmgr.ConnectionMangerEventsSMP{
		EncryptionChanged: func(conn *hciconnmgr.Connection, ev *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
			select {
			case encA <- encChange{ev.Status, ev.EncryptionEnabled}:
			default:
			}
			return ev
		},
	})
	stackB.Controller.ConnMgr.SetEventsSMP(hciconnmgr.ConnectionMangerEventsSMP{
		LEEncryptionGetKey: func(conn *hciconnmgr.Connection, ev *hcievents.LELongTermKeyRequestEvent) ([]byte, *hcievents.LELongTermKeyRequestEvent) {
			return ltkB[:], ev // returns a different LTK than A supplied
		},
		EncryptionChanged: func(conn *hciconnmgr.Connection, ev *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
			select {
			case encB <- encChange{ev.Status, ev.EncryptionEnabled}:
			default:
			}
			return ev
		},
	})

	if err := connA.Encrypt(0, 0, ltkA); err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	expect := func(t *testing.T, label string, ch chan encChange) {
		t.Helper()
		select {
		case e := <-ch:
			if e.status != 0x06 || e.enabled != 0x00 {
				t.Errorf("%s EncryptionChanged on mismatch: got {status=%#x enabled=%d} want {status=0x06 enabled=0}",
					label, e.status, e.enabled)
			}
		case <-time.After(2 * time.Second):
			t.Errorf("%s: no EncryptionChanged event received", label)
		}
	}
	expect(t, "A", encA)
	expect(t, "B", encB)
}

// Encryption (HCI level): peripheral has no LTK and replies negatively.
// Both sides see EncryptionChange with status=0x06.
func TestLoopbackEncryption_NegReply(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	stackA, stopA := startStack(t, a, false, false)
	defer stopA()
	stackB, stopB := startStack(t, b, false, true)
	defer stopB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connA, _ := connectStacks(t, ctx, stackA, stackB)

	var ltkA [16]byte
	ltkA[0] = 0x77

	type encChange struct {
		status  uint8
		enabled uint8
	}
	encA := make(chan encChange, 4)
	encB := make(chan encChange, 4)

	stackA.Controller.ConnMgr.SetEventsSMP(hciconnmgr.ConnectionMangerEventsSMP{
		EncryptionChanged: func(conn *hciconnmgr.Connection, ev *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
			select {
			case encA <- encChange{ev.Status, ev.EncryptionEnabled}:
			default:
			}
			return ev
		},
	})
	stackB.Controller.ConnMgr.SetEventsSMP(hciconnmgr.ConnectionMangerEventsSMP{
		LEEncryptionGetKey: func(conn *hciconnmgr.Connection, ev *hcievents.LELongTermKeyRequestEvent) ([]byte, *hcievents.LELongTermKeyRequestEvent) {
			// Returning nil triggers LELongTermKeyRequestNegativeReply
			// from the connmgr.
			return nil, ev
		},
		EncryptionChanged: func(conn *hciconnmgr.Connection, ev *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
			select {
			case encB <- encChange{ev.Status, ev.EncryptionEnabled}:
			default:
			}
			return ev
		},
	})

	if err := connA.Encrypt(0, 0, ltkA); err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	select {
	case e := <-encA:
		if e.status != 0x06 || e.enabled != 0x00 {
			t.Errorf("A EncryptionChanged on neg-reply: got {status=%#x enabled=%d} want {status=0x06 enabled=0}",
				e.status, e.enabled)
		}
	case <-time.After(2 * time.Second):
		t.Error("A: no EncryptionChanged event received")
	}

	select {
	case e := <-encB:
		if e.status != 0x06 || e.enabled != 0x00 {
			t.Errorf("B EncryptionChanged on neg-reply: got {status=%#x enabled=%d} want {status=0x06 enabled=0}",
				e.status, e.enabled)
		}
	case <-time.After(2 * time.Second):
		t.Error("B: no EncryptionChanged event received")
	}
}

// gattWiring wires L2CAP+SMP+ATT on one side of a connection.
type gattWiring struct {
	l2       *blel2cap.L2CAP
	smpConn  *blesmp.SMPConn
	gatt     *bleatt.GattDevice
	smpReady chan struct{}
}

func wireGATTWithSMP(
	stack *ble.BluetoothStack,
	conn *bleconnecter.BLEConnection,
	structure *attstructure.Structure,
	cfg *bleatt.GattDeviceConfig,
	smpCfg *blesmp.SMPConnConfig,
	useGattWithConn bool,
) *gattWiring {
	w := &gattWiring{smpReady: make(chan struct{})}

	if useGattWithConn {
		w.gatt = bleatt.NewGattDeviceWithConn(conn, structure, cfg)
	} else {
		w.gatt = bleatt.NewGattDevice(structure, cfg)
	}

	w.l2 = blel2cap.New(conn, nil, func(psm blel2cap.PSMType, accept blel2cap.L2CAPConnAccepter) {
		switch psm {
		case blel2cap.PSMTypeSecurityManager:
			w.smpConn = stack.SMP.AddConn(accept(), smpCfg)
			w.gatt.SetSMP(w.smpConn)
			close(w.smpReady)
		case blel2cap.PSMTypeATT:
			w.gatt.AddConnWithSMP(accept(), w.smpConn)
		}
	})
	return w
}

func (w *gattWiring) startL2() { go w.l2.Run() }

// Full-stack loopback with encryption and key caching:
//
//  1. Build two stacks with SMP key persistence pointing at temp files.
//  2. Session 1: connect, wire L2CAP+SMP+ATT, central initiates pairing
//     via GoSecure, both sides reach StateSecure (LE SC Just Works),
//     LTK gets stored in both SMP key stores. Subscribe + receive
//     notification over the encrypted link. Disconnect.
//  3. Session 2: new connection. SMP.AddConn on the central side
//     consults the cached LTK and immediately calls leEncrypt, so the
//     link comes up already encrypted without any pairing exchange.
//     Verify SMP state is Secure and that GATT operations still work.
//  4. Verify the on-disk persistence file is non-empty (gob-encoded
//     LTK was saved).
func TestLoopbackGATTNotificationEncryptedWithKeyCache(t *testing.T) {
	_, a, b := NewWorld(silentLogger())

	// Use temp files for SMP key persistence — the loopback drives the
	// real on-disk gob persistence path.
	keysDir := t.TempDir()
	keysPathA := filepath.Join(keysDir, "smp-A.gob")
	keysPathB := filepath.Join(keysDir, "smp-B.gob")

	mkStackCfg := func(advertiserOn bool, keysPath string) *ble.BluetoothStackConfig {
		cfg := ble.DefaultConfig()
		cfg.BLEScannerUse = false
		cfg.BLEAdvertiserUse = advertiserOn
		cfg.BLEConnecterUse = true
		cfg.BLEAdvertiserConfig = bleadvertiser.DefaultConfig()
		cfg.BLEAdvertiserConfig.AlwaysAdvertising = false
		cfg.BLEScannerConfig = &blescanner.BLEScannerConfig{
			StoreGAPMap:         true,
			ScanCycleDurationMs: 100,
			ScanCycleActiveDuty: 1,
		}
		cfg.SMPConfig = blesmp.DefaultConfig()
		cfg.SMPConfig.StoredKeysPath = keysPath
		// Bonding ON, MITM OFF → LE SC Just Works.
		cfg.SMPConfig.DefaultConnConfig = &blesmp.SMPConnConfig{
			AuthReq:        0x01,
			StaticPasscode: -1,
			MinKeySize:     16,
		}
		cfg.HCIControllerConfig.WatchdogTimeout = 0
		cfg.HCIControllerConfig.AwaitStartup = false
		cfg.HCIControllerConfig.PrivacyAdvertise = false
		cfg.HCIControllerConfig.PrivacyConnect = false
		cfg.HCIControllerConfig.PrivacyScan = false
		cfg.HCIControllerConfig.LERandomAddrBits = 32
		return cfg
	}

	cfgA := mkStackCfg(false, keysPathA)
	cfgB := mkStackCfg(true, keysPathB)

	stackA := ble.New(silentLogger(), cfgA, a)
	stackB := ble.New(silentLogger(), cfgB, b)

	// Run both stacks.
	aReady := make(chan struct{})
	bReady := make(chan struct{})
	go func() { stackA.Run(func() { close(aReady) }) }()
	go func() { stackB.Run(func() { close(bReady) }) }()
	defer stackA.Close()
	defer stackB.Close()
	<-aReady
	<-bReady

	// B's GATT structure: one notify+read characteristic. Read does
	// not require encryption (so initial discovery succeeds before
	// pairing); the user-visible sensitivity here is "subscribe is
	// only meaningful after pairing", which we verify by exchanging
	// notifications over the encrypted channel.
	bStruct := attstructure.NewStructure()
	bSvc := bStruct.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	bChar := bSvc.AddCharacteristic(
		bleutil.UUIDFromStringPanic("2a3a"),
		attstructure.CharacteristicRead|attstructure.CharacteristicNotify,
		attstructure.ValueConfig{LengthMax: 32},
	)

	smpCfg := &blesmp.SMPConnConfig{
		AuthReq:        0x01, // bond, no MITM
		StaticPasscode: -1,
		MinKeySize:     16,
	}

	// ---- SESSION 1: pair, encrypt, exchange a notification ----
	t.Run("Session1_Pair", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		connA, connB := connectStacks(t, ctx, stackA, stackB)

		gattCfgA := &bleatt.GattDeviceConfig{MTU: 247, DeviceName: "A", DiscoverRemoteOnConnect: true}
		gattCfgB := &bleatt.GattDeviceConfig{MTU: 247, DeviceName: "B", DiscoverRemoteOnConnect: false}

		// Wire B first and wait for its SMP to be live before starting
		// A's L2CAP — A's leTryEncryptLTK / GoSecure can fire HCI traffic
		// the moment its SMP PSM callback runs, and B's AddConn is the
		// only thing that publishes B.SMPConn (with addresses) to the
		// connmgr where the LELongTermKeyRequest event handler reads it.
		wB := wireGATTWithSMP(stackB, connB, bStruct, gattCfgB, smpCfg, false)
		defer wB.l2.Close()
		wB.startL2()
		select {
		case <-wB.smpReady:
		case <-time.After(2 * time.Second):
			t.Fatal("B SMP not ready")
		}

		wA := wireGATTWithSMP(stackA, connA, attstructure.NewStructure(), gattCfgA, smpCfg, true)
		defer wA.l2.Close()
		wA.startL2()
		select {
		case <-wA.smpReady:
		case <-time.After(2 * time.Second):
			t.Fatal("A SMP not ready")
		}

		// Central initiates pairing.
		state, err := wA.smpConn.GoSecure(ctx, true)
		if err != nil {
			t.Fatalf("GoSecure: %v", err)
		}
		if state != blesmp.StateSecure {
			t.Fatalf("GoSecure returned state %d, want StateSecure", state)
		}

		// Both sides should now report encrypted+bonded.
		enc, _, bonded := wA.smpConn.GetSecurity()
		if !enc || !bonded {
			t.Errorf("A GetSecurity: enc=%v bonded=%v want true/true", enc, bonded)
		}

		// GATT discovery + subscribe + notification round trip.
		remote := wA.gatt.ClientGetStructure(ctx)
		if remote == nil {
			t.Fatal("A did not discover B's GATT structure")
		}
		rChar := remote.GetService(bleutil.UUIDFromStringPanic("180a")).GetCharacteristic(bleutil.UUIDFromStringPanic("2a3a"))
		if rChar == nil {
			t.Fatal("A did not discover the 2a3a characteristic")
		}
		got := make(chan []byte, 4)
		if err := rChar.Subscribe(ctx, func(value []byte) {
			select {
			case got <- append([]byte(nil), value...):
			default:
			}
		}); err != nil {
			t.Fatalf("Subscribe: %v", err)
		}
		want := []byte("encrypted-hello")
		if _, err := bChar.SetValue(ctx, want); err != nil {
			t.Fatalf("SetValue: %v", err)
		}
		select {
		case rx := <-got:
			if string(rx) != string(want) {
				t.Errorf("session1 notify: got %q want %q", rx, want)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("session1: no notification")
		}

		// Tear down the connection. B's connecter notices the disconnect
		// and exits the peripheral Connect.
		_ = connA.Close()
	})

	// Verify the LTK was persisted to disk.
	if st, err := os.Stat(keysPathA); err != nil || st.Size() == 0 {
		t.Errorf("A keys file: %v size=%d (expected non-empty)", err, st.Size())
	}
	if st, err := os.Stat(keysPathB); err != nil || st.Size() == 0 {
		t.Errorf("B keys file: %v size=%d (expected non-empty)", err, st.Size())
	}

	// Brief settle so the previous connection is fully torn down on
	// both sides before we reconnect.
	time.Sleep(200 * time.Millisecond)

	// ---- SESSION 2: reconnect, expect cached-LTK auto-encryption ----
	t.Run("Session2_CachedLTK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		connA, connB := connectStacks(t, ctx, stackA, stackB)

		gattCfgA := &bleatt.GattDeviceConfig{MTU: 247, DeviceName: "A", DiscoverRemoteOnConnect: true}
		gattCfgB := &bleatt.GattDeviceConfig{MTU: 247, DeviceName: "B", DiscoverRemoteOnConnect: false}

		// Wire B first (and wait for its SMP) so that A's auto-encrypt
		// path doesn't fire LEEnableEncryption before B's SMPConn is
		// fully published to the connmgr.
		wB := wireGATTWithSMP(stackB, connB, bStruct, gattCfgB, smpCfg, false)
		defer wB.l2.Close()
		wB.startL2()
		select {
		case <-wB.smpReady:
		case <-time.After(2 * time.Second):
			t.Fatal("B SMP not ready in session 2")
		}

		wA := wireGATTWithSMP(stackA, connA, attstructure.NewStructure(), gattCfgA, smpCfg, true)
		defer wA.l2.Close()
		wA.startL2()
		select {
		case <-wA.smpReady:
		case <-time.After(2 * time.Second):
			t.Fatal("A SMP not ready in session 2")
		}

		// Encryption should engage via the cached LTK as part of
		// SMP.AddConn → leTryEncryptLTK. GoSecure(false) returns when
		// state == Secure without initiating pairing.
		state, err := wA.smpConn.GoSecure(ctx, false)
		if err != nil {
			t.Fatalf("session 2 GoSecure(false): %v", err)
		}
		if state != blesmp.StateSecure {
			t.Fatalf("session 2 expected StateSecure (cached LTK), got %d", state)
		}
		enc, _, bonded := wA.smpConn.GetSecurity()
		if !enc {
			t.Errorf("session 2 not encrypted")
		}
		if !bonded {
			t.Errorf("session 2 not bonded")
		}

		// GATT round trip works over the auto-encrypted link.
		remote := wA.gatt.ClientGetStructure(ctx)
		if remote == nil {
			t.Fatal("session 2: GATT discovery failed")
		}
		rChar := remote.GetService(bleutil.UUIDFromStringPanic("180a")).GetCharacteristic(bleutil.UUIDFromStringPanic("2a3a"))
		got := make(chan []byte, 4)
		if err := rChar.Subscribe(ctx, func(value []byte) {
			select {
			case got <- append([]byte(nil), value...):
			default:
			}
		}); err != nil {
			t.Fatalf("session 2 Subscribe: %v", err)
		}
		want := []byte("session2-cached")
		if _, err := bChar.SetValue(ctx, want); err != nil {
			t.Fatalf("session 2 SetValue: %v", err)
		}
		select {
		case rx := <-got:
			if string(rx) != string(want) {
				t.Errorf("session 2 notify: got %q want %q", rx, want)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("session 2: no notification")
		}

		_ = connA.Close()
	})
}

// Ensure the test compile imports stay tidy.
var _ = (*ble.BluetoothStack)(nil)
