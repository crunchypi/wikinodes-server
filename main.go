package main

import (
	"fmt"
	"log"
	"wikinodes-server/config"
	"wikinodes-server/db/neo4j"
	"wikinodes-server/db/redis"
	"wikinodes-server/wapi"
)

func main() {

	r := redis.New(config.RedisIP, config.RedisPort, config.RedisPWD, config.RedisDB)
	n, err := neo4j.New(config.Neo4jURI, config.Neo4jUSR, config.Neo4jPWD)
	if err != nil {
		msg := fmt.Sprint("neo4j setup err:", err)
		log.Fatal(msg)
	}

	if err = wapi.Start(n, r); err != nil {
		log.Fatal(err)
	}

}
