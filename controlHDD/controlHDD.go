package main

import (
	"log"

	"github.com/warthog618/gpiod"
)

func getHardDiskRelay(hdd string, relay []int) *gpiod.Lines {
	firstHDD, _ := gpiod.RequestLines(hdd, relay, gpiod.AsOutput(0, 1))

	return firstHDD
}

func getInfoFromRelay(hdd *gpiod.Lines) bool {
	err, _ := hdd.Info()

	if err != nil {
		log.Fatal(err)
	}
	return err[0].Used
}

func hardDiskStatusManager(hdd string, relay []int) bool {
	return getInfoFromRelay(getHardDiskRelay(hdd, relay))
}

func hardDiskShutDown(hdd *gpiod.Lines) {
	err := hdd.Reconfigure(gpiod.WithLines([]int{5, 6}, gpiod.AsActiveLow))

	if err != nil {
		log.Fatal(err)
	}
}

func hardDiskStartUp(hdd *gpiod.Lines) {
	err := hdd.Reconfigure(gpiod.WithLines([]int{5, 6}, gpiod.AsActiveHigh))

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	hdd := "HDD1"
	relay := []int{5, 6}

	if hardDiskStatusManager(hdd, relay) {
		hardDiskShutDown(getHardDiskRelay(hdd, relay))
	} else {
		hardDiskStartUp(getHardDiskRelay(hdd, relay))
	}
}
