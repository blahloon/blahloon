package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
	"github.com/nats-io/nats"
)

func main() {

	gbot := gobot.NewGobot()
	e := edison.NewEdisonAdaptor("edison")
	blueLed := gpio.NewGroveLedDriver(e, "led", "4")
	greenLed := gpio.NewGroveLedDriver(e, "led", "2")
	redLed := gpio.NewGroveLedDriver(e, "led", "3")

	log.Printf("Default URL: [%s]\n", nats.DefaultURL)

	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	log.Printf("Connecting to %s ...\n", *urls)

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	work := func() {
		var i int
		subj := "hello"

		nc.Subscribe(subj, func(msg *nats.Msg) {
			i++
			printMsg(msg, i)

			m := string(msg.Data)
			switch {
			case m == "blue-on":
				log.Println("Turning on blue LED...")
				blueLed.On()
			case m == "blue-off":
				log.Println("Turning off blue LED...")
				blueLed.Off()
			case m == "green-on":
				log.Println("Turning on green LED...")
				greenLed.On()
			case m == "green-off":
				log.Println("Turning off green LED...")
				greenLed.Off()
			case m == "red-on":
				log.Println("Turning on red LED...")
				redLed.On()
			case m == "red-off":
				log.Println("Turning off red LED...")
				redLed.Off()
			}
		})
		nc.Flush()

		log.Printf("Listening on [%s]\n", subj)

		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		}
	}

	blueLedBot := gobot.NewRobot("blueLedBot",
		[]gobot.Connection{e},
		[]gobot.Device{blueLed},
		work,
	)
	gbot.AddRobot(blueLedBot)

	greenLedBot := gobot.NewRobot("greenLedBot",
		[]gobot.Connection{e},
		[]gobot.Device{greenLed},
		work,
	)
	gbot.AddRobot(greenLedBot)

	redLedBot := gobot.NewRobot("redLedBot",
		[]gobot.Connection{e},
		[]gobot.Device{redLed},
		work,
	)
	gbot.AddRobot(redLedBot)

	gbot.Start()

	runtime.Goexit()
}

func usage() {
	log.Fatalf("Usage: blinkit [-s server]\n")
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}
