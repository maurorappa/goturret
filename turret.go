// Copyright 2015, Truveris Inc. All Rights Reserved.

package turret

import (
	"errors"
	"time"

	"github.com/truveris/gousb/usb"
)

const (
	DEVICE_TYPE_THUNDER = iota
	DEVICE_TYPE_CLASSIC = iota
)

var (
	err_no_devices = errors.New("no devices found")
	err_no_led     = errors.New("device has no led")
)

type Turret struct {
	*usb.Device
	Type  int
	Input chan Command
}

func NewTurret(dev *usb.Device) *Turret {
	t := &Turret{Device: dev}
	t.Input = make(chan Command, 64)
	return t
}

func (t *Turret) Command(cmdtype, cmd byte) error {
	if t.Type == DEVICE_TYPE_THUNDER {
		_, err := t.Device.Control(0x21, 0x09, 0, 0, []byte{cmdtype, cmd, 0, 0, 0, 0, 0, 0})
		if err != nil {
			return err
		}
	} else if t.Type == DEVICE_TYPE_CLASSIC {
		if cmdtype != CMD_TYPE_TURRET {
			return err_no_led
		}

		_, err := t.Device.Control(0x21, 0x09, 0x0200, 0, []byte{cmd})
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Turret) QueueCommand(cmdtype, cmd byte, duration time.Duration) {
	t.Input <- Command{Type: cmdtype, Value: cmd, Duration: duration}
}

func (t *Turret) ConsumeCommands() {
	for cmd := range t.Input {
		t.Command(cmd.Type, cmd.Value)
		if cmd.Duration > 0 {
			time.Sleep(cmd.Duration)
		}
	}
}

func (t *Turret) Close() {
	t.Device.Close()
	close(t.Input)
}

func (t *Turret) Light(on bool) {
	var cmd byte

	if on {
		cmd = CMD_LED_ON
	} else {
		cmd = CMD_LED_OFF
	}

	t.QueueCommand(CMD_TYPE_LED, cmd, 0)
}

func (t *Turret) Left(duration time.Duration) {
	t.QueueCommand(CMD_TYPE_TURRET, CMD_TURRET_LEFT, duration)
}

func (t *Turret) Right(duration time.Duration) {
	t.QueueCommand(CMD_TYPE_TURRET, CMD_TURRET_LEFT, duration)
}

func (t *Turret) Up(duration time.Duration) {
	t.QueueCommand(CMD_TYPE_TURRET, CMD_TURRET_UP, duration)
}

func (t *Turret) Down(duration time.Duration) {
	t.QueueCommand(CMD_TYPE_TURRET, CMD_TURRET_DOWN, duration)
}

func (t *Turret) Stop() {
	t.QueueCommand(CMD_TYPE_TURRET, CMD_TURRET_STOP, 0)
}

func (t *Turret) Fire(shots int) {
	for i := 0; i < shots; i++ {
		t.QueueCommand(CMD_TYPE_TURRET, CMD_TURRET_FIRE, 4500)
	}
}
