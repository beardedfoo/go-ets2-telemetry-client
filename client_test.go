package ets2

import "io/ioutil"
import "os"
import "testing"

func TestParse(t *testing.T) {
	f, _ := os.Open("testdata/disconnected.json")
	defer f.Close()
	json, _ := ioutil.ReadAll(f)
	data, err := parseTelemetry(json)
	if err != nil {
		t.Fatal(err)
	}

	if data.Game.Connected != false {
		t.Fatalf("data.Game.Connected != false")
	}
}
