package blescannerjson

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/BertoldVdb/go-ble/blescanner"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

type JSONGapEntry struct {
	GAPType   uint8
	EventType string
	Payload   string
}

type JSONScanDevice struct {
	Address     string
	Name        string
	Flags       uint8
	Connectable int
	RSSI        int8
	LastSeenMs  int64

	Services []string
	GAP      []JSONGapEntry
}

type JSONScanResults struct {
	Now      int64
	ScanType string
	Devices  []JSONScanDevice
}

type ScanJSONGenerator struct {
	sync.Mutex

	/* These elements are kept to reduce allocations a bit */
	results    JSONScanResults
	knownAddrs []bleutil.BLEAddr
	services   []bleutil.UUID
	gapTypes   []int
	gapRecord  *blescanner.GAPRecord

	lastUpdate time.Time
	lastResult []byte

	scanner *blescanner.BLEScanner
}

func New(scanner *blescanner.BLEScanner) *ScanJSONGenerator {
	return &ScanJSONGenerator{
		scanner: scanner,
	}
}

func (jg *ScanJSONGenerator) generateJSONLocked() ([]byte, error) {
	now := time.Now()

	jg.knownAddrs = jg.scanner.KnownDevicesAddresses(jg.knownAddrs)
	sort.SliceStable(jg.knownAddrs, func(i, j int) bool {
		return jg.knownAddrs[i].IsLess(jg.knownAddrs[j])
	})

	jg.results.Devices = jg.results.Devices[:0]
	jg.results.Now = now.UnixNano() / 1e6
	switch jg.scanner.GetScanType() {
	case -1:
		jg.results.ScanType = "Off"
	case 0:
		jg.results.ScanType = "Passive"
	case 1:
		jg.results.ScanType = "Active"
	}

	for _, addr := range jg.knownAddrs {
		dev := jg.scanner.GetDevice(addr)
		if dev == nil {
			continue
		}

		device := JSONScanDevice{
			Address:    addr.String(),
			Name:       dev.GetName(),
			Flags:      dev.GetFlags(),
			RSSI:       dev.GetRSSI(),
			LastSeenMs: now.Sub(dev.LastSeen()).Milliseconds(),
		}

		device.Connectable = 0
		if dev.IsConnectable() {
			device.Connectable = 1
		}

		jg.services = dev.GetServices(-1, jg.services)
		for _, m := range jg.services {
			device.Services = append(device.Services, m.String())
		}

		jg.gapTypes = dev.GetGAPTypes(jg.gapTypes)
		sort.IntSlice(jg.gapTypes).Sort()
		for _, i := range jg.gapTypes {
			jg.gapRecord = dev.GetGAPRecord(i, jg.gapRecord)
			if jg.gapRecord != nil {
				device.GAP = append(device.GAP, JSONGapEntry{
					GAPType:   jg.gapRecord.Type,
					EventType: jg.gapRecord.EventType.String(),
					Payload:   hex.EncodeToString(jg.gapRecord.Data),
				})
			}
		}

		jg.results.Devices = append(jg.results.Devices, device)

		dev.Release()
	}

	jsb, err := json.MarshalIndent(jg.results, "", "  ")
	if err == nil {
		jg.lastResult = jsb
		jg.lastUpdate = now
	}
	return jg.lastResult, err
}

func (jg *ScanJSONGenerator) GenerateJSONThrottled() ([]byte, error) {
	jg.Lock()
	defer jg.Unlock()

	var err error
	if time.Now().After(jg.lastUpdate.Add(400 * time.Millisecond)) {
		_, err = jg.generateJSONLocked()
	}

	return jg.lastResult, err
}

func (jg *ScanJSONGenerator) GenerateJSON() ([]byte, error) {
	jg.Lock()
	defer jg.Unlock()

	return jg.generateJSONLocked()
}

func (jg *ScanJSONGenerator) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	result, err := jg.GenerateJSONThrottled()

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
