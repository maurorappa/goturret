// Copyright 2015, Truveris Inc. All Rights Reserved.

package turret

import (
	"errors"

	"github.com/truveris/gousb/usb"
)

var (
	errNoDevices = errors.New("no devices found")
)

// Find returns all the turrets we could find on the system. You must call
// Close() or Shutdown() for all the turrets returned from this function.
func Find(ctx *usb.Context) ([]*Turret, error) {
	var turrets []*Turret

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

	// ListDevices can occaionally fail, so be sure to check its return value.
	if err != nil {
		return turrets, err
	}

	if len(devs) == 0 {
		return turrets, errNoDevices
	}

	for _, dev := range devs {
		t := NewTurret(dev)

		if dev.Descriptor.Vendor == 0x2123 && dev.Descriptor.Product == 0x1010 {
			t.Type = DeviceTypeThunder
		} else if dev.Descriptor.Vendor == 0x0a81 && dev.Descriptor.Product == 0x0701 {
			t.Type = DeviceTypeClassic
		}

		turrets = append(turrets, t)
	}

	return turrets, nil
}
