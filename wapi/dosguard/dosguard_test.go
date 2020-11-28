package dosguard

import (
	"testing"
	"time"
)

func TestRegisterCheck(t *testing.T) {

	ip := "1.0.0.0"
	// # Over access.
	for i := 0; i < accessPerLimit+1; i++ {
		status := Control.RegisterCheck(ip)
		if i >= accessPerLimit && status {
			t.Error("IP accessed excessively and passed", i, status)
		}
	}
	// # Try flush
	flushDeltaSeconds = 0 // # flush is time-based
	time.Sleep(time.Second * 1)
	Control.tryFlushStale()
	if _, ok := Control.at[ip]; ok {
		t.Error("Flush failed")
	}

	// # Under access (not using loops for more control).
	flushDeltaSeconds = 1e3 // # 'disabling' this

	limitDeltaSeconds = 2
	accessPerLimit = 1
	if status := Control.RegisterCheck(ip); !status {
		t.Error("#1, should be ok")
	}
	if status := Control.RegisterCheck(ip); status {
		t.Error("#2, should be excessive")
	}
	time.Sleep(time.Second * time.Duration(limitDeltaSeconds))
	if status := Control.RegisterCheck(ip); !status {
		t.Error("#3, should be ok after reset")
	}
}
