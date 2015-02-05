// Copyright 2015, Truveris Inc. All Rights Reserved.

package turret

import (
	"time"
)

// Type of commands (control different parts of the turret).
const (
	CmdTypeTurret = 0x02
	CmdTypeLight  = 0x03
	CmdTypeDone   = 0x99
)

// All the turret (motor) commands.
const (
	CmdTurretDown  = 0x01
	CmdTurretUp    = 0x02
	CmdTurretLeft  = 0x04
	CmdTurretRight = 0x08
	CmdTurretFire  = 0x10
	CmdTurretStop  = 0x20
)

// All the light/LED commands.
const (
	CmdLightOff = 0x00
	CmdLightOn  = 0x01
)

// Command contains all the command information needed to be passed to a
// turret.
type Command struct {
	Type  byte
	Value byte
	time.Duration
}
