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
	NumActiveGhosts  uint8
	TrustedClientIPs []string
}

// Read from the config.json file in the base directory
func GetConfig() Configuration {

	// Look in the file "config.json" in the top directory
	file, err := os.Open("../config.json")
	if err != nil {
		log.Fatalln("FATAL: could not open config.json:", err)
	}
	defer file.Close()

	// Decode the JSON arguments
	decoder := json.NewDecoder(file)
	config := Configuration{}
	if err := decoder.Decode(&config); err != nil {
		log.Fatalln("FATAL: could not parse config.json:", err)
	}

	// Validate fields that would otherwise cause runtime panics downstream
	// (e.g. a zero GameFPS divides by zero when computing the tick interval)
	if config.GameFPS <= 0 {
		log.Fatalln("FATAL: config.json GameFPS must be a positive value")
	}

	// Return the configuration when done
	return config
}
