package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var (
	ctx = context.Background()
	// How long to keep ip:queryid(wiki) alive.
	// This is used to keep track of which Wikipedia
	// article (neo4j db) IDs are used by IPs for
	// the purpose of recommendations.
	queryIDExpiration = time.Second * 20
	// Namespace or ip:queryid(wiki) keys.
	namespaceWikiID = "wikiID"
	// These two are used to prevent service
	// spam. An IP is allowed to make x amount
	// of requests per t amount of time, where
	// x = dosguardAllowance and
	// t = dosguardExpiration
	dosguardExpiration = time.Second * 10
	dosguardAllowance  = 100
	namespaceDosguard  = "dg"
)

type RedisManager struct {
	c *redis.Client
}

// New sets up- and returns a RedisManager with a Redis client
func New(ip, port, pwd string, db int) *RedisManager {
	return &RedisManager{
		c: redis.NewClient(&redis.Options{
			Addr:     ip + ":" + port,
			Password: pwd,
			DB:       db,
		})}
}

// SetLastQueryID tries to set a query id for an ip. Intenden
// to be used for keeping track of which Wikipedia Articles a
// front-end client searches for, for the purpose of article
// recommendation.
func (r *RedisManager) SetLastQueryID(ip string, id int64) bool {
	err := r.c.Set(ctx, namespaceWikiID+ip, id, queryIDExpiration).Err()
	if err != nil {
		return false
	}
	return true
}

// LastQueryID is the counterpart of SetLastQueryID, it simply
// tries to retrieve a Wikipedia Article for a given IP.
func (r *RedisManager) LastQueryID(ip string) (int64, bool) {
	v, err := r.c.Get(ctx, namespaceWikiID+ip).Result()
	if err != nil {
		return 0, false
	}
	res, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}
	return res, true
}

// Used to prevent service spam. Calling this method will
// increment the counter for an IP and check if it has
// exceeded an allowance over a time period (see pkg vars
// dosguardAllowance & dosguardExpiration). If True is
// returned, then the IP is good for more requests.
func (r *RedisManager) CheckRegDOSIP(ip string) (bool, error) {
	v, err := r.c.Get(ctx, namespaceDosguard+ip).Result()
	// # If key doesn't exist, create it with fresh allowance.
	if err != nil {
		r.c.Set(ctx, namespaceDosguard+ip, 1, dosguardExpiration)
		return true, nil
	}
	// # Shouldn't be a problem if this method is self-contained
	// # since an int is guaranteed(?), given the block above.
	// # But still...
	count, err := strconv.Atoi(v)
	if err != nil {
		return false, err
	}
	// # Allowance exceeded.
	if count+1 > dosguardAllowance {
		return false, nil
	}
	// # Ok: Increment and allow.
	err = r.c.Incr(ctx, namespaceDosguard+ip).Err()
	return true, err
}
