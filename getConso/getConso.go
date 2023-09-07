package getConso

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// var hddId = 0
var DiskConsoList [1]int = [1]int{} //voir pour moduler taille list

const (
	IdlePower      = 30.0 // Watts in idle mode
	ReadWritePower = 60.0 // Watts during read/write operations
)

// [0] : conso par seconde incr par seconde
// [1] : nbr de seconde incr par seconde (if statement si 3600 sec)
// [2] : conso par heure totale = DiskConsoList[0] / DiskConsoList[1]
// [3] : nbr d'heures totale

// return DiskConsoList[2] / DiskConsoList[3] pour avoir conso par heure en kW/h dans une autre fonction

func findArduino() string {
	outputDirRead, _ := os.Open("/dev")
	outputDirFiles, _ := outputDirRead.ReadDir(0)
	for outputIndex := range outputDirFiles {
		outputFileHere := outputDirFiles[outputIndex]
		// Get name of file.
		f := outputFileHere.Name()
		if strings.Contains(f, "ttyUSB") || strings.Contains(f, "ttyACM0") {
			return "/dev/" + f
		}
	}
	return ""
}

func GetConso(hddId int) {
	var (
		totalEnergy  float64
		startTime    time.Time
		newState     string
		currentState string
	)

	c := &serial.Config{Name: findArduino(), Baud: 9600}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
	currentState = "idle"

	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatal(err)
		}
		str := scanner.Text()
		split_str := strings.Split(str, " ")
		value, err := strconv.ParseFloat(split_str[3], 32)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1000 * time.Millisecond)
		DiskConsoList[hddId] = int(value)
		fmt.Println("debug: ", int(value))

		// State detection
		if int(value) >= IdlePower {
			newState = "readwrite"
		} else {
			newState = "idle"
		}

		if newState != currentState {
			// Calculate the time spent in the current state
			elapsedTime := time.Since(startTime)

			// Update the total energy consumption based on the previous state
			if currentState == "idle" {
				totalEnergy += float64(elapsedTime) / float64(time.Second) * value
				DiskConsoList[hddId] = int(totalEnergy)
			} else if currentState == "readwrite" {
				totalEnergy += float64(elapsedTime) / float64(time.Second) * value
				DiskConsoList[hddId] = int(totalEnergy)
			}

			// Update the current state
			currentState = newState

			// Update the start time to mark the beginning of the new state
			startTime = time.Now()
		}

		// Print the current total energy consumption
		fmt.Printf("Total Energy Consumption: %.2f Watt-seconds\n", totalEnergy)

		// Sleep for a while to control the measurement rate
		time.Sleep(time.Second)
	}
}

func SendConso(hddId int) int {
	// if DiskConsoList[hddId][3] < 2 {
	// 	return DiskConsoList[hddId][2]
	// } else {
	// 	return DiskConsoList[hddId][2] / DiskConsoList[hddId][3]
	// }
	return DiskConsoList[hddId]
}
