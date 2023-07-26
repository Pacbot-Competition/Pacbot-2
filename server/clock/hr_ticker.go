package clock

import (
	"fmt"
	"sync"
	"time"

	"github.com/loov/hrtime"
)

/*
A high-resolution ticker, which I designed as an alternative to time.Ticker
in case it ends up being inaccurate on Windows machines
*/
type HighResTicker struct {
	ReadyCh        chan struct{}
	TotalTicks     int32
	quitCh         chan struct{}
	paused         bool
	usBetweenTicks time.Duration
	startTime      time.Duration
	lastTick       time.Duration
	muHRT          sync.Mutex // A mutex to protect values during public methods
}

// Create a new high-resolution ticker
func NewHighResTicker(ticksPerSecond int32) *HighResTicker {
	return &HighResTicker{
		ReadyCh:        make(chan struct{}, 10),
		TotalTicks:     0,
		quitCh:         make(chan struct{}, 1),
		paused:         false,
		usBetweenTicks: 1000000 * time.Microsecond / time.Duration(ticksPerSecond),
		startTime:      hrtime.Now(),
		lastTick:       hrtime.Now(),
	}
}

// Quit the ticker loop
func (hrt *HighResTicker) Quit() {
	<-hrt.quitCh
}

// Pause the ticker loop
func (hrt *HighResTicker) Pause() {
	hrt.muHRT.Lock()
	hrt.paused = true
	hrt.muHRT.Unlock()
}

// Unpause the ticker loop
func (hrt *HighResTicker) Play() {
	hrt.muHRT.Lock()
	hrt.paused = false
	hrt.muHRT.Unlock()
}

// Start the ticker loop
func (hrt *HighResTicker) Start() {

	defer close(hrt.ReadyCh)

	hrt.startTime = hrtime.Now()
	hrt.lastTick = hrtime.Now()

	for {

		if hrt.paused {

			fmt.Printf("Ticker paused\n")

			for hrt.paused {
				time.Sleep(100 * time.Millisecond)
			}

			fmt.Printf("Ticker unpaused\n")
		}

		select {
		case <-hrt.quitCh:
			return
		default:
		}

		if hrtime.Since(hrt.lastTick) > hrt.usBetweenTicks {
			hrt.ReadyCh <- struct{}{}
			hrt.TotalTicks++
			hrt.lastTick = hrtime.Now()
		}
	}
}

// Read the time elapsed since the high-resolution ticker started
func (hrt *HighResTicker) Lifetime() time.Duration {
	return hrtime.Since(hrt.startTime)
}
