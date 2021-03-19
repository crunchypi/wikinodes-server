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
