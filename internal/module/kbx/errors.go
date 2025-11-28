package kbx

import "errors"

// ErrTerminal
var (
	ErrTerminal = errors.New("logz: pipeline already in terminal stage")
	ErrNotFound = errors.New("logz: entry not found")
	ErrInvalid  = errors.New("logz: invalid entry")
	ErrClosed   = errors.New("logz: entry already closed")
	ErrTimeout  = errors.New("logz: entry processing timeout")
	ErrCanceled = errors.New("logz: entry processing canceled")
	ErrAccess   = errors.New("logz: access denied")
	ErrExists   = errors.New("logz: entry already exists")
	ErrLocked   = errors.New("logz: entry is locked")
	ErrCorrupt  = errors.New("logz: entry is corrupt")
	ErrUnknown  = errors.New("logz: unknown error")
)
