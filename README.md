# observergo
An implementation of the observer pattern in go.

We need the observer pattern in situations where we have one data object 
being used by multiple users at the same time.

In this example we assume 3 users.
- User 1 is interested in the state only for example a log file
- User 2 and 3 are interested in the state of the data object
  and both create a series of events with different time in between.

And we have one state of the system. The state is a plain data object.
At any point in time the system has only one state and this state is 
represented by the data object.

As we have multiple users we need to make sure all users see the same
state and the state maintains the sequence of changes as the changes are
happening.

For this we have the object change manager. The change manager is responsible
for synchronising the events as they are send from the view handlers, so the state
is updated in proper sequence. And the change manager needs to inform all 
views interested in any state updates.

The system connects to the users with so called view handlers. Each view handler is
responsible for only one user.

The basic flow is: A view handler recieves an input from his associated user. The view
handler creates an event an sends it to the change manager. The change manager updates
the state and pushes the state to all view handlers which have registered for state changes.

This implementation is using no locks or semaphores. It is based purely on Go channels.
By the proper sequence of chaining of the use of the channels the system is 
operating with out any blocks and event losses. For details see the sequence diagram.

It is exspected this design is universally usable, whenever multiple users use 
concurrently a system. In terms of this design a user is any interface to a system
externally to the go runtime. For example files, or timer are considered user.

Please note the change manager does not depend on any user interface. The state
does not depend on any user intrface. The view handlers encapsulate the specific interface
to a user.

