package main

import (
	"fmt"
	rpio "github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

var (
	// Use mcu pin 22, corresponds to GPIO 3 on the pi
	pin = rpio.Pin(17)
)

func main() {
	fmt.Println("start")

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge) // enable falling edge event detection

	fmt.Println("press a button")

	for  {
		if pin.EdgeDetected() { // check if event occured
			fmt.Println("button pressed")
		}
		time.Sleep(100* time.Millisecond)
	}

	pin.Detect(rpio.NoEdge) // disable edge event detection


	fmt.Println("ende")
}
