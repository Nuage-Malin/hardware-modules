package getConso

import (
	"bufio"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kraken-hpc/go-fork"
	"github.com/tarm/serial"
)

// var hddId = 0
var DiskConsoList [1][4]int = [1][4]int{} //voir pour moduler taille list

// [0] : conso par seconde incr par seconde
// [1] : nbr de seconde incr par seconde (if statement si 3600 sec)
// [2] : conso par heure totale = DiskConsoList[0] / DiskConsoList[1]
// [3] : nbr d'heures totale

// return DiskConsoList[2] / DiskConsoList[3] pour avoir conso par heure en kW/h dans une autre fonction

func findArduino() string {
	contents, _ := ioutil.ReadDir("/dev")

	for _, f := range contents {
		if strings.Contains(f.Name(), "ttyUSB") || strings.Contains(f.Name(), "ttyACM0") {
			return "/dev/" + f.Name()
		}
	}
	return ""
}

func child(hddId int) {
	c := &serial.Config{Name: findArduino(), Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatal(err)
		}
		if DiskConsoList[hddId][1] < 3600 {
			str := scanner.Text()
			split_str := strings.Split(str, " ")
			value, err := strconv.ParseFloat(split_str[0], 32)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(1000 * time.Millisecond) //à voir si nécessaire
			DiskConsoList[hddId][0] = int(value)
			DiskConsoList[hddId][1]++
			// DiskConsoList[hddId][2] += (DiskConsoList[hddId][0] / DiskConsoList[hddId][1]) // fonctionnel ?
		} else if DiskConsoList[hddId][1] >= 3600 { //trouver un moyen de quand meme envoyer des kW/h sans attendre directement 1 heure
			//DiskConsoList[hddId][2] += (DiskConsoList[hddId][0] / DiskConsoList[hddId][1])
			//DiskConsoList[hddId][1] = 0
			DiskConsoList[hddId][1] = 0
			DiskConsoList[hddId][3]++
		}
	}
}

func GetConso(hddId int) {
	fork.RegisterFunc("child", child)
	fork.Init()
	if err := fork.Fork("child", 1); err != nil {
		log.Fatalf("failed to fart: %v", err)
	}
}

func SendConso(hddId int) int {
	// if DiskConsoList[hddId][3] < 2 {
	// 	return DiskConsoList[hddId][2]
	// } else {
	// 	return DiskConsoList[hddId][2] / DiskConsoList[hddId][3]
	// }
	return DiskConsoList[hddId][0]
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
