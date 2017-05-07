package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Data struct {
	Origin    string
	Payload   string
	TimeEntry time.Time
}

func main() {
	d := Data{
		Origin:    "fake",
		Payload:   "test",
		TimeEntry: time.Now(),
	}
	j, _ := json.MarshalIndent(d, "", "\t")
	fmt.Println(string(j))
}
