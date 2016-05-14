package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	ets2		"github.com/cyleriggs/go-ets2-telemetry-client"
	serial 	"github.com/goburrow/serial"
)

var baseUrl = ""
var updateFreq = 0
var serialPort = ""

func init() {
	flag.StringVar(&baseUrl, "baseUrl", "http://localhost:25555", "HTTP url for telemetry server")
	flag.IntVar(&updateFreq, "updateFreq", 1000/24, "Update frequency in milliseconds")
	flag.StringVar(&serialPort, "serialPort", "COM3", "A serial port to write updates to, one var at a time")
}

func main() {
	flag.Parse()

	var err error

	// Open serial port (optional)
	fmt.Printf("Connecting to serial...")
	serialConfig := &serial.Config{
		Address: serialPort, BaudRate: 115200, StopBits: 1,
		Timeout: 100 * time.Millisecond, Parity: "N"}
	fSerial, err := serial.Open(serialConfig)
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
		sendCmd(fh, fmt.Sprintf("rpm=%f%s", t.Truck.EngineRpm, valueSep))
		sendCmd(fh, fmt.Sprintf("kmh=%f%s", t.Truck.Speed, valueSep))
		sendCmd(fh, fmt.Sprintf("fuel=%f%s", t.Truck.FuelCapacity/t.Truck.Fuel, valueSep))

		time.Sleep(time.Duration(updateFreq * int(time.Millisecond)))

	}
}

func sendCmd(fh io.ReadWriter, cmd string) {
	if _, err := io.WriteString(fh, cmd); err != nil {
		log.Fatalf("Failed to write cmd: %v", err)
		os.Exit(-1)
	} else {
		scanner := bufio.NewScanner(fh)
		scanner.Scan()
		log.Printf("serial: %v", scanner.Text())
		scanner.Scan()
		log.Printf("serial: %v", scanner.Text())
	}
}
