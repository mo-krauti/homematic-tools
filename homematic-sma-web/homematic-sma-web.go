package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/mo-pyy/homematicutils"
	"github.com/mxmCherry/movavg"
)

var ma = movavg.ThreadSafe(movavg.NewSMA(15))
var sma_pass string

func debugReq(req *http.Request) {
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("REQUEST:\n%s", string(reqDump))
}

func debugResp(resp *http.Response) {
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RESPONSE:\n%s", string(respDump))
}

func main() {
	sma_url := os.Getenv("SMA_WEB_HOST")
  sma_pass = os.Getenv("SMA_WEB_PASSWORD")

	var homematic = &homematicutils.HomematicInfo{
		Hostname: os.Getenv("HOMEMATIC_HOST"),
		User:     os.Getenv("HOMEMATIC_USER"),
		Password: os.Getenv("HOMEMATIC_PASSWORD"),
  }

	client := &http.Client{Timeout: time.Second * 10}
	sid := authenticate(sma_url)
	var val int = 0
	var count int = 0
	for {
		val, sid = value(sma_url, sid)
		ma.Add(float64(val))
		count += 1
		if count == 5 {
			count = 0
			homematicutils.SetIntVar(*client, *homematic, "PV-Leistung", int(ma.Avg()))
		}
		time.Sleep(20 * time.Second)
	}
}

type ValueObj struct {
	Value int `json:"val"`
}

type ValueRes struct {
	Error  int `json:"err"`
	Result struct {
		V1 struct {
			V2 struct {
				List []ValueObj `json:"1"`
			} `json:"6100_40263F00"`
		} `json:"012F-730A4D39"`
	} `json:"result"`
}

func value(sma_url string, sid string) (int, string) {
	var jsonData = []byte(`{"destDev":[],"keys":["6100_40263F00"]}`)
	resp, err := http.Post(sma_url+"/dyn/getValues.json?sid="+sid, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var value_res ValueRes
	json.NewDecoder(resp.Body).Decode(&value_res)
	if value_res.Error != 0 {
		log.Println("error in value response")
		time.Sleep(10 * time.Second)
		sid := authenticate(sma_url)
		return value(sma_url, sid)
	}
	return value_res.Result.V1.V2.List[0].Value, sid
}

type AuthRes struct {
	Result map[string]string `json:"result"`
	Error  int               `json:"err"`
}

func authenticate(sma_url string) string {
	log.Println("trying to authenticate")
	postBody, _ := json.Marshal(map[string]string{
		"right": "usr",
		"pass":  sma_pass,
	})
	reqBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(sma_url+"/dyn/login.json", "application/json", reqBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var auth AuthRes
	json.NewDecoder(resp.Body).Decode(&auth)
	if auth.Error != 0 || auth.Result["sid"] == "" {
		log.Println("authentication failed")
		time.Sleep(10 * time.Second)
		return authenticate(sma_url)
	}
	log.Println("authentication successful")
	return auth.Result["sid"]
}
