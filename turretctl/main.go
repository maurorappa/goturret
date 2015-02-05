// Copyright 2015, Truveris Inc. All Rights Reserved.

package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/truveris/goturret"
	"github.com/truveris/gousb/usb"
)

func usage() {
	log.Fatal("usage: turretctl cmd [value]")
}

func getShots() int {
	var shots int

	if len(os.Args) == 3 {
		shots, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		if shots > 4 {
			shots = 4
		}
	} else {
		shots = 1
	}
	return shots
}

func getDuration() time.Duration {
	if len(os.Args) == 3 {
		ms, err := time.ParseDuration(os.Args[2] + "ms")
		if err != nil {
			log.Fatal(err)
		}
		return ms
	}
	return 0
}

func getBoolean() bool {
	if len(os.Args) == 3 {
		if os.Args[2] == "on" {
			return true
		}
	}
	return false
}

func main() {
	var cmd string

	if len(os.Args) < 2 || len(os.Args) > 3 {
		usage()
	}

	cmd = os.Args[1]

	ctx := usb.NewContext()
	defer ctx.Close()

	turrets, err := turret.Find(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer func() {
		for _, t := range turrets {
			t.Close()
		}
	}()

	for _, t := range turrets {
		switch cmd {
		case "left":
			duration := getDuration()
			t.Left(0)
			time.Sleep(duration)
			t.Stop()
		case "right":
			duration := getDuration()
			t.Right(0)
			time.Sleep(duration)
			t.Stop()
		case "up":
			duration := getDuration()
			t.Up(0)
			time.Sleep(duration)
			t.Stop()
		case "down":
			duration := getDuration()
			t.Down(0)
			time.Sleep(duration)
			t.Stop()
		case "light":
			light := getBoolean()
			t.Light(light)
		case "fire":
			shots := getShots()
			t.Fire(shots)
		case "reset":
			t.Reset()
		default:
			log.Fatal("error: unknown command")
		}
	}

	for _, t := range turrets {
		t.Shutdown()
	}
}
