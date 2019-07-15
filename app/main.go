// Package main is a command line tool for querying and controlling
// Daikin AC/Heatpump units.
package main

import (
	"flag"
	"fmt"
	"github.com/buxtronix/go-daikin"
	"github.com/golang/glog"
)

var (
	ifName  = flag.String("interface", "", "Interface to scan on")
	address = flag.String("address", "", "Use device at specific address")

	powerOn  = flag.Bool("on", false, "Turn unit on")
	powerOff = flag.Bool("off", false, "Turn unit off")

	modeHeat = flag.Bool("heat", false, "Set to heating mode")
	modeCool = flag.Bool("cool", false, "Set to cooling mode")
	modeFan  = flag.Bool("fan", false, "Set to fan mode")

	fanRate = flag.String("speed", "", "Fan speed (A, B, 1, 2, 3, 4, 5)")

	fanVertical   = flag.Bool("vertical", false, "Sweep louvres vertically")
	fanHorizontal = flag.Bool("horizontal", false, "Sweep louvres horizontally")

	setTemp = flag.Float64("temp", 22.0, "Temperature to set to")
)

func main() {
	flag.Parse()
	d, err := daikin.NewNetwork(
		daikin.InterfaceOption(*ifName),
		daikin.AddressOption(*address))
	if err != nil {
		glog.Exit(err)
	}
	if err := d.Discover(); err != nil {
		glog.Exit(err)
	}

	fmt.Printf("Devices:\n")
	for a, d := range d.Devices {
		if err := d.GetControlInfo(); err != nil {
			glog.Error(err)
			continue
		}
		if err := d.GetSensorInfo(); err != nil {
			glog.Error(err)
			continue
		}
		fmt.Printf("Current %s:\n%s\n\n", a, d)
		if *powerOn || *powerOff {
			if *powerOn {
				d.ControlInfo.Power = daikin.PowerOn
			}
			if *powerOff {
				d.ControlInfo.Power = daikin.PowerOff
			}
			if *modeHeat {
				d.ControlInfo.Mode = daikin.ModeHeat
			}
			if *modeCool {
				d.ControlInfo.Mode = daikin.ModeCool
			}
			if *modeFan {
				d.ControlInfo.Mode = daikin.ModeFan
			}

			switch *fanRate {
			case "A":
				d.ControlInfo.Fan = daikin.FanAuto
			case "B":
				d.ControlInfo.Fan = daikin.FanSilent
			case "1":
				d.ControlInfo.Fan = daikin.Fan1
			case "2":
				d.ControlInfo.Fan = daikin.Fan2
			case "3":
				d.ControlInfo.Fan = daikin.Fan3
			case "4":
				d.ControlInfo.Fan = daikin.Fan4
			case "5":
				d.ControlInfo.Fan = daikin.Fan5
			case "":
				// Noop.
			default:
				glog.Exitf("Unsupported fan rate: %s", *fanRate)
			}

			switch {
			case *fanHorizontal && *fanVertical:
				d.ControlInfo.FanDir = daikin.FanDirBoth
			case *fanVertical:
				d.ControlInfo.FanDir = daikin.FanDirVertical
			case *fanHorizontal:
				d.ControlInfo.FanDir = daikin.FanDirHorizontal
			default:
				d.ControlInfo.FanDir = daikin.FanDirStopped
			}

			if *setTemp > 0 {
				d.ControlInfo.Temperature = daikin.Temperature(*setTemp)
			}
			fmt.Printf("Setting to new values:\n%s\n\n", d)

			if err := d.SetControlInfo(); err != nil {
				glog.Exitf("Error setting aircon: %v", err)
			}

			if err := d.GetControlInfo(); err != nil {
				glog.Exitf("Error getting aircon data: %v", err)
			}
			if err := d.GetSensorInfo(); err != nil {
				glog.Exitf("Error getting aircon data: %v", err)
			}
			fmt.Printf("New values %s:\n%s\n\n", a, d)
		}
	}
}
