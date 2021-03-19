package redis

import (
	"testing"
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
