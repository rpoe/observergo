// Example of the application of the observer pattern.
// Multiple viewHandler display the state of a subject and
// modify the state of the subject concurrently.
// When one viewHandler is changing the state of the subject,
// each viewHandler receives the updated subject for updating its view.
// The viewHandler register with a ChangeManager, which will push
// the updated subject to the registered viewHandlers
// Example makes use of go routines and channels.
//
// Working of the example can be seen on stdout. You can filter with grep
// for the different packages (main, state, cm, Tick, Log).
//
// Ralf Poeppel 2024-08-08

package main

import (
	"fmt"
	"time"

	"github.com/rpoe/observergo/changeMgr"
	"github.com/rpoe/observergo/log"
	"github.com/rpoe/observergo/tick"
)

const TimeMilli = "15:04:05.999999" // zero time in golang used as format

// main managing, ChangeManager and ViewHandlers
func main() {
	fmt.Println("main start", time.Now().Format(TimeMilli))

	// we need a buffered chan with capacity of count of ViewHandlers
	countOfViewHandlers := 3
	// chan from ChangeManager to main telling end of View
	cviewend := make(chan bool, countOfViewHandlers)
	// create ChangeManager
	creg, cevt := changeMgr.ChangeManager(cviewend)

	// chan to stop log a ViewHandler not stopping by itself
	cend := make(chan bool)
	// kick of 3 view handlers
	go log.ViewHandlerLog(creg, cend)

	time.Sleep(13 * time.Millisecond) // wait to have different start time

	go tick.ViewHandlerTick("     v2", 5, 5*time.Millisecond, creg, cevt)
	countViewHandlerTick := 1 // Bookkeeping of started ViewHandler, ending by themselves

	time.Sleep(3 * time.Millisecond) // wait to have different start time

	go tick.ViewHandlerTick("     v3", 7, 3*time.Millisecond, creg, cevt)
	countViewHandlerTick++

	// if all active views are closed program should end
	for i := 1; i <= countViewHandlerTick; { // wait for ViewHandlerTick
		select {
		case <-cviewend:
			i++
			fmt.Println("main received view end")
		}
	}
	fmt.Println("main all ViewHandlerTick have ended")

	// stop ViewHandlerLog
	cend <- true
	// wait for ViewHandlerLog end
	<-cviewend
	fmt.Println("main received ViewHandlerLog end")

	fmt.Println("main end", time.Now().Format(TimeMilli))
}
