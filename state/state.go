// Example of the application of the observer pattern.
//
// The one an only State
//
// This is the one and only one object where ALL data of the system resides.
// ----------------------------------------------------------------------------
//
// Ralf Poeppel 2024-08-06

package state

import (
	"time"
)

const TimeMilli = "15:04:05.999999" // zero time in golang used as format

// state a sample subject having a state with a name and the name of the view and time
// of the last update. We can make all members public, as ChangeManager takes care on
// simultaneous access.
type State struct {
	Name  string    // name of subject
	View  string    // name of view doing last update
	Vtype string    // type of event doing last update
	Tims  time.Time // time provided with last update
}

func (s State) String() string {
	return s.Name + "," + s.View + " " + s.Vtype + ", " + s.Tims.Format(TimeMilli)
}
