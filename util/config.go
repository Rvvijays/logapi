package util

import (
	"encoding/json"
	"fmt"
	"os"
)

var Config config

// read configs from config.json to store it global Config
func InitConfig() {

	file, err := os.Open("config.json")
	if err != nil {
		// fmt.Println("error while loading config file")
		panic("error while loading config file")
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Config)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("config:", Config)

}
