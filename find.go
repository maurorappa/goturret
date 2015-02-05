// Copyright 2015, Truveris Inc. All Rights Reserved.

package turret

import (
	"github.com/truveris/gousb/usb"
)

func Find() ([]*Turret, error) {
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
		return turrets, err
	}

	if len(devs) == 0 {
		return turrets, err_no_devices
	}

	for _, dev := range devs {
		t := NewTurret(dev)

		if dev.Descriptor.Vendor == 0x2123 && dev.Descriptor.Product == 0x1010 {
			t.Type = DEVICE_TYPE_THUNDER
		} else if dev.Descriptor.Vendor == 0x0a81 && dev.Descriptor.Product == 0x0701 {
			t.Type = DEVICE_TYPE_CLASSIC
		}

		turrets = append(turrets, t)
	}

	return turrets, nil
}
