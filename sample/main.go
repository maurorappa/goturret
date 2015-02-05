// Copyright 2015, Truveris Inc. All Rights Reserved.

package main

import (
	"log"
	"time"

	"github.com/truveris/goturret"
)

func main() {
	turrets, err := turret.Find()
	if err != nil {
		log.Fatal(err)
	}

	for _, turret := range turrets {
		turret.Light(true)

		turret.Left(1 * time.Second)
		turret.Up(1 * time.Second)
		turret.Right(1 * time.Second)
		turret.Down(1 * time.Second)

		turret.Light(false)

		turret.Close()
	}
}
