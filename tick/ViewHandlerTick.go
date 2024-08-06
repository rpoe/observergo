// Example of the application of the observer pattern.
//
// Event Tick and ViewHandlerTick
//
// Ralf Poeppel 2024-08-06

package tick

import (
	"fmt"
	"time"

	"github.com/rpoe/observergo/changeMgr"
	"github.com/rpoe/observergo/state"
)

const TimeMilli = "15:04:05.999999" // zero time in golang used as format

//
// ViewHandler creates Tick Events and sends them to the ChangeManager
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Event Tick
// ----------------------------------------------------------------------------

// eventTick the data send from the ViewHandlerTick to update the subject
// holding the name of the ViewHandlerTick and a timestamp.
type eventTick struct {
	vname string
	vtime time.Time
}

func (e eventTick) String() string {
	return "eventTick:" + e.vname + ", " + e.vtime.Format(TimeMilli)
}

// Update will be called by the ChangeManager for the processing of
// an event of type Tick
func (e eventTick) Update(s *state.State) {
	s.View = e.vname
	s.Tims = e.vtime
	s.Vtype = "Tick"
}

// ----------------------------------------------------------------------------
// ViewHandler Tick
// ----------------------------------------------------------------------------

// ViewHandlerTick a function running ViewHandlerTick.
// It registers for updates of the subject and updates the subject.
func ViewHandlerTick(name string, count int, delta time.Duration,
	creg chan<- changeMgr.Notifiable, cevt chan<- changeMgr.StateUpdater) {
	fmt.Println("\n" + name + " start")
	// Create a receive channel for the subject.
	cst := make(chan state.State)
	// receive and process events in a dedicated go routine, chan is closed by ChangeManager
	go func() {
		for s := range cst { // received/got notified of a subject state update
			fmt.Println(name, "show state", s)
		}
	}()
	ntf := changeMgr.Notifiable{name, cst} // new notifiable for this view
	creg <- ntf                            // register it
	ntf.Vchan = nil                        // prepare for ending
	defer func() { creg <- ntf }()         // deregister
	// create an arbitrary view input of count ticks
	ticker := time.NewTicker(delta)
	defer ticker.Stop()
	for i := 0; i < count; {
		select {
		case t := <-ticker.C:
			i++
			fmt.Println("\n"+name+" Tick", i,
				"   ", t.Format(TimeMilli))
			// send a view specific event to the event chan of the ChangeManager
			cevt <- eventTick{name, t}
		}
	}
	fmt.Println(name, "end")
}
