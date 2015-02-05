// Copyright 2015, Truveris Inc. All Rights Reserved.

package turret

import (
	"time"
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

type Command struct {
	Type  byte
	Value byte
	time.Duration
}
