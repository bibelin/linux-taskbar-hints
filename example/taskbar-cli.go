package main

import (
	"flag"
	"fmt"

	"github.com/bibelin/taskbar"
)

func main() {
	desktopPtr := flag.String("desktop", "", "Desktop file name for libunity backend")
	xidPtr := flag.Int("xid", 0, "X11 windows id for Xap backend")
	progressPtr := flag.Int("progress", 0, "Progress value (0-100)")
	pulsePtr := flag.Bool("pulse", false, "Pulse state")
	countPtr := flag.Int("count", 0, "Counter value")
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

	if err := tb.SetProgress(*progressPtr); err != nil {
		fmt.Println(err)
	}
	if err := tb.SetPulse(*pulsePtr); err != nil {
		fmt.Println(err)
	}
	if err := tb.SetCount(*countPtr); err != nil {
		fmt.Println(err)
	}
}
