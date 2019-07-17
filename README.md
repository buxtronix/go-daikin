go-daikin - Daikin AirCon controller
======

This package provides an API for probing and controlling
Daikin split systems fitted with a Wifi controller
module.

This package is tested against the BRP072A42 interface.

Features
=====

 * Discover devices on the local network
 * Query and set current operating parameters
 * Query current sensor values

Basic usage
====

Querying status for discovered devices:

```go
package main

import (
  "fmt"

  "github.com/buxtronix/go-daikin"
  "github.com/golang/glog"
)

func main() {
  d, err := daikin.NewNetwork()
  if err != nil {
    glog.Exit(err)
  }
  if err := d.Discover(); err != nil {
    glog.Exit(err)
  }

  for addr, dev := range d.Devices {
    if err := dev.GetControlInfo(); err != nil {
      glog.Error(err)
      continue
    }
    if err := dev.GetSensorInfo(); err != nil {
      glog.Error(err)
      continue
    }
    fmt.Printf("%s:\n%s", addr, dev.String())
  }
}
```

Set state for a specific device:

```go
package main

import(
  "flag"
  "fmt"

  "github.com/golang/glog"
  daikin "github.com/buxtronix/go-daikin"
)

var (
  address = flag.String("address", "", "Address of daikin unit")
)

func main() {
  flag.Parse()
  d, err := daikin.NewNetwork(daikin.AddressOption(*address))
  if err != nil {
    glog.Exit(err)
  }
  dev := d.Devices[*address]
  if err := dev.GetControlInfo(); err != nil {
    glog.Error(err)
    continue
  }
  if err := dev.GetSensorInfo(); err != nil {
    glog.Error(err)
    continue
  }
  fmt.Printf("%s:\n%s", addr, dev.String())

  dev.ControlInfo.Power = daikin.PowerOn
  dev.ControlInfo.Mode = daikin.ModeHeat
  dev.ControlInfo.Fan = daikin.FanAuto
  dev.ControlInfo.Temperature = 22.5

  if err := dev.SetControlInfo(); err != nil {
    glog.Exit(err)
  }
}
```
