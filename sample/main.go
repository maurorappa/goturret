// Copyright 2015, Truveris Inc. All Rights Reserved.

package main

import (
	"log"
	"time"

	"github.com/maurorappa/goturret"
	"github.com/truveris/gousb/usb"
)

func main() {
	ctx := usb.NewContext()
	defer ctx.Close()

	turrets, err := turret.Find(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		for _, t := range turrets {
			t.Close()
		}
	}()

	for _, t := range turrets {
		//t.BlinkOn(4)

		t.Left(1 * time.Second)
		t.Up(1 * time.Second)
		t.Right(1 * time.Second)
		t.Down(1 * time.Second)

		t.Light(false)
	}

}
