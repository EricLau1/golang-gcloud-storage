package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func Console(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(bytes))
	}
}
