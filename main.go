package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/crgimenes/goConfig"
	l "github.com/crgimenes/logSys"
)

type Config struct {
	Server string `json:"server" cfg:"server" cfgDefault:"localhost:8080"`
}

type Data struct {
	Origin    string
	Payload   string
	TimeEntry time.Time
}

var cfg = &Config{}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "...")
}

func statusHandle(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "status\n")
}

func main() {
	l.Println(l.Message, "Starting")

	c := make(chan Data)

	/******************************
	 ** Load configuration
	 ******************************/

	goConfig.PrefixEnv = "ROY"
	err := goConfig.Parse(cfg)
	if err != nil {
		l.Println(l.Error, err)
		return
	}

	/******************************
	 ** Start queues
	 ******************************/

	/******************************
	 ** Start sensor scheduler
	 ******************************/

	go func() { // fake sensor
		for {
			time.Sleep(time.Second)

			// run sensor

			cmd := exec.Command("./sensors/fake/fake")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				l.Println(l.Error, err, stderr.String())
				continue
			}
			l.Println(l.Debug, "out:", stdout.String(), "err:", stderr.String())

			// convert stdout from sensor to send to dispatcher
			var d Data
			err = json.Unmarshal(stdout.Bytes(), &d)
			if err != nil {
				l.Println(l.Error, err)
				continue
			}

			c <- d
		}
	}()

	/******************************
	 ** Start actuator dispatcher
	 ******************************/

	go func() {
		for {
			d := <-c
			fmt.Println(">", d.Origin, d.Payload, d.TimeEntry)
		}
	}()

	/******************************
	 ** Start HTTP server
	 ******************************/

	http.HandleFunc("/", mainHandle)
	http.HandleFunc("/status", statusHandle)

	l.Println(l.Message, "Listen on http://", cfg.Server)
	err = http.ListenAndServe(cfg.Server, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
