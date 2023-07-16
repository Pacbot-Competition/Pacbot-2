package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	ServerIP        string
	TcpPort         int
	WebSocketPort   int
	OneBrowserPerIP bool
}

// Read from the config.json file in the base directory
func GetConfig() Configuration {

	// Look in the file "config.json" in the top directory
	file, _ := os.Open("../config.json")
	defer file.Close()

	// Decode the JSON arguments
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("JSON read error:", err)
	}

	// Return the configuration when done
	return config
}
