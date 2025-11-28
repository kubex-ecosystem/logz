// Package control defines the stages of log entry processing.
package control

import (
	ctl "github.com/kubex-ecosystem/logz/internal/module/control"
)

const (
	StepValidate ctl.JobFlag = 1 << iota
	StepPreHooks
	StepFormat
	StepPostHooks
	StepWrite
	StepDone
	StepFailed
)

const (
	StepTerminal = StepDone | StepFailed
)

type StageReg32 struct{ v ctl.FlagReg32[ctl.JobFlag] }

func NewStageReg32() *StageReg32 {
	stgR := &StageReg32{}
	return stgR
}

// func (r *StageReg32) Load() StageFlag      { return r.v.Load() }
// func (r *StageReg32) Store(val StageFlag)  { r.v.Store(val) }
// func (r *StageReg32) Set(mask StageFlag)   { r.v.Set(mask) }
// func (r *StageReg32) Clear(mask StageFlag) { r.v.Clear(mask) }
// func (r *StageReg32) SetIf(clearMask, setMask StageFlag) bool {
// 	return r.v.SetIf(clearMask, setMask)
// }
// func (r *StageReg32) Any(mask StageFlag) bool { return r.v.Any(mask) }
// func (r *StageReg32) All(mask StageFlag) bool { return r.v.All(mask) }

// type JobState = ctl.JobState

type ManagerState struct{ v ctl.FlagReg32[ctl.JobFlag] }

func (r *ManagerState) Load() ctl.JobFlag      { return r.v.Load() }
func (r *ManagerState) Store(val ctl.JobFlag)  { r.v.Store(val) }
func (r *ManagerState) Set(mask ctl.JobFlag)   { r.v.Set(mask) }
func (r *ManagerState) Clear(mask ctl.JobFlag) { r.v.Clear(mask) }
func (r *ManagerState) SetIf(clearMask, setMask ctl.JobFlag) bool {
	return r.v.SetIf(clearMask, setMask)
}
func (r *ManagerState) Any(mask ctl.JobFlag) bool { return r.v.Any(mask) }
func (r *ManagerState) All(mask ctl.JobFlag) bool { return r.v.All(mask) }

type ManagerControl struct {
	Stage ManagerState
	State ManagerState
}

func NewManagerControl() *ManagerControl {
	return &ManagerControl{
		Stage: ManagerState{v: *ctl.NewFlagReg32[ctl.JobFlag]()},
		State: ManagerState{v: *ctl.NewFlagReg32[ctl.JobFlag]()},
	}
}
