package main

import (
	"encoding/json"
	"log"
)

func (d *console) marshalJson() string {

	output, err := json.Marshal(d)

	if err != nil {

		log.Fatal("json falied")
	}

	return string(output)
}

func (d *console) unmarshalJson(data string) {

	d.clear()
	err := json.Unmarshal([]byte(data), d)

	if err != nil {
		log.Fatal("json falied")
	}
}
