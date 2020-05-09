package hcievents

import (
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
)

/* Remap some events, rest of the remapping is done in the spec parser (Why can't they put them sequentially?) */
func remapEvent(id int) int {
	if id == 109 {
		return 112
	}
	if id == 110 {
		return 113
	}
	if id == 111 {
		return 109
	}
	if id == 112 {
		return 110
	}
	if id == 113 {
		return 111
	}
	return id
}

func updateEventMask(current uint64, bit int, enabled bool) (uint64, bool) {
	mask := uint64(1) << bit

	var out uint64
	if enabled {
		out = current | mask
	} else {
		out = current & ^mask
	}

	return out, out != current
}

func (e *EventHandler) eventChangedMask(id int, enabled bool) error {
	eventMask, changed := updateEventMask(e.eventMask, id, enabled)
	if changed {
		err := e.cmds.BasebandSetEventMaskSync(hcicommands.BasebandSetEventMaskInput{
			EventMask: eventMask,
		})

		/* Normally the controller should ignore unknown bits, but some seem to
		   return a parameter invalid error. If we keep the bit set we will never
		   be able to change any events.
		   Some also return that error, but apply the new mask regardless... */
		if err == nil {
			e.eventMask = eventMask
		}

		return err
	}
	return nil
}

func (e *EventHandler) eventChangedMask2(id int, enabled bool) error {
	eventMask2, changed := updateEventMask(e.eventMask2, id-100, enabled)
	if changed {
		err := e.cmds.BasebandSetEventMaskPage2Sync(hcicommands.BasebandSetEventMaskPage2Input{
			EventMaskPage2: eventMask2,
		})

		if err == nil {
			e.eventMask2 = eventMask2
		}

		return err
	}
	return nil
}

func (e *EventHandler) eventChangedMaskLe(id int, enabled bool) error {
	eventMaskLe, changed := updateEventMask(e.eventMaskLe, id-200, enabled)
	if changed {
		/* LE events have a global mask that needs to be updated */
		err := e.eventChangedMask(61, eventMaskLe != 0)
		if err != nil {
			return err
		}

		err = e.cmds.LESetEventMaskSync(hcicommands.LESetEventMaskInput{
			LEEventMask: eventMaskLe,
		})

		if err == nil {
			e.eventMaskLe = eventMaskLe
		}

		return err
	}
	return nil
}

func (e *EventHandler) eventChanged(id int, enabled bool) error {
	/* If we can send commands we will enable/disable events */
	if e.cmds != nil {
		if id == 0xD || id == 0xE {
			/* These are the status/complete events and they should not be touched */
			return nil
		}

		id = remapEvent(id)

		e.enableMutex.Lock()
		defer e.enableMutex.Unlock()

		if id >= 200 {
			return e.eventChangedMaskLe(id, enabled)
		} else if id >= 100 {
			return e.eventChangedMask2(id, enabled)
		} else {
			return e.eventChangedMask(id, enabled)
		}
	}
	return nil
}
