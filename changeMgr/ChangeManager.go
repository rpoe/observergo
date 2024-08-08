// Example of the application of the observer pattern.
//
// ChangeManager is the synchronization instance, it defines the interfaces:
// Notifiable and StateUpdater
//
// Ralf Poeppel 2024-08-08

package changeMgr

import (
	"fmt"

	"github.com/rpoe/observergo/state"
)

// Notifiable each ViewHandler registers with the ChangeManager
type Notifiable struct {
	Vname string           // name of the view handler
	Vchan chan state.State // channel of the view handler for receiving updates from the ChangeManager
}

// StateUpdater the interface each event needs to implement
type StateUpdater interface {
	Update(s *state.State) // update the state for this event
}

// ChangeManager synchronizes the changes of the subject and
// notifies the ViewHandlers. It returns a chan notifiable for
// receiving subscriptions of ViewHandlers and a chan stateUpdater
// for receiving subject updates. ChangeManger ends when
// chan cviewend is closed.
func ChangeManager(cviewend chan<- bool) (chan<- Notifiable, chan<- StateUpdater) {
	fmt.Println("cm start")
	creg := make(chan Notifiable)   // channel for registering a notifiable ViewHandler
	cevt := make(chan StateUpdater) // channel to receive events
	go func() {
		vs := make(map[string]Notifiable) // Map of observing views
		s := state.State{Name: "state"}   // create subject, our state
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
					v.Vchan <- s
				}
			case n := <-creg: // register a view
				fmt.Println("\ncm reg:", n)
				if n.Vchan == nil { // registering with a nil channel is a deregister
					// close chan
					v := vs[n.Vname]
					close(v.Vchan)
					delete(vs, n.Vname)
					cviewend <- true // tell main a ViewHandler has ended
				} else {
					vs[n.Vname] = n
				}
				fmt.Println("cm views:", vs)
			}
		}
	}()
	return creg, cevt
}
