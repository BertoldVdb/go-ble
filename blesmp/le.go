package blesmp

import (
	"context"
	"time"

	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
)

func (c *SMPConn) leEncryptWait(work func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

drain:
	for {
		select {
		case <-c.encUpdateChan:
		default:
			break drain
		}
	}

	err := work()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case value := <-c.encUpdateChan:
			if value {
				return nil
			}
		}
	}
}

func (s *SMPConn) leSetKeyFlagsFromLTK(ltk smpStoredLTK) {
	s.keyMutex.Lock()
	s.keyIsBonded = ltk.Bonded
	s.keyIsAuthenticated = ltk.Authenticated
	s.keyMutex.Unlock()
}

func (c *SMPConn) leEncrypt(ltk smpStoredLTK) error {
	raw := c.rawConnLE()

	return c.leEncryptWait(func() error {
		err := raw.Encrypt(ltk.EDIV, ltk.Rand, ltk.LTK)
		if err == nil {
			c.leSetKeyFlagsFromLTK(ltk)
		}
		return err
	})
}

func (s *SMP) connmgrLEGetKey(conn *hciconnmgr.Connection, event *hcievents.LELongTermKeyRequestEvent) ([]byte, *hcievents.LELongTermKeyRequestEvent) {
	smpConn := smpConnFromConn(conn)
	if smpConn == nil {
		// Connection has no SMP yet — answer with NegativeReply by
		// returning a non-16-byte key. The connmgr's LTK request
		// handler interprets that as "no key available".
		return nil, event
	}

	s.storedKeysPersist.Lock()
	ltk, ok := s.storedKeys[makeSMPStoredLTKMapKey(false, smpConn.addrLELocal, smpConn.addrLERemote, event.EncryptedDiversifier, event.RandomNumber)]
	s.storedKeysPersist.Unlock()
	if !ok {
		return nil, event
	}

	smpConn.leSetKeyFlagsFromLTK(ltk)

	return ltk.LTK[:], event
}

func (c *SMPConn) leTryEncryptLTK() error {
	c.parent.storedKeysPersist.Lock()
	ltk, ok := c.parent.storedKeys[makeSMPStoredLTKMapKey(true, c.addrLELocal, c.addrLERemote, 0, 0)]
	c.parent.storedKeysPersist.Unlock()
	if !ok {
		return nil
	}

	err := c.leEncrypt(ltk)
	if err != nil {
		return err
	}

	/* leEncryptWait drained the EncryptionChange signal that smpHandler
	   would otherwise consume to transition state to Secure, so do it
	   here. Without this the cached-LTK fast path leaves SMP in
	   StateInsecure even though the link is fully encrypted, which makes
	   GetSecurity() report unencrypted and breaks GoSecure(ctx, false). */
	c.setState(StateSecure)

	return nil
}
