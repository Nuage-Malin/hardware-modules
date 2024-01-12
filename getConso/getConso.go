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

var DiskConsoList [1]float64 = [1]float64{} //voir pour moduler taille list

var (
	HddIdlePower      = 30.0	// Watts in idle mode
	HddReadWritePower = 60.0	// Watts during read/write operations
	ServerPower		  = 100.0	// Watts the machine alone draws
	NoisePower		  = 2.0		// Watts measured when everything is shut down
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

func GetConso(hddId int, bench bool) {
	var (
		startTime    time.Time
		// currentState string
		// newState     string
	)

	c := &serial.Config{Name: findArduino(), Baud: 9600}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
	// currentState = "idle"
	startTime = time.Now()

	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatal(err)
		}
		str := scanner.Text()
		print(str)
		split_str := strings.Split(str, " ")
		value, err := strconv.ParseFloat(split_str[3], 32)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("debug: ", float64(value))

		if (!bench) {
			// State detection
			// if float64(value) >= HddIdlePower {
				// newState = "readwrite"
			// } else {
				// newState = "idle"
			// }

			// Update the total energy consumption based on the current state
			// if newState != currentState {
				// Calculate the energy consumed during the previous state
				// elapsedTime := time.Since(startTime)
				// fmt.Println(elapsedTime)
				// fmt.Println(time.Second)
				// if currentState == "idle" {
					// DiskConsoList[hddId] += float64(elapsedTime) / float64(time.Second) * HddIdlePower
				// } else if currentState == "readwrite" {
					// DiskConsoList[hddId] += float64(elapsedTime) / float64(time.Second) * HddReadWritePower
				// }

				// Update the current state
				// currentState = newState

				// Update the start time to mark the beginning of the new state
				// startTime = time.Now()

				// Print the current total energy consumption
				DiskConsoList[hddId] = float64(value)
				fmt.Printf("Total Energy Consumption: %.2f Watt-seconds\n", DiskConsoList[hddId])
			//}
		} else {
			// Give 20 seconds to the device to stabilize it's consumption at each step of the measurment
			elapsedTime := time.Since(startTime)
			switch {
			case (elapsedTime >= 20 && elapsedTime <= 30):
				HddReadWritePower += float64(elapsedTime) - 20 / float64(time.Second) * value
			case (elapsedTime >= 50 && elapsedTime <= 60):
				HddIdlePower += float64(elapsedTime) - 50 / float64(time.Second) * value
			case (elapsedTime >= 80 && elapsedTime <= 90):
				ServerPower += float64(elapsedTime) - 80 / float64(time.Second) * value
			case (elapsedTime >= 110 && elapsedTime <= 120):
				NoisePower += float64(elapsedTime) - 110 / float64(time.Second) * value
			case (elapsedTime > 120):
				return
			}
			
		}
		// Sleep for a while to control the measurement rate
		time.Sleep(time.Second)
	}
}

func SendConso(hddId int) float64 {
	// if DiskConsoList[hddId][3] < 2 {
	// 	return DiskConsoList[hddId][2]
	// } else {
	// 	return DiskConsoList[hddId][2] / DiskConsoList[hddId][3]
	// }
	return DiskConsoList[hddId]
}
