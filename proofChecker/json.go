package main

import (
	"encoding/json"
	"log"
)

func marshalJson() string {

	output, err := json.Marshal(dsp)

	if err != nil {

		log.Fatal("json falied")
	}

	return string(output)
}

func unmarshalJson(data string) {

	dsp.clear()
	err := json.Unmarshal([]byte(data), dsp)

	if err != nil {
		log.Fatal("json falied")
	}
}
