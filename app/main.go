package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/calvernaz/rak811"
	"github.com/tarm/serial"
)

func main()  {
	conf := &serial.Config{
		Name: "/dev/ttyS0",
	}

	lora, err := rak811.New(conf)
	if err != nil {
		log.Fatal("failed to instantiate rak811")
	}

	fmt.Println("Initialise RAK811 module...")
	app_eui := os.Getenv("APP_EUI")
	app_key := os.Getenv("APP_KEY")

	lora.HardReset()
	lora.SetMode(0) // LoRaWAN mode
	lora.SetBand("EU868")
	lora.GetConfig("dev_eui")
	lora.SetConfig(fmt.Sprintf("app_eui=%s,app_key=%s", app_eui, app_key))
	lora.JoinOTAA()
	lora.SetDataRate("5")

	for {
		f, err := os.Open("/sys/class/thermal/thermal_zone0/temp")
		if err != nil {
			break
		}
		r, err := ioutil.ReadAll(f)
		if err != nil {
			continue
		}
		val, err := strconv.ParseFloat(string(r), 32)
		if err != nil {
			continue
		}
		temp := val / 1000

		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, 1)
		binary.BigEndian.PutUint32(buf, 103)
		binary.BigEndian.PutUint32(buf, uint32(temp * 10 + 0.5))

		lora.Send(string(buf))

		time.Sleep(300 * time.Millisecond)
	}

	log.Println("closing...")
	lora.Close()
}
