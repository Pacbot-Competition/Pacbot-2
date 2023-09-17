package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ServerIP         string
	TcpPort          int
	WebSocketPort    int
	OneClientPerIP   bool
	GameFPS          int32
	TrustedClientIPs []string
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
		log.Println("JSON read error:", err)
	}

	// Return the configuration when done
	return config
}
