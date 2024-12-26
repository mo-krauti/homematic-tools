package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
  "fmt"
  "os"

	"github.com/mo-pyy/homematicutils"
)


type apiResponse struct {
	Value int `json:"state"`
}

var api_url = "https://api.stromgedacht.de/v1/now?zip="

func getApiValue(client http.Client) int {
	req, err := http.NewRequest("GET", api_url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var resValue apiResponse
	json.NewDecoder(resp.Body).Decode(&resValue)
	return resValue.Value
}

func main() {
	var homematic = &homematicutils.HomematicInfo{
		Hostname: os.Getenv("HOMEMATIC_HOST"),
		User:     os.Getenv("HOMEMATIC_USER"),
		Password: os.Getenv("HOMEMATIC_PASSWORD"),
	}
	client := &http.Client{Timeout: time.Second * 10}
  api_url += os.Getenv("STOMGEDACHT_PLZ")
	var value int = 0
	for {
		value = getApiValue(*client)
    homematicutils.SetIntVar(*client, *homematic, "stromgedacht_wert", value)
    fmt.Println(value)
		time.Sleep(60 * time.Second)
	}
}
