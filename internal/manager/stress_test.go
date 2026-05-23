package flow_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	flow "github.com/kubex-ecosystem/logz/internal/manager"
	"github.com/kubex-ecosystem/logz/internal/manager/fsm"
)

func TestAtomicSuperpowersStress(t *testing.T) {
	const (
		numGoroutines = 1000
		opsPerRoutine = 100
	)

	f := flow.NewFlow()

	const (
		StateIdle fsm.State = 1 << iota
		StateProcessing
		StateSyncing
		StateDone
	)

	const (
		EvStart fsm.Event = iota
		EvSync
		EvFinish
	)

	f.FSM = fsm.NewFSM(StateIdle, []fsm.Transition{
		{From: StateIdle, Event: EvStart, To: StateProcessing},
		{From: StateProcessing, Event: EvSync, To: StateSyncing},
		{From: StateSyncing, Event: EvFinish, To: StateDone},
	})

	var (
		successTransitions atomic.Uint64
		failedTransitions  atomic.Uint64
	)

	wg := sync.WaitGroup{}
	wg.Add(numGoroutines)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerRoutine; j++ {
				if f.FSM.Trigger(EvStart) {
					successTransitions.Add(1)
					if f.FSM.Trigger(EvSync) {
						successTransitions.Add(1)
					}
					if f.FSM.Trigger(EvFinish) {
						successTransitions.Add(1)
					}
					f.FSM.Reset(StateIdle)
				} else {
					failedTransitions.Add(1)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("Executed %d operations in %v", numGoroutines*opsPerRoutine, duration)
	t.Logf("Success Transitions: %d", successTransitions.Load())
	t.Logf("Failed Transitions (Conflict/Invalid): %d", failedTransitions.Load())

	if successTransitions.Load() == 0 {
		t.Error("Zero successful transitions recorded - machine is stuck!")
	}
}

func TestConcurrentHookSafety(t *testing.T) {
	hooks := &flow.LogHooks{}

	var callCount atomic.Int32

	for i := 0; i < 50; i++ {
		hooks.Add(flow.HookFunc[any](func(record any) error {
			callCount.Add(1)
			return nil
		}))
	}

	wg := sync.WaitGroup{}
	numCallers := 100
	wg.Add(numCallers)

	for i := 0; i < numCallers; i++ {
		go func() {
			defer wg.Done()
			hooks.Fire(struct{ ID int }{ID: 1})
		}()
	}

	wg.Wait()

	expectedCalls := int32(50 * numCallers)
	if callCount.Load() != expectedCalls {
		t.Errorf("Hook call mismatch: expected %d, got %d", expectedCalls, callCount.Load())
	}
}
