// Copyright 2015, Truveris Inc. All Rights Reserved.

package main

import (
	"log"
	"time"

	"github.com/truveris/gousb/usb"
)

const (
	DEVICE_TYPE_THUNDER = iota
	DEVICE_TYPE_CLASSIC = iota
)

const (
	CMD_TYPE_TURRET = 0x02
	CMD_TYPE_LED    = 0x03
)

const (
	CMD_TURRET_DOWN  = 0x01
	CMD_TURRET_UP    = 0x02
	CMD_TURRET_LEFT  = 0x04
	CMD_TURRET_RIGHT = 0x08
	CMD_TURRET_FIRE  = 0x10
	CMD_TURRET_STOP  = 0x20
)

const (
	CMD_LED_OFF = 0x00
	CMD_LED_ON  = 0x01
)

type Turret struct {
	*usb.Device
	Type int
}

func (t *Turret) Command(cmdtype, cmd byte) error {
	_, err := t.Device.Control(0x21, 0x09, 0, 0, []byte{cmdtype, cmd, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	turrets := make([]*Turret, 0)

	ctx := usb.NewContext()
	defer ctx.Close()

	// ListDevices is used to find the devices to open.
	devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
		if desc.Vendor == 0x2123 && desc.Product == 0x1010 {
			return true
		}

		if desc.Vendor == 0x0a81 && desc.Product == 0x0701 {
			return true
		}

		return false
	})

	// All Devices returned from ListDevices must be closed.
	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	// ListDevices can occaionally fail, so be sure to check its return value.
	if err != nil {
		log.Fatalf("list: %s", err)
	}

	if len(devs) == 0 {
		log.Fatalf("no devices found")
	}

	for _, dev := range devs {
		t := &Turret{Device: dev}

		if dev.Descriptor.Vendor == 0x2123 && dev.Descriptor.Product == 0x1010 {
			t.Type = DEVICE_TYPE_THUNDER
		} else if dev.Descriptor.Vendor == 0x0a81 && dev.Descriptor.Product == 0x0701 {
			t.Type = DEVICE_TYPE_CLASSIC
		}

		turrets = append(turrets, t)
	}

	for _, turret := range turrets {
		turret.Command(CMD_TYPE_LED, CMD_LED_ON)
		turret.Command(CMD_TYPE_TURRET, CMD_TURRET_LEFT)
		time.Sleep(1 * time.Second)
		turret.Command(CMD_TYPE_TURRET, CMD_TURRET_RIGHT)
		time.Sleep(1 * time.Second)
		turret.Command(CMD_TYPE_TURRET, CMD_TURRET_STOP)
		turret.Command(CMD_TYPE_LED, CMD_LED_OFF)
	}
}
