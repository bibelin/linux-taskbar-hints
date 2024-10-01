package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/bibelin/taskbar"
)

func main() {
	desktopPtr := flag.String("desktop", "", "Desktop file name for libunity backend")
	xidPtr := flag.Int("xid", 0, "X11 windows id for Xap backend")
	progressPtr := flag.Int("progress", 0, "Progress value (0-100)")
	pulsePtr := flag.Bool("pulse", false, "Pulse state")
	countPtr := flag.Int("count", 0, "Counter value")
	demoPtr := flag.Bool("demo", false, "Runs demonstration. Properties flags will be ignored.")
	flag.Parse()

	tb, err := taskbar.Connect(*desktopPtr, *xidPtr)
	if err != nil {
		fmt.Println(err)
		return
	}
	// That's how one may need to gracefully disconnect from taskbar.
	// Commented here to not reset properties on program end for demonstration
	// purpose.
	//
	// defer func() {
	// 	if err := tb.Disconnect(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()

	if *demoPtr {
		fmt.Println("Starting demonstration...")
		for i := 0; i < 100; i++ {
			progress := int(math.Floor(float64(i)/10.0) * 10)
			update(tb, progress, false, i)
			time.Sleep(100 * time.Millisecond)
		}
		pulse := true
		for i := 0; i < 6; i++ {
			update(tb, 0, pulse, 0)
			pulse = !pulse
			time.Sleep(1 * time.Second)
		}
		tb.Disconnect()
		fmt.Println("Finish")
	} else {
		update(tb, *progressPtr, *pulsePtr, *countPtr)
	}
}

func update(tb *taskbar.Taskbar, progress int, pulse bool, count int) {
	if err := tb.SetProgress(progress); err != nil {
		fmt.Println(err)
	}
	if err := tb.SetPulse(pulse); err != nil {
		fmt.Println(err)
	}
	if err := tb.SetCount(count); err != nil {
		fmt.Println(err)
	}
}
