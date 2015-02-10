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
	Input  chan Command
	Done   chan bool
	Closed bool
}

// NewTurret creates a Turret from a USB device.
func NewTurret(dev *usb.Device) *Turret {
	t := &Turret{
		Device: dev,
		Closed: false,
		Input:  make(chan Command, 64),
		Done:   make(chan bool, 0),
	}
	go t.ConsumeCommands()
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

// QueueCommand schedules a command in the Input channel of this Turret.  This
// is the prefered way to schedule commands.
func (t *Turret) QueueCommand(cmdtype, cmd byte, duration time.Duration) {
	t.Input <- Command{Type: cmdtype, Value: cmd, Duration: duration}
}

// NormalizeDuration adjust the provided duration to be within reasonable limits.
// For example, most turrets take 8 seconds to do a full horizontal cycle.
// Because of that, we can safely limit the maximum horizontal movement
// duration to 8 seconds.
func (t *Turret) NormalizeDuration(cmd byte, duration time.Duration) time.Duration {
	switch cmd {
	case CmdTurretLeft, CmdTurretRight:
		if duration > 8*time.Second {
			duration = 8 * time.Second
		}
	case CmdTurretUp, CmdTurretDown:
		if duration > 2*time.Second {
			duration = 2 * time.Second
		}
	}
	return duration
}

// ConsumeCommands is a go routine started by NewTurret which executes the
// Turret commands sequentially from the Input channel.
func (t *Turret) ConsumeCommands() {
	for cmd := range t.Input {
		if cmd.Type == CmdTypeDone {
			t.Done <- true
			return
		}
		t.Command(cmd.Type, cmd.Value)
		duration := t.NormalizeDuration(cmd.Value, cmd.Duration)
		if duration > 0 {
			time.Sleep(duration)
		}
	}
}

// Shutdown sends a Command to identify the end of the queue, then waits for
// the ConsumeCommands function to terminate via the Turret Done channel.
func (t *Turret) Shutdown() {
	if t.Closed {
		return
	}
	t.QueueCommand(CmdTypeDone, 0, 0)
	<-t.Done
	t.Close()
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
	close(t.Input)
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

	t.QueueCommand(CmdTypeLight, cmd, 0)
}

// BlinkOn will turn the light of the turret on and off a few times, ending
// with a lit light.
func (t *Turret) BlinkOn(times int) {
	for i := 0; i < times; i++ {
		t.QueueCommand(CmdTypeLight, CmdLightOff, 200*time.Millisecond)
		t.QueueCommand(CmdTypeLight, CmdLightOn, 200*time.Millisecond)
	}
}

// BlinkOff will turn the light of the turret on and off a few times, ending
// with a turned off light.
func (t *Turret) BlinkOff(times int) {
	for i := 0; i < times; i++ {
		t.QueueCommand(CmdTypeLight, CmdLightOn, 200*time.Millisecond)
		t.QueueCommand(CmdTypeLight, CmdLightOff, 200*time.Millisecond)
	}
}

// Left rotates the turret left for the specified duration.
func (t *Turret) Left(duration time.Duration) {
	t.QueueCommand(CmdTypeTurret, CmdTurretLeft, duration)
}

// Right rotates the turret right for the specified duration.
func (t *Turret) Right(duration time.Duration) {
	t.QueueCommand(CmdTypeTurret, CmdTurretRight, duration)
}

// Up tilts the turret up for the specified duration.
func (t *Turret) Up(duration time.Duration) {
	t.QueueCommand(CmdTypeTurret, CmdTurretUp, duration)
}

// Down tilts the turret down for the specified duration.
func (t *Turret) Down(duration time.Duration) {
	t.QueueCommand(CmdTypeTurret, CmdTurretDown, duration)
}

// Stop interrupts the turret movements.
func (t *Turret) Stop() {
	t.QueueCommand(CmdTypeTurret, CmdTurretStop, 50*time.Millisecond)
}

// Reset rotates and tilts the turret in the given direction until it gets
// parked in the leftmost and lowest position.
func (t *Turret) Reset() {
	t.Left(8 * time.Second)
	t.Down(2 * time.Second)
}

// Fire one or multiple shots.
func (t *Turret) Fire(shots int) {
	for i := 0; i < shots; i++ {
		t.QueueCommand(CmdTypeTurret, CmdTurretFire, 4500*time.Millisecond)
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
