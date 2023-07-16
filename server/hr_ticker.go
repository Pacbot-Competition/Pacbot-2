package main

import (
	"time"

	"github.com/loov/hrtime"
)

/*
A high-resolution ticker, which I designed as an alternative to time.Ticker
in case it is inaccurate on Windows
*/
type HighResTicker struct {
	QuitCh         chan struct{}
	ReadyCh        chan struct{}
	TotalTicks     int32
	UsBetweenTicks time.Duration
	StartTime      time.Duration
	LastTick       time.Duration
}

// Create a new high-resolution ticker
func NewHighResTicker(ticksPerSecond int32) *HighResTicker {
	return &HighResTicker{
		QuitCh:         make(chan struct{}, 1),
		ReadyCh:        make(chan struct{}, 10),
		UsBetweenTicks: time.Duration(1000000 / ticksPerSecond),
		TotalTicks:     0,
		StartTime:      hrtime.Now(),
		LastTick:       hrtime.Now(),
	}
}

// Start the ticker loop
func (hrt *HighResTicker) Start() {

	defer close(hrt.ReadyCh)

	hrt.StartTime = hrtime.Now()
	hrt.LastTick = hrtime.Now()

mainLoop:
	for {
		select {
		case <-hrt.QuitCh:
			break mainLoop
		default:
		}

		if hrtime.Since(hrt.LastTick) > hrt.UsBetweenTicks*time.Microsecond {
			hrt.ReadyCh <- struct{}{}
			hrt.TotalTicks++
			hrt.LastTick = hrtime.Now()
		}
	}
}

// Read the time elapsed since the high-resolution ticker started
func (hrt *HighResTicker) Lifetime() time.Duration {
	return hrtime.Since(hrt.StartTime)
}
