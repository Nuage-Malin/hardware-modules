package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// var hddId = 0
var DiskConsoList [][]int = [][]int{}

// [0] : conso par seconde incr par seconde
// [1] : nbr de seconde incr par seconde (if statement si 3600 sec)
// [2] : conso par heure totale = DiskConsoList[0] / DiskConsoList[1]
// [3] : nbr d'heures totale

// return DiskConsoList[2] / DiskConsoList[3] pour avoir conso par heure en kW/h dans une autre fonction

func findArduino() string {
	contents, _ := ioutil.ReadDir("/dev")

	for _, f := range contents {
		if strings.Contains(f.Name(), "ttyUSB") ||
			strings.Contains(f.Name(), "ttyACM0") {
			fmt.Printf("Linked to Arduino: %s\n", f.Name())
			return "/dev/" + f.Name()
		}
	}
	return ""
}

func main() {
	c := &serial.Config{Name: findArduino(), Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("connexion arduino")
	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		time.Sleep(1000 * time.Millisecond) //par seconde
		if scanner.Err() != nil {
			log.Fatal(err)
		}
		if DiskConsoList[0][1] < 3600 {
			value, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			} else {
				DiskConsoList[0][0] += value
				DiskConsoList[0][1]++
			}
		} else if DiskConsoList[0][1] >= 3600 {
			DiskConsoList[0][2] += (DiskConsoList[0][0] / DiskConsoList[0][1])
			DiskConsoList[0][1] = 0
			DiskConsoList[0][0] = 0
			DiskConsoList[0][3]++
		}
		fmt.Print(DiskConsoList[0][0])
	}
}

// arg int -> id hdd
func SendConso(hdd int) int {
	return DiskConsoList[hdd][2] / DiskConsoList[hdd][3]
}

// func main() {
// 	c := &serial.Config{Name: FindArduino(), Baud: 9600}
// 	s, err := serial.OpenPort(c)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scanner := bufio.NewScanner(s)
// 	scanner.Split(bufio.ScanWords)
// 	for scanner.Scan() {
// 		fmt.Print("Energy consumntion: ")
// 		time.Sleep(1000 * time.Millisecond) //à voir si nécessaire
// 		if scanner.Err() != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(scanner.Text())
// 	}
// }
