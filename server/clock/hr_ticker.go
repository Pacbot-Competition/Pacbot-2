package clock

import (
	"time"

	"github.com/loov/hrtime"
)

/*
A high-resolution ticker, which I designed as an alternative to time.Ticker
in case it ends up being inaccurate on Windows machines
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
		UsBetweenTicks: 1000000 * time.Microsecond / time.Duration(ticksPerSecond),
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

	for {
		select {
		case <-hrt.QuitCh:
			return
		default:
		}

		if hrtime.Since(hrt.LastTick) > hrt.UsBetweenTicks {
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
