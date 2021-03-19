package redis

import (
	"testing"
	"time"
)

var (
	ip   = "localhost"
	port = "6379"
	pwd  = ""
	db   = 0
	r    *RedisManager
)

func init() {
	r = New(ip, port, pwd, db)
}

func TestSetGetLastQueryID(t *testing.T) {
	ip := "0.0.0.0"
	id := int64(1)

	// # Prep.
	r.c.Del(ctx, ip)

	_, ok := r.LastQueryID(ip)
	if ok {
		t.Fatal("unexpected query success")
	}

	if ok := r.SetLastQueryID(ip, id); !ok {
		t.Fatal("failed while setting k:v")
	}

	v, ok := r.LastQueryID(ip)
	if !ok || v != id {
		t.Fatalf("unexpected query fail: %v, %v", v, ok)
	}
}

func TestCheckRegDOSIP(t *testing.T) {
	// # Backup pkg vars so it's safe to reduce
	// # allowance and expiration (so test is quicker).
	dguardExpBackup := dosguardExpiration
	dguardAllowBackup := dosguardAllowance

	ip := "0.0.0.0"
	expire := time.Second * 3
	allow := 2

	dosguardExpiration = expire
	dosguardAllowance = allow

	// # Use up allowance.
	for i := 0; i < allow; i++ {
		if ok, err := r.CheckRegDOSIP(ip); !ok || err != nil {
			t.Fatalf("checkreg step 1 (iter %v) fail: %v, %v", i, ok, err)
		}
	}
	// # Exceed allowance
	if ok, err := r.CheckRegDOSIP(ip); ok || err != nil {
		t.Fatalf("checkreg step 2 fail: %v, %v", ok, err)
	}
	// # Wait until ip expires.
	time.Sleep(expire + 1)
	if ok, err := r.CheckRegDOSIP(ip); !ok || err != nil {
		t.Fatalf("checkreg step 3 fail: %v, %v", ok, err)
	}
	// # Cleanup.
	dosguardExpiration = dguardExpBackup
	dosguardAllowance = dguardAllowBackup
}
