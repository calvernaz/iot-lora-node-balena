package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/TheThingsNetwork/go-cayenne-lib/cayennelpp"
	"github.com/calvernaz/rak811"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

func main()  {
	conf := &serial.Config{
		Name: "/dev/ttyS0",
	}

	lora, err := rak811.New(conf)
	if err != nil {
		fmt.Printf("%s\n", errors.Wrap(err, "failed to create serial connection"))
	}

	fmt.Println("Initialise RAK811 module...")
	appEui := os.Getenv("APP_EUI")
	appKey := os.Getenv("APP_KEY")

	fmt.Printf("AppEUI: %s\n", appEui)
	fmt.Printf("AppKey: %s\n", appKey)
	lora.HardReset()
	fmt.Println("Set LoRaWAN")
	lora.SetMode(0) // LoRaWAN mode
	fmt.Println("Set Band")
	lora.SetBand("EU868")
	fmt.Println("Get DevEUI")
	lora.GetConfig("dev_eui")
	fmt.Println("Set AppEUI,AppKey")
	lora.SetConfig(fmt.Sprintf("app_eui=%s,app_key=%s", appEui, appKey))
	fmt.Println("JoinOTAA")
	lora.JoinOTAA()
	fmt.Println("Set DataRate")
	lora.SetDataRate("5")

	fmt.Println("before loop")
	enc := cayennelpp.NewEncoder()
	for {
		f, err := os.Open("/sys/class/thermal/thermal_zone0/celcius")
		if err != nil {
			fmt.Printf("%s\n", errors.Wrap(err, "open termal file"))
			break
		}
		r, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Printf("%s\n", errors.Wrap(err, "read termal file"))
			continue
		}
		val, err := strconv.ParseFloat(string(r), 32)
		fmt.Printf("value %f", val)
		if err != nil {
			continue
		}
		celcius := val / 1000

		enc.AddTemperature(1, float32(celcius))
		lora.Send(string(enc.Bytes()))
		fmt.Println("sending data")
		time.Sleep(300 * time.Millisecond)
	}

	log.Println("closing...")
	lora.Close()
}
