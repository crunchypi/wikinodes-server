package dosguard

import (
	"sync"
	"time"
)

// This module is a simple tool for guarding against DOS attacks,
// or api over-use, in general.
//
// Package variable Control is an instance of accessControl, which
// keeps a map, containing IPs as keys and accessTime as values.
// Only intended usage of Control is by calling Control.RegisterCheck(ip)
// to check if an IP is abusive (returns false, else true). This method
// registers, updates and flushes the map automatically on each call.

// accessTime is intended to be the value of the map contained in
// accessControl.
type accessTime struct {
	// # Amount of accesses registered for an IP
	number int
	// # Time initialised, will dictate when this
	// # object will be flushed.
	time int64
}

// accessControl is intended to hold the access map and
// the method which registers and maintains said map.
type accessControl struct {
	at map[string]accessTime
	sync.Mutex
}

var (
	// Control is the main access point of this mod, makes it unnecessary to init new.
	Control accessControl = accessControl{at: make(map[string]accessTime)}
	// flushDeltaSeconds dictates how many seconds should pass before an accessControl
	// entry becomes stale and can be removed
	flushDeltaSeconds int64 = 60
	// limitDeltaSeconds dictates time in: access n/time
	limitDeltaSeconds int64 = 60
	// accessPerLimit dictates access in: access n/time
	accessPerLimit int = 120
)

// RegisterCheck manages access limitation to accessControl.
// Returns true if access from an IP is not excessive, else
// false. Method flushes stale ip entries, registers a new
// one if one does not exist, checks ip access excessiveness
// and ticks entries that are ok. Thread safe.
func (a *accessControl) RegisterCheck(ip string) bool {
	// # Async safety.
	a.Lock()
	defer a.Unlock()
	// # Remove all stale.
	a.tryFlushStale()

	v, ok := a.at[ip]
	// # New.
	if !ok {
		a.at[ip] = accessTime{number: 1, time: time.Now().Unix()}
		return true
	}
	// # Guard excessive.
	if v.time+limitDeltaSeconds > time.Now().Unix() {
		if v.number >= accessPerLimit {
			return false
		}
	}
	// # Tick and normal exit.
	a.at[ip] = accessTime{number: v.number + 1, time: v.time}
	return true
}

// tryFlushStale flushes all stale IPs. Frequency is set by
// flushDeltaSeconds (package var)
func (a *accessControl) tryFlushStale() {
	for key, val := range a.at {
		if val.time+flushDeltaSeconds <= time.Now().Unix() {
			delete(a.at, key)
		}
	}
}
