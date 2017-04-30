package main

import (
	"io"
	"log"
	"net/http"

	"github.com/crgimenes/goConfig"
	l "github.com/crgimenes/logSys"
)

type Config struct {
	Server string `json:"server" cfg:"server" cfgDefault:"localhost:8080"`
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

	/******************************
	 ** Start actuator dispatcher
	 ******************************/

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
