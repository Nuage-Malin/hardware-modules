package controlHDD

import (
	"log"
	"fmt"

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

/**
 * @brief Hibernate or wake up the computer to which disks are connected.
 *
 * The function calls the server-manager `hsmlib` functions.
 * It either wakes up the computer using WoLAN magic packet or hibernates it
 * using the `systemctl` command.
 *
 * @param isAsleep The state the computer must be changed to.
 */
func findPinsFromDisk(hdd string) []int {
	for i, row := range DiskPins {
		if DiskList[i] == hdd {
			return row
		}
	}
	return nil
}

/**
 * @brief Hibernate or wake up the computer to which disks are connected.
 *
 * The function calls the server-manager `hsmlib` functions.
 * It either wakes up the computer using WoLAN magic packet or hibernates it
 * using the `systemctl` command.
 *
 * @param isAsleep The state the computer must be changed to.
 */
func getHardDiskRelay(hdd string) bool {
	if (GpioSocketInstance == nil) {
		print("Error getHardDiskRelay\n")
		return false
	}
	
	gpioLine, err := GpioSocketInstance.Info()

	if err != nil {
		print("getHardDiskRelay gpio library request error")
		log.Fatal(err)
	}

	state := gpioLine[0].Used

	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = state
		}
	}
	return state
}

// Disk status request
/**
 * @brief Hibernate or wake up the computer to which disks are connected.
 *
 * The function calls the server-manager `hsmlib` functions.
 * It either wakes up the computer using WoLAN magic packet or hibernates it
 * using the `systemctl` command.
 *
 * @param isAsleep The state the computer must be changed to.
 */
func HardDiskStatusManager(hdd string) bool {
	return getHardDiskRelay(hdd)
}

/**
 * @brief Hibernate or wake up the computer to which disks are connected.
 *
 * The function calls the server-manager `hsmlib` functions.
 * It either wakes up the computer using WoLAN magic packet or hibernates it
 * using the `systemctl` command.
 *
 * @param isAsleep The state the computer must be changed to.
 */
func HardDiskShutDown(hdd string) {
	hddLines := findPinsFromDisk(hdd)
	//_, err := gpiod.RequestLines("gpiochip0", hddLines, gpiod.AsOutput(1, 1))
	print("HDD? ", hdd, "\nHDD Lines? ", hddLines[0], hddLines[1])
	//err := firstHDD.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveLow))
	if (GpioSocketInstance == nil) {
		print("Error HardDiskShutDown\n")
		return
	}
	for i, linePrint := range hddLines {
		print(linePrint[i])
	}
	GpioSocketInstance.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveLow))

	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = false
		}
	}
}

/**
 * @brief .
 *
 * The function calls the server-manager `hsmlib` functions.
 * It either wakes up the computer using WoLAN magic packet or hibernates it
 * using the `systemctl` command.
 *
 * @param isAsleep The state the computer must be changed to.
 */
func HardDiskStartUp(hdd string) error {
	hddLines := findPinsFromDisk(hdd)

	if (GpioSocketInstance == nil) {
		print("Error HardDiskStartUp\n")
		return
	}
	for i, linePrint := range hddLines {
		print(linePrint[i])
	}
	GpioSocketInstance.Reconfigure(gpiod.WithLines(hddLines, gpiod.AsActiveHigh))

	for i, hddName := range DiskList {
		if hddName == hdd {
			DiskStatusList[i] = true
		}
	}
}

/**
 * @brief Constructor for gpioLines of each HDD.
 *
 * The function initializes a gpioLines object for each HDD.
 * Available gpio pin tuples are fetched from `AvailableDiskPins` depending on HDD name.
 * In our current implementation, Raspberry Pi 3b+ allows only 4 HDD to be controlled.
 *
 * @param hdd The HDD to initialize the line for
 */
func HDDRelaySocketConstructor(hdd string) error {
	hddLines := findPinsFromDisk(hdd)
	var err error
	GpioSocketInstance, err = gpiod.RequestLines("gpiochip0", hddLines, gpiod.AsOutput(0, 0))
	if err != nil {
		return fmt.Errorf("HDDRelaySocketConstructor: Line initialization error\n%s", err.Error())
	}
	return nil
}