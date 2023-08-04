package game

// Currently unused - more useful for frame rates higher than 60fps

import (
	"time"

	"github.com/loov/hrtime"
)

/*
A high-resolution ticker, which I designed as an alternative to time.Ticker
in case it ends up being inaccurate on Windows machines
*/
type HighResTicker struct {
	C              chan struct{}
	quitCh         chan struct{}
	usBetweenTicks time.Duration
	lastTick       time.Duration
}

// Create a new high-resolution ticker
func NewHighResTicker(_tickTime time.Duration) *HighResTicker {
	return &HighResTicker{
		C:              make(chan struct{}, 10),
		quitCh:         make(chan struct{}, 1),
		usBetweenTicks: _tickTime,
		lastTick:       hrtime.Now(),
	}
}

// Quit the ticker loop
func (hrt *HighResTicker) Quit() {
	<-hrt.quitCh
}

// Start the ticker loop (should be called as a go-routine)
func (hrt *HighResTicker) Start() {

	// Close the channel once we close the loop
	defer close(hrt.C)

	// Record the time of the last tick to be now
	hrt.lastTick = hrtime.Now()

	// "While" loop, once for each tick
	for {

		// Quit the loop if we get a quit signal
		select {
		case <-hrt.quitCh:
			return
		default:
		}

		// If enough time has elapsed, send an object into the ready channel
		if hrtime.Since(hrt.lastTick) > hrt.usBetweenTicks {
			hrt.C <- struct{}{}

			// Update the last tick to be now, and restart the loop
			hrt.lastTick = hrtime.Now()
		}
	}
}
