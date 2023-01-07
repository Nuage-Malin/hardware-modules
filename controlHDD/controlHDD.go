package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/warthog618/gpiod"
)

func main() {
	firstRelay, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	secondRelay, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	v := 0
	ll, err := gpiod.RequestLines("gpiochip0", []int{firstRelay, secondRelay}, gpiod.AsOutput(v, v))
	if err != nil {
		panic(err)
	}
	defer func() {
		ll.Reconfigure(gpiod.AsInput)
		ll.Close()
	}()
	values := map[int]string{0: "inactive", 1: "active"}
	fmt.Printf("Set pin %d %s\n", firstRelay, values[v])
	fmt.Printf("Set pin %d %s\n", secondRelay, values[v])

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	for {
		select {
		case <-time.After(2 * time.Second):
			v ^= 1
			ll.SetValues([]int{v, v})
			fmt.Printf("Set pin %d %s\n", firstRelay, values[v])
			fmt.Printf("Set pin %d %s\n", secondRelay, values[v])
		case <-quit:
			return
		}
	}
}
