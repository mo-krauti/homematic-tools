package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mo-pyy/homematicutils"
	"github.com/mxmCherry/movavg"
)

type PowerfoxInfo struct {
	Username string
	Password string
}

type currentResponse struct {
	Value int `json:"Watt"`
}

var powerfox_api_url = "https://backend.powerfox.energy/api/2.0/my/main/current"

func getPowerValue(client http.Client, powerfox PowerfoxInfo) int {
	req, err := http.NewRequest("GET", powerfox_api_url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.SetBasicAuth(powerfox.Username, powerfox.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var resValue currentResponse
	json.NewDecoder(resp.Body).Decode(&resValue)
	return resValue.Value
}

func main() {
	var homematic = &homematicutils.HomematicInfo{
		Hostname: os.Getenv("HOMEMATIC_HOST"),
		User:     os.Getenv("HOMEMATIC_USER"),
		Password: os.Getenv("HOMEMATIC_PASSWORD"),
	}
	var powerfox = &PowerfoxInfo{
		Username: os.Getenv("POWERFOX_USER"),
		Password: os.Getenv("POWERFOX_PASSWORD"),
	}
	var ma = movavg.ThreadSafe(movavg.NewSMA(4))
	client := &http.Client{Timeout: time.Second * 10}
	var count int = 0
	var value int = 0
	for {
		value = getPowerValue(*client, *powerfox)
		ma.Add(float64(value))
		log.Println(value, ma.Avg())
		count += 1
		if count == 4 {
			count = 0
			homematicutils.SetIntVar(*client, *homematic, "powerfox_wert", int(ma.Avg()))
		}
		time.Sleep(30 * time.Second)
	}
}
