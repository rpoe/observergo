// Example of the application of the observer pattern.
//
// ViewHandlerLog
//
// Ralf Poeppel 2024-08-06

package log

import (
	"fmt"

	"github.com/rpoe/observergo/changeMgr"
	"github.com/rpoe/observergo/state"
)

// ViewHandlerLog a function running ViewHandlerLog.
// It registers for updates of the subject and logs the current state.
func ViewHandlerLog(creg chan<- changeMgr.Notifiable, cend <-chan bool) {
	name := "Log"
	fmt.Println("\n" + name + " start")
	// Create a receive channel for the subject.
	cst := make(chan state.State)
	ntf := changeMgr.Notifiable{name, cst} // new notifiable for this view
	creg <- ntf                            // register it
	ntf.Vchan = nil                        // prepare for ending
	defer func() { creg <- ntf }()         // deregister
	defer fmt.Println("   ", name, "end")
	for {
		select {
		case s := <-cst: // event of state update received
			fmt.Println("   ", name, "show state", s)
		case <-cend:
			return
		}
	}
}
