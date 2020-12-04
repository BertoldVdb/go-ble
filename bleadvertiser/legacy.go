package bleadvertiser

import (
	"bytes"
	"errors"
	"time"

	"github.com/BertoldVdb/go-ble/hci"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

var (
	ErrorClosed  = errors.New("Advertiser is closed")
	ErrorExpired = errors.New("Dataset is not current, cannot apply unless forced")
)

type LegacyAdvertisingData struct {
	Active       bool
	BeaconPacket []byte
	ScanPacket   []byte
	PeerAddr     bleutil.BLEAddr
	Type         LegacyAdvertisementType
	IntervalMin  uint16
	IntervalMax  uint16
	UseAllowlist bool
	AddrType     bleutil.MacAddrType

	version uint64
}

func (a *BLEAdvertiser) legacyAdvertisingConfigure(data *LegacyAdvertisingData) error {
	var advParams hcicommands.LESetAdvertisingParametersInput
	if data != nil {
		filterPolicy := uint8(0)
		if data.UseAllowlist {
			filterPolicy = 2
		}

		if data.IntervalMin > data.IntervalMax {
			data.IntervalMin = data.IntervalMax
		}

		if data.IntervalMin == 0 {
			data.IntervalMin = 0x20
			data.IntervalMax = 0x40
		}

		advParams = hcicommands.LESetAdvertisingParametersInput{
			AdvertisingIntervalMin:  data.IntervalMin,
			AdvertisingIntervalMax:  data.IntervalMax,
			AdvertisingType:         uint8(data.Type),
			PeerAddress:             data.PeerAddr.MacAddr,
			PeerAddressType:         data.PeerAddr.MacAddrType,
			OwnAddressType:          data.AddrType,
			AdvertisingChannelMap:   7,
			AdvertisingFilterPolicy: filterPolicy,
		}
	} else {
		advParams = hcicommands.LESetAdvertisingParametersInput{
			AdvertisingIntervalMin: 0x30,
			AdvertisingIntervalMax: 0x60,
			AdvertisingChannelMap:  7,
		}
	}

	/* There is no way to test if advertising is enabled (and some hardware disables it on its own),
	   but this command cannot be run when it is, so we can know */
	err := a.ctrl.Cmds.LESetAdvertisingParametersSync(advParams)

	if err != nil {
		err := a.ctrl.Cmds.LESetAdvertisingEnableSync(hcicommands.LESetAdvertisingEnableInput{
			AdvertisingEnable: 0,
		})
		if err != nil || data == nil {
			return err
		}

		err = a.ctrl.Cmds.LESetAdvertisingParametersSync(advParams)
		if err != nil {
			return err
		}
	}

	if data == nil {
		return nil
	}

	var buffer [31]byte

	if !bytes.Equal(data.BeaconPacket, a.legacyAdvertisingData) {
		bleutil.Assert(len(data.BeaconPacket) <= 31, "Base advertising data is too long")

		a.legacyAdvertisingData = a.legacyAdvertisingData[:len(data.BeaconPacket)]
		copy(a.legacyAdvertisingData, data.BeaconPacket)
		copy(buffer[:], data.BeaconPacket)

		err = a.ctrl.Cmds.LESetAdvertisingDataSync(hcicommands.LESetAdvertisingDataInput{
			AdvertisingDataLength: uint8(len(data.BeaconPacket)),
			AdvertisingData:       buffer,
		})
		if err != nil {
			return err
		}
	}

	if !bytes.Equal(data.ScanPacket, a.legacyAdvertisingScanData) {
		bleutil.Assert(len(data.ScanPacket) <= 31, "Scan response advertising data is too long")

		a.legacyAdvertisingScanData = a.legacyAdvertisingScanData[:len(data.ScanPacket)]
		copy(a.legacyAdvertisingScanData, data.ScanPacket)
		copy(buffer[:], data.ScanPacket)

		err = a.ctrl.Cmds.LESetScanResponseDataSync(hcicommands.LESetScanResponseDataInput{
			ScanResponseDataLength: uint8(len(data.ScanPacket)),
			ScanResponseData:       buffer,
		})
		if err != nil {
			return err
		}
	}

	return a.ctrl.Cmds.LESetAdvertisingEnableSync(hcicommands.LESetAdvertisingEnableInput{
		AdvertisingEnable: 1,
	})
}

func (a *BLEAdvertiser) LegacyAdvertisingSetConnection(useAllowlist bool, peerAddr *bleutil.BLEAddr) (func() error, error) {
	a.legacyAdvertisingInit()

	data, _ := a.legacyAdvertisingBaseSlot.GetData()

	a.logger.WithFields(logrus.Fields{
		"1useAllowlist": useAllowlist,
		"2peerAddr":     peerAddr,
	}).Info("Setting beacon mode")

	data.Active = true
	data.Type = LegacyAdvertisementTypeInd
	data.UseAllowlist = useAllowlist

	if peerAddr != nil {
		data.PeerAddr = *peerAddr
		data.Type = LegacyAdvertisementTypeDirectIndLowDuty
	}

	newData, err := a.legacyAdvertisingBaseSlot.ReplaceData(true, data)
	if err != nil {
		return nil, err
	}

	cancelFunc := func() error {
		newData.Active = a.legacyAdvertisingAlwaysOn
		newData.Type = LegacyAdvertisementTypeScanInd
		_, _ = a.legacyAdvertisingBaseSlot.ReplaceData(false, data)
		return nil
	}

	return cancelFunc, nil
}

type LegacyAdvertisingSlot struct {
	parent *BLEAdvertiser
	valid  bool
	index  int
	data   LegacyAdvertisingData
}

func (s *LegacyAdvertisingSlot) Close() {
	s.parent.legacyAdvertisingMutex.Lock()
	defer s.parent.legacyAdvertisingMutex.Unlock()

	s.valid = false

	select {
	case s.parent.legacyAdvertisingUpdateChan <- -1:
	default:
	}
}

func (s *LegacyAdvertisingSlot) GetData() (LegacyAdvertisingData, error) {
	s.parent.legacyAdvertisingMutex.Lock()
	defer s.parent.legacyAdvertisingMutex.Unlock()

	return s.data, nil
}

func (s *LegacyAdvertisingSlot) ReplaceData(force bool, new LegacyAdvertisingData) (LegacyAdvertisingData, error) {
	s.parent.legacyAdvertisingMutex.Lock()
	defer s.parent.legacyAdvertisingMutex.Unlock()

	old := s.data

	if !force && new.version > 0 && new.version < old.version {
		return new, ErrorExpired
	}

	new.version = old.version + 1
	s.data = new

	select {
	case s.parent.legacyAdvertisingUpdateChan <- s.index:
	default:
	}

	return new, nil
}

func (a *BLEAdvertiser) LegacyAdvertisingGetSlot() *LegacyAdvertisingSlot {
	a.legacyAdvertisingMutex.Lock()
	defer a.legacyAdvertisingMutex.Unlock()

	new := LegacyAdvertisingSlot{
		parent: a,
		valid:  true,
	}

	for i, m := range a.legacyAdvertisingSlots {
		if !m.valid {
			new.index = i
			a.legacyAdvertisingSlots[i] = new
			return &a.legacyAdvertisingSlots[new.index]
		}
	}

	new.index = len(a.legacyAdvertisingSlots)
	a.legacyAdvertisingSlots = append(a.legacyAdvertisingSlots, new)

	return &a.legacyAdvertisingSlots[new.index]
}

func (a *BLEAdvertiser) legacyAdvertisingInit() {
	a.legacyAdvertisingInitOnce.Do(func() {
		a.legacyAdvertisingBaseSlot = a.LegacyAdvertisingGetSlot()
		a.legacyAdvertisingAlwaysOn = a.config.AlwaysAdvertising

		a.legacyAdvertisingData = make([]byte, 0, 31)
		a.legacyAdvertisingScanData = make([]byte, 0, 31)

		beaconData := LegacyAdvertisingData{
			Active:   a.legacyAdvertisingAlwaysOn,
			Type:     LegacyAdvertisementTypeScanInd,
			AddrType: a.ctrl.GetLERecommenedOwnAddrType(hci.LEAddrUsageAdvertise),
		}

		/* TODO: Basic beacon, can be improved with bin packer utility function later */
		beaconData.BeaconPacket = UtilPDUAddRecord(beaconData.BeaconPacket, 1, []byte{a.config.DeviceFlags})

		if bleutil.UUIDBase != a.config.DeviceService {
			service := a.config.DeviceService.UUIDToBytes()
			serviceType := uint8(3)
			if len(service) > 2 {
				serviceType = 7
			}
			beaconData.BeaconPacket = UtilPDUAddRecord(beaconData.BeaconPacket, serviceType, service)
		}

		nameType := uint8(9)
		nameDev := a.config.DeviceName
		if len(a.config.DeviceName) > 27 {
			nameType = 8
			nameDev = nameDev[0:27]
		}
		beaconData.ScanPacket = UtilPDUAddRecord(beaconData.ScanPacket, nameType, []byte(nameDev))

		a.legacyAdvertisingBaseSlot.ReplaceData(true, beaconData)
	})
}

func (a *BLEAdvertiser) legacyAdvertisingManager() error {
	a.legacyAdvertisingInit()

	advIndex := 0
	var advTick *time.Ticker

	setupNextAdv := func() error {
		a.legacyAdvertisingMutex.Lock()
		defer a.legacyAdvertisingMutex.Unlock()

		numElements := 0

		for _, m := range a.legacyAdvertisingSlots {
			if m.valid && m.data.Active {
				numElements++
			}
		}

		if numElements > 1 {
			if advTick == nil {
				advTick = time.NewTicker(100 * time.Millisecond)
			}
		} else {
			if advTick != nil {
				advTick.Stop()
				advTick = nil
			}
		}

		for range a.legacyAdvertisingSlots {
			index := advIndex
			advIndex++
			if advIndex >= len(a.legacyAdvertisingSlots) {
				advIndex = 0
			}

			slot := a.legacyAdvertisingSlots[index]

			if slot.valid && slot.data.Active {
				return a.legacyAdvertisingConfigure(&a.legacyAdvertisingSlots[index].data)
			}
		}

		return a.legacyAdvertisingConfigure(nil)
	}

loop:
	for {
		var tickChan <-chan (time.Time)
		if advTick != nil {
			tickChan = advTick.C
		}

		select {
		case <-a.closeflag.Chan():
			break loop
		case <-tickChan:
			setupNextAdv()
		case index, ok := <-a.legacyAdvertisingUpdateChan:
			if !ok {
				break loop
			}
			if index >= 0 {
				advIndex = index
			}

			setupNextAdv()
		}
	}

	if advTick != nil {
		advTick.Stop()
	}

	return nil
}
