package main

import (
	"fmt"
	"log"
	"os"
	"wikinodes-server/db/neo4j"
	"wikinodes-server/db/redis"
	"wikinodes-server/wapi"
)

func main() {

	r := redis.New("localhost", "6379", "", 0)
	if len(os.Args) != 4 {
		msg := "API uses Neo4j, requires 3 args: <uri> <usr> <pwd>"
		log.Fatal(msg)
	}
	uri, usr, pwd := os.Args[1], os.Args[2], os.Args[3]
	n, err := neo4j.New(uri, usr, pwd)
	if err != nil {
		msg := fmt.Sprint("neo4j setup err:", err)
		log.Fatal(msg)
	}

	if err = wapi.Start(n, r); err != nil {
		log.Fatal(err)
	}

}
