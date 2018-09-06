package main

import (
	"fmt"
	rpio "github.com/stianeikeland/go-rpio"
	"os"
	"time"
	"github.com/sisi485/ring/hue"
)


//const (
//	AWAITTIME      = 800
//	TRANSITIONTIME = float64(AWAITTIME) / 1000 / 0.1
//	//SENDMSGTO 	   = "491787807924"
//	SENDMSGTO 	   = "4915154689867"
//	MSGTEXT 	   = "mhhh et klingelt bruder"
//)


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
	pin.Detect(rpio.AnyEdge) // enable falling edge event detection

	fmt.Println("press a button")

	for  {
		if pin.Read() == 0 {
			fmt.Printf("i channged")
			ring()
		}
		time.Sleep(128 * time.Millisecond)
	}

	pin.Detect(rpio.NoEdge) // disable edge event detection

	fmt.Println("ende")
}

func ring() {

	lights, err := hue.GetLights()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not get lgihts: %v\n", err)
		return
	}

	keys := make([]string, 0, len(lights))
	for k := range lights {
		keys = append(keys, k)
	}

	scene, err := hue.CreateScene("test", keys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not create scene: %v\n", err)
		return
	}

	lightsColor := map[string]interface{}{
		"on":             true,
		"sat":            255,
		"hue":            34000,
		"bri":            255,
		"transitiontime": TRANSITIONTIME,
	}

	lightsOn := map[string]interface{}{
		"on":             true,
		"sat":            255,
		"hue":            34000,
		"bri":            255,
		"transitiontime": TRANSITIONTIME,
	}

	lightsOff := map[string]interface{}{
		"bri":            0,
		"transitiontime": TRANSITIONTIME,
	}

	id := scene[0].(map[string]interface{})["success"].(map[string]interface{})["id"].(string)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not convert id: %v\n", err)
		return
	}

	lightsDefault := map[string]interface{}{
		"scene":          id,
		"transitiontime": TRANSITIONTIME,
	}

	hue.SetGroup(1, lightsOff)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsColor)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsOff)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsOn)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsOff)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsDefault)

	_, err = hue.DeleteScene(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not delete scene: %v\n", err)
		return
	}

	fmt.Println("im done, closing..")
}
