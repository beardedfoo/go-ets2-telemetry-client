package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	ets2 "github.com/cyleriggs/go-ets2-telemetry-client"
	serial "github.com/tarm/serial"
)

var baseUrl = ""
var updateFreq = 0
var serialPort = ""

func init() {
	flag.StringVar(&baseUrl, "baseUrl", "http://localhost:25555", "HTTP url for telemetry server")
	flag.IntVar(&updateFreq, "updateFreq", 1000/30, "Update frequency in milliseconds")
	flag.StringVar(&serialPort, "serialPort", "", "A serial port to write updates to, one var at a time")
}

func main() {
	flag.Parse()

	var err error

	// Open serial port (optional)
	var fSerial io.ReadWriteCloser
	fmt.Printf("Connecting to serial...")
	fSerial, err = serial.OpenPort(&serial.Config{Name: serialPort, Baud: 115200})
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
		os.Exit(1)
	}
	defer fSerial.Close()
	time.Sleep(5 * time.Second)
	fmt.Printf("OK\n")
	if _, err := io.WriteString(fSerial, "\n"); err != nil {
		log.Fatalf("Error writing to serial: %v", err)
		os.Exit(-1)
	}

	log.Printf("reading from serial")
	in := make([]byte, 64)
	if _, err := fSerial.Read(in); err != nil {
		log.Printf("serial read error: %v", err)
		os.Exit(-1)
	}
	log.Printf("serial: %v", in)

	monitor(fSerial, "\n", "")
}

func monitor(fh io.ReadWriteCloser, valueSep string, groupSep string) {
	c := ets2.NewClient(baseUrl)

	// Do the monitoring
	for {
		t, err := c.GetTelemetry()
		if err != nil {
			log.Fatalf("Error reading telemetry data: %v", err)
			os.Exit(-1)
		}
		log.Printf("sending cmds")

		sendCmd(fh, fmt.Sprintf("rpm=%f%s", t.Truck.EngineRpm, valueSep))
		sendCmd(fh, fmt.Sprintf("kmh=%f%s", t.Truck.Speed, valueSep))

		time.Sleep(time.Duration(updateFreq * int(time.Millisecond)))

		log.Printf("reading from serial")
		in := make([]byte, 64)
		if _, err := fh.Read(in); err != nil {
			log.Printf("serial read error: %v", err)
			os.Exit(-1)
		}
		log.Printf("serial: %v", in)
	}
}

func sendCmd(fh io.Writer, cmd string) {
	if _, err := io.WriteString(fh, cmd); err != nil {
		log.Fatalf("Failed to write cmd: %v", err)
		os.Exit(-1)
	}
}
