package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
)

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
	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Print("Energy consumntion: ")
		time.Sleep(1000 * time.Millisecond) //à voir si nécessaire
		if scanner.Err() != nil {
			log.Fatal(err)
		}
		fmt.Println(scanner.Text())
	}
}
