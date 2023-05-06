package controlHDD

import (
	"log"

	"github.com/warthog618/gpiod"
)

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
	for i, hddtmp := range DiskList {
		print(hddtmp)
		if hddtmp == hdd {
			return DiskStatusList[i]
		}
	}
	firstHDD, _ := gpiod.RequestLines(hdd, findPinsFromDisk(hdd), gpiod.AsOutput(0, 1))

	err, _ := firstHDD.Info()

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
	for i, hddtmp := range DiskList {
		if hddtmp == hdd {
			DiskStatusList[i] = false
			return
		}
	}
	firstHDD, _ := gpiod.RequestLines(hdd, findPinsFromDisk(hdd), gpiod.AsOutput(0, 1))

	hddLines := findPinsFromDisk(hdd)
	err := firstHDD.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveLow))

	if err != nil {
		log.Fatal(err)
	}
	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = false
		}
	}
}

func HardDiskStartUp(hdd string) {
	for i, hddtmp := range DiskList {
		if hddtmp == hdd {
			DiskStatusList[i] = true
			return
		}
	}
	firstHDD, _ := gpiod.RequestLines(hdd, findPinsFromDisk(hdd), gpiod.AsOutput(0, 1))

	hddLines := findPinsFromDisk(hdd)
	err := firstHDD.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveHigh))

	if err != nil {
		log.Fatal(err)
	}
	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = true
		}
	}
}
