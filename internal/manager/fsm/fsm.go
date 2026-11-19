package fsm

import (
	"github.com/kubex-ecosystem/logz/internal/manager/control"
)

// State represents an atomic state based on bitflags.
type State uint32

// Event represents an event that can trigger a state transition.
type Event uint32

// Transition represents an atomic transition between states.
type Transition struct {
	From  State
	Event Event
	To    State
}

// FSM represents an atomic finite state machine.
type FSM struct {
	current control.FlagReg32[State]
	table   map[State]map[Event]State
}

// NewFSM creates a new atomic FSM with the given initial state and transitions.
func NewFSM(initial State, transitions []Transition) *FSM {
	fsm := &FSM{
		table: make(map[State]map[Event]State),
	}
	fsm.current.Store(initial)

	for _, t := range transitions {
		if fsm.table[t.From] == nil {
			fsm.table[t.From] = make(map[Event]State)
		}
		fsm.table[t.From][t.Event] = t.To
	}

	return fsm
}

// Current returns the current state of the FSM atômico.
func (f *FSM) Current() State {
	return f.current.Load()
}

// Trigger transitions the FSM to a new state atomically using CAS.
func (f *FSM) Trigger(event Event) bool {
	for {
		curr := f.current.Load()
		if next, ok := f.table[curr][event]; ok {
			if f.current.CompareAndSwap(curr, next) {
				return true
			}
			// If CAS fails, someone else changed the state, retry the logic from the new state.
			continue
		}
		return false
	}
}

// Can checks if the given event can trigger a transition from the current state.
func (f *FSM) Can(event Event) bool {
	_, ok := f.table[f.current.Load()][event]
	return ok
}

// Reset forces the FSM to the given state.
func (f *FSM) Reset(state State) {
	f.current.Store(state)
}
