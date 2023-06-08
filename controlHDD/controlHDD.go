package controlHDD

import (
	"log"

	"github.com/warthog618/gpiod"
)

var GpioSocketInstance *gpiod.Lines = nil

// Available RaspBerry Pi GPIO pins for HDD relays
var AvailableDiskPins [][]int = [][]int{
	{5, 6},
	{27, 22},
	{23, 17},
	{24, 25},
}
var DiskPins [][]int = [][]int{}   //[1][2]int
var DiskList []string = []string{} //[1]string
var DiskStatusList []bool = []bool{}

func findPinsFromDisk(hdd string) []int {
	for i, row := range DiskPins {
		if DiskList[i] == hdd {
			return row
		}
	}
	return nil
}

func getHardDiskRelay(hdd string) bool {
	//hddLines := findPinsFromDisk(hdd)
	//firstHDD, _ := gpiod.RequestLines("gpiochip0", hddLines, gpiod.AsInput)
	if (GpioSocketInstance == nil) {
		print("Error\n")
		return false
	}
	err, _ := GpioSocketInstance.Info()

	if err != nil {
		log.Fatal(err)
	}

	state := err[0].Used

	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = state
		}
	}
	return state
}

// Disk status request

func HardDiskStatusManager(hdd string) bool {
	return getHardDiskRelay(hdd)
}

func HardDiskShutDown(hdd string) {
	hddLines := findPinsFromDisk(hdd)
	//_, err := gpiod.RequestLines("gpiochip0", hddLines, gpiod.AsOutput(1, 1))
	print("HDD? ", hdd, "\nHDD Lines? ", hddLines[0], hddLines[1])
	//err := firstHDD.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveLow))
	if (GpioSocketInstance == nil) {
		print("Error\n")
		return
	}
	GpioSocketInstance.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveLow))

	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = false
		}
	}
}

func HardDiskStartUp(hdd string) {
	hddLines := findPinsFromDisk(hdd)
	//_, err := gpiod.RequestLines("gpiochip0", hddLines, gpiod.AsOutput(0, 0))
	//err := firstHDD.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveHigh))
	if (GpioSocketInstance == nil) {
		print("Error\n")
		return
	}
	GpioSocketInstance.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveHigh))

	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = true
		}
	}
}

func HDDRelaySocketConstructor(hdd string) {
	hddLines := findPinsFromDisk(hdd)
	var err error
	GpioSocketInstance, err = gpiod.RequestLines("gpiochip0", hddLines, gpiod.AsOutput(0, 0))
	if err == nil {
		print("Line initialization error\n")
		GpioSocketInstance = nil
	}
}
