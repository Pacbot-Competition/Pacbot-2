package game

import "sync"

// Determines whether commands received in the input channel get printed
var commandLogEnable bool = false

// Mutex accompanying the above variable
var muLC sync.RWMutex

// Getter method for commandLogEnable
func getCommandLogEnable() bool {
	muLC.RLock()
	defer muLC.RUnlock()
	return commandLogEnable
}

// Setter method for commandLogEnable
func SetCommandLogEnable(en bool) {
	muLC.Lock()
	{
		commandLogEnable = en
	}
	muLC.Unlock()
}
