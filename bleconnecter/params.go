package bleconnecter

import (
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"

	"github.com/sirupsen/logrus"
)

type BLEConnectionParametersActual struct {
	Interval            uint16
	Latency             uint16
	Timeout             uint16
	MasterClockAccuracy uint8
}

type BLEConnectionParametersRequested struct {
	ConnectionIntervalMin uint16
	ConnectionIntervalMax uint16
	ConnectionLatency     uint16
	SupervisionTimeout    uint16
	MinCELength           uint16
	MaxCELength           uint16
}

func (c *BLEConnecter) leConnectionParameterRequestHandler(event *hcievents.LERemoteConnectionParameterRequestEvent) *hcievents.LERemoteConnectionParameterRequestEvent {
	c.logger.WithField("0handle", event.ConnectionHandle).Debug("Accepting connection parameters update")

	go c.ctrl.Cmds.LERemoteConnectionParameterRequestReplySync(
		hcicommands.LERemoteConnectionParameterRequestReplyInput{
			ConnectionHandle: event.ConnectionHandle,
			IntervalMin:      event.IntervalMin,
			IntervalMax:      event.IntervalMax,
			Latency:          event.Latency,
			Timeout:          event.Timeout,
		}, nil)
	return event
}

func (c *BLEConnecter) leConnectionUpdateCompleteHandler(event *hcievents.LEConnectionUpdateCompleteEvent) *hcievents.LEConnectionUpdateCompleteEvent {
	hwConn := c.ctrl.ConnMgr.FindConnectionByHandle(event.ConnectionHandle)
	if hwConn != nil {
		switch conn := hwConn.AppConn.(type) {
		case *BLEConnection:
			conn.parametersMutex.Lock()
			conn.parametersActual = BLEConnectionParametersActual{
				Interval: event.ConnectionInterval,
				Latency:  event.ConnectionLatency,
				Timeout:  event.SupervisionTimeout,
			}

			c.logger.WithFields(logrus.Fields{
				"0handle": event.ConnectionHandle,
				"1params": conn.parametersActual,
			}).Debug("Connection parameters changed")

			conn.parametersMutex.Unlock()
		}
	}

	return event
}

func (r *BLEConnectionParametersRequested) makeValid() {
	r.ConnectionIntervalMin = bleutil.ClampUint16(r.ConnectionIntervalMin, 0x6, 0xC80)
	r.ConnectionIntervalMax = bleutil.ClampUint16(r.ConnectionIntervalMax, 0x6, 0xC80)
	if r.ConnectionIntervalMin > r.ConnectionIntervalMax {
		r.ConnectionIntervalMin = r.ConnectionIntervalMax
	}
	r.ConnectionLatency = bleutil.ClampUint16(r.ConnectionLatency, 0, 0x1F3)

	minSupervisionTimeout := uint16(((1+uint64(r.ConnectionLatency))*(1250*uint64(r.ConnectionIntervalMax))*2)/10000) + 1
	if minSupervisionTimeout < 0xA {
		minSupervisionTimeout = 0xA
	}
	r.SupervisionTimeout = bleutil.ClampUint16(r.SupervisionTimeout, minSupervisionTimeout, 0xC80)
}

func (c *BLEConnection) UpdateParams(request BLEConnectionParametersRequested) error {
	request.makeValid()

	return c.connecter.ctrl.Cmds.LEConnectionUpdateSync(hcicommands.LEConnectionUpdateInput{
		ConnectionHandle:      c.event.ConnectionHandle,
		ConnectionIntervalMin: request.ConnectionIntervalMin,
		ConnectionIntervalMax: request.ConnectionIntervalMax,
		ConnectionLatency:     request.ConnectionLatency,
		SupervisionTimeout:    request.SupervisionTimeout,
		MinCELength:           request.MinCELength,
		MaxCELength:           request.MaxCELength,
	})
}
