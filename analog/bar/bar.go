package main

import (
	"fmt"
	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/platforms/firmata"
)

const startMsg = "\n*** Change input to analog pin %s to generate events. ***\n"
const inputPin = "0"

var firstEvent = true

func handleChange(data interface{}) {
	if firstEvent {
		fmt.Printf(startMsg, inputPin)
		firstEvent = false
	}
	value := data.(int)
	// ticks := int(float32(value) / 1024 * 70)
	ticks := int(gobot.ToScale(gobot.FromScale(float64(value), 0, 1024), 0, 70))
	bar := strings.Repeat("#", ticks)
	fmt.Printf("%4d %s\n", value, bar)
}

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/cu.usbmodem1411")
	sensor := aio.NewAnalogSensorDriver(firmataAdaptor, inputPin, 10)
	work := func() {
		sensor.On(sensor.Event("data"), handleChange)
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor},
		work,
	)
	robot.Start()
}
