// Copyright 2015, Truveris Inc. All Rights Reserved.

package turret

import (
	"errors"
	"fmt"
	"time"

	"github.com/truveris/gousb/usb"
)

// Device types
const (
	DeviceTypeThunder = iota
	DeviceTypeClassic = iota
)

var (
	errNoLight = errors.New("device has no light")
)

// Turret is a wrapper around a USB Device.
type Turret struct {
	*usb.Device
	Type   int
	Closed bool
}

// NewTurret creates a Turret from a USB device.
func NewTurret(dev *usb.Device) *Turret {
	t := &Turret{
		Device: dev,
		Closed: false,
	}
	return t
}

// Command immediately passes the given command to a turret device. If you plan
// on using this function concurrently across multiple go-routines, you may
// want to use the QueueCommand and ConsumeCommands facility instead.
func (t *Turret) Command(cmdtype, cmd byte) error {
	if t.Type == DeviceTypeThunder {
		_, err := t.Device.Control(0x21, 0x09, 0, 0, []byte{cmdtype, cmd, 0, 0, 0, 0, 0, 0})
		if err != nil {
			return err
		}
	} else if t.Type == DeviceTypeClassic {
		if cmdtype != CmdTypeTurret {
			return errNoLight
		}

		_, err := t.Device.Control(0x21, 0x09, 0x0200, 0, []byte{cmd})
		if err != nil {
			return err
		}
	}

	return nil
}

// NormalizeDuration adjust the provided duration to be within reasonable limits.
// For example, most turrets take 8 seconds to do a full horizontal cycle.
// Because of that, we can safely limit the maximum horizontal movement
// duration to 8 seconds.
func (t *Turret) NormalizeDuration(cmd byte, duration time.Duration) time.Duration {
	switch cmd {
	case CmdTurretLeft, CmdTurretRight:
		if duration > 8*time.Second {
			duration = 6800 * time.Millisecond
		}
	case CmdTurretUp, CmdTurretDown:
		if duration > 2*time.Second {
			duration = 1500 * time.Millisecond
		}
	}
	return duration
}

// TimedCommand executes the given command and wait a little.
func (t *Turret) TimedCommand(cmdType, cmd byte, duration time.Duration) {
	t.Command(cmdType, cmd)
	duration = t.NormalizeDuration(cmd, duration)
	if duration > 0 {
		time.Sleep(duration)
	}
}

// Close ends the connection to the USB device, closes the Input channel. You
// should never call Close() directly if you queue commands and uses
// ConsumeCommands.  Use it if you plan on feeding the turret directly with
// commands.
func (t *Turret) Close() {
	if t.Closed {
		return
	}
	t.Device.Close()
	t.Closed = true
}

// Light is used to turn the turret light on or off.
func (t *Turret) Light(on bool) {
	var cmd byte

	if on {
		cmd = CmdLightOn
	} else {
		cmd = CmdLightOff
	}

	t.TimedCommand(CmdTypeLight, cmd, 0)
}

// BlinkOn will turn the light of the turret on and off a few times, ending
// with a lit light.
func (t *Turret) BlinkOn(times int) {
	for i := 0; i < times; i++ {
		t.TimedCommand(CmdTypeLight, CmdLightOff, 200*time.Millisecond)
		t.TimedCommand(CmdTypeLight, CmdLightOn, 200*time.Millisecond)
	}
}

// BlinkOff will turn the light of the turret on and off a few times, ending
// with a turned off light.
func (t *Turret) BlinkOff(times int) {
	for i := 0; i < times; i++ {
		t.TimedCommand(CmdTypeLight, CmdLightOn, 200*time.Millisecond)
		t.TimedCommand(CmdTypeLight, CmdLightOff, 200*time.Millisecond)
	}
}

// Left rotates the turret left for the specified duration.
func (t *Turret) Left(duration time.Duration) {
	t.TimedCommand(CmdTypeTurret, CmdTurretLeft, duration)
}

// Right rotates the turret right for the specified duration.
func (t *Turret) Right(duration time.Duration) {
	t.TimedCommand(CmdTypeTurret, CmdTurretRight, duration)
}

// Up tilts the turret up for the specified duration.
func (t *Turret) Up(duration time.Duration) {
	t.TimedCommand(CmdTypeTurret, CmdTurretUp, duration)
}

// Down tilts the turret down for the specified duration.
func (t *Turret) Down(duration time.Duration) {
	t.TimedCommand(CmdTypeTurret, CmdTurretDown, duration)
}

// Stop interrupts the turret movements.
func (t *Turret) Stop() {
	t.TimedCommand(CmdTypeTurret, CmdTurretStop, 50*time.Millisecond)
}

// Reset rotates and tilts the turret in the given direction until it gets
// parked in the leftmost and lowest position.
func (t *Turret) Reset() {
	t.Down(2 * time.Second)
	t.Left(8 * time.Second)
}

// Fire one or multiple shots.
func (t *Turret) Fire(shots int) {
	for i := 0; i < shots; i++ {
		t.TimedCommand(CmdTypeTurret, CmdTurretFire, 4500*time.Millisecond)
	}
}

// HumanReadableType returns the model of the turret in a human readable format.
func (t *Turret) HumanReadableType() string {
	switch t.Type {
	case DeviceTypeThunder:
		return "Dream Cheeky Thunder"
	case DeviceTypeClassic:
		return "Classic"
	default:
		name := fmt.Sprintf("Unknown Turret (0x%H:0x%H)", t.Device.Descriptor.Vendor, t.Device.Descriptor.Product)
		return name
	}
}
