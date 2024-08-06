// Example of the application of the observer pattern.
// Multiple viewHandler display the state of a subject and
// modify the state of the subject concurrently.
// When one viewHandler is changing the state of the subject,
// each viewHandler receives the updated subject for updating its view.
// The viewHandler register with a ChangeManager, which will push
// the updated subject to the registered viewHandlers
// Example makes use of go routines and channels
//
// Ralf Poeppel 2024-08-06 (Never forget Hiroshima)

package main

import (
	"fmt"
	"time"
)

const TimeMilli = "15:04:05.999999" // zero time in golang used as format

// ============================================================================
// The one an only State
//
// This is the one and only one object where ALL data of the system resides.
// ----------------------------------------------------------------------------

// state a sample subject having a state with a name and the name of the view and time
// of the last update
type state struct {
	name  string    // name of subject
	view  string    // name of view doing last update
	vtype string    // type of event doing last update
	tims  time.Time // time provided with last update
}

func (s state) String() string {
	return s.name + "," + s.view + " " + s.vtype + ", " + s.tims.Format(TimeMilli)
}

// ============================================================================
// notifiable and stateUpdate define the interface used by the ChangeManager
// ----------------------------------------------------------------------------
// notifiable each ViewHandler registers with the ChangeManager
type notifiable struct {
	vname string     // name of the view handler
	vchan chan state // channel of the view handler for receiving updates from the ChangeManager
}

// ----------------------------------------------------------------------------
// State Updater Interface
// ----------------------------------------------------------------------------

// stateUpdater the interface each event needs to implement
type stateUpdater interface {
	Update(s *state) // update the state for this event
}

// ============================================================================
// ViewHandler, a ViewHandler may define a specific event for input data
// from its view.
// ----------------------------------------------------------------------------

// ViewHandlerLog a function running ViewHandlerLog.
// It registers for updates of the subject and logs the current state.
func ViewHandlerLog(creg chan<- notifiable, cend <-chan bool) {
	name := "Log"
	fmt.Println("\n" + name + " start")
	// Create a receive channel for the subject.
	cst := make(chan state)
	ntf := notifiable{name, cst}   // new notifiable for this view
	creg <- ntf                    // register it
	ntf.vchan = nil                // prepare for ending
	defer func() { creg <- ntf }() // deregister
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

// ============================================================================
// 
// Event Tick and ViewHandlerTick
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
func (e eventTick) Update(s *state) {
	s.view = e.vname
	s.tims = e.vtime
	s.vtype = "Tick"
}

// ----------------------------------------------------------------------------
// ViewHandler Tick
// ----------------------------------------------------------------------------

// ViewHandlerTick a function running ViewHandlerTick.
// It registers for updates of the subject and updates the subject.
func ViewHandlerTick(name string, count int, delta time.Duration,
	creg chan<- notifiable, cevt chan<- stateUpdater) {
	fmt.Println("\n" + name + " start")
	// Create a receive channel for the subject.
	cst := make(chan state)
	// receive and process events in a dedicated go routine, chan is closed by ChangeManager
	go func() {
		for s := range cst { // received/got notified of a subject state update
			fmt.Println(name, "show state", s)
		}
	}()
	ntf := notifiable{name, cst}   // new notifiable for this view
	creg <- ntf                    // register it
	ntf.vchan = nil                // prepare for ending
	defer func() { creg <- ntf }() // deregister
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

// ----------------------------------------------------------------------------
// The one and only Change Manager
// ----------------------------------------------------------------------------

// ChangeManager synchronizes the changes of the subject and
// notifies the ViewHandlers. It returns a chan notifiable for
// receiving subscriptions of ViewHandlers and a chan stateUpdater
// for receiving subject updates. ChangeManger ends when
// chan cviewend is closed.
func ChangeManager(cviewend chan<- bool) (chan<- notifiable, chan<- stateUpdater) {
	fmt.Println("cm start")
	creg := make(chan notifiable)   // channel for registering a notifiable ViewHandler
	cevt := make(chan stateUpdater) // channel to receive events
	go func() {
		vs := make(map[string]notifiable) // Map of observing views
		s := state{name: "state"}         // create subject, our state
		for {
			fmt.Println("cm waiting with state:", s)
			select {
			case e := <-cevt: // event of state update received
				fmt.Println("cm rceiv", e)
				// update state of subject
				e.Update(&s)
				fmt.Println("\nnew state of subject:", s)
				// notify observers/views
				for _, v := range vs {
					v.vchan <- s
				}
			case n := <-creg: // register a view
				fmt.Println("\ncm reg:", n)
				if n.vchan == nil { // registering with a nil channel is a deregister
					// close chan
					v := vs[n.vname]
					close(v.vchan)
					delete(vs, n.vname)
					cviewend <- true // tell main a ViewHandler has ended
				} else {
					vs[n.vname] = n
				}
				fmt.Println("cm views:", vs)
			}
		}
	}()
	return creg, cevt
}

// ----------------------------------------------------------------------------

// main managing, ChangeManager and ViewHandlers
func main() {
	countOfViewHandlers := 3
	// we need a buffered chan with capacity of count of ViewHandlers
	cviewend := make(chan bool, countOfViewHandlers) // chan from ChangeManager to main telling end of View
	cend := make(chan bool)                          // chan to stop log a ViewHandler not stopping by itself
	// create ChangeManager
	creg, cevt := ChangeManager(cviewend)

	// kick of 3 view handlers
	go ViewHandlerLog(creg, cend)
	time.Sleep(13 * time.Millisecond) // wait to have different start time
	go ViewHandlerTick("     v2", 5, 5*time.Millisecond, creg, cevt)
    countViewHandlerTick := 1  // Bookkeeping of started ViewHandler, ending by themselves
	time.Sleep(3 * time.Millisecond) // wait to have different start time
	go ViewHandlerTick("     v3", 7, 3*time.Millisecond, creg, cevt)
    countViewHandlerTick++

	// if all active views are closed program should end
	for i := 1; i <= countViewHandlerTick; { // wait for ViewHandlerTick
		select {
		case <-cviewend:
			i++
			fmt.Println("Main received view end")
		}
	}
	fmt.Println("All ViewHandlerTick have ended")

	// stop ViewHandlerLog
	cend <- true
	// wait for ViewHandlerLog end
	<-cviewend
	fmt.Println("Main received ViewHandlerLog end")

	fmt.Println("main end", time.Now().Format(TimeMilli))
}
