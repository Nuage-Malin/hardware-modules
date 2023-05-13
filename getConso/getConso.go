package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var hddId = 0
var DiskConsoList [1][4]float64 = [1][4]float64{}

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
	// fmt.Print("connexion arduino\n")
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		// time.Sleep(1000 * time.Millisecond) //par seconde
		if scanner.Err() != nil {
			log.Fatal(err)
		}
		if DiskConsoList[hddId][1] < 3600 {
			str := scanner.Text()
			split_str := strings.Split(str, " ")
			// fmt.Print(split_str[0])
			// fmt.Print("\n")
			value, err := strconv.ParseFloat(split_str[0], 32)
			if err != nil {
				log.Fatal(err)
			} else if DiskConsoList[hddId][3] == 0 {
				time.Sleep(1000 * time.Millisecond) //à voir si nécessaire
				fmt.Print(split_str[0])
				DiskConsoList[hddId][0] += math.Round(value*100) / 100
				DiskConsoList[hddId][1]++
				if DiskConsoList[hddId][1] != 0 {
					DiskConsoList[hddId][2] += (DiskConsoList[hddId][0] / DiskConsoList[hddId][1]) // fonctionnel ?
				}
				fmt.Print("=====================================\n")
				fmt.Print("Total conso sec : ")
				fmt.Print(DiskConsoList[hddId][0])
				fmt.Print("\nNbr secondes (égal à 60 normalement) : ")
				fmt.Print(DiskConsoList[hddId][1])
				fmt.Print("\nTotal conso heure : ")
				fmt.Print(DiskConsoList[hddId][2])
				fmt.Print("\nNbr heures (égal à 0 normalement) : ")
				fmt.Print(DiskConsoList[hddId][3])
				fmt.Print("\n")
				fmt.Print("=====================================\n")
			} else {
				time.Sleep(1000 * time.Millisecond) //à voir si nécessaire
				DiskConsoList[hddId][0] += math.Round(value*100) / 100
				DiskConsoList[hddId][1]++
			}
		} else if DiskConsoList[hddId][1] >= 3600 { //trouver un moyen de quand meme envoyer des kW/h sans attendre directement 1 heure
			DiskConsoList[hddId][2] += (DiskConsoList[hddId][0] / DiskConsoList[hddId][1])
			DiskConsoList[hddId][1] = 0
			DiskConsoList[hddId][0] = 0
			DiskConsoList[hddId][3]++
		}
		// fmt.Print(DiskConsoList[hddId][1])
		// fmt.Print("\n")
		// if DiskConsoList[hddId][1] == 60 {

		// }
	}
}

// arg int -> id hdd
func SendConso(hdd int) float64 {
	return DiskConsoList[hdd][2] / DiskConsoList[hdd][3]
}

// func main() {
// 	c := &serial.Config{Name: findArduino(), Baud: 9600}
// 	s, err := serial.OpenPort(c)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scanner := bufio.NewScanner(s)
// 	// scanner.Split(bufio.ScanWords)
// 	for scanner.Scan() {
// 		fmt.Print("Energy consumntion: ")
// 		// time.Sleep(1000 * time.Millisecond) //à voir si nécessaire
// 		if scanner.Err() != nil {
// 			log.Fatal(err)
// 		}
// 		str := scanner.Text()
// 		// fmt.Print(str)
// 		split_str := strings.Split(str, " ")
// 		fmt.Print(split_str[0])
// 		fmt.Print("\n")
// 	}
// }
