package config

import (
	"time"
)

// Neo4j block.
var (
	Neo4jURI = "neo4j://localhost:7687" // Default.
	Neo4jUSR = "neo4j"                  // Default.
	Neo4jPWD = "neo4j"                  // Default.
)

// Redis block.
var (
	RedisIP   = "localhost" // Default.
	RedisPort = "6379"      // Default.
	RedisPWD  = ""          // Default.
	RedisDB   = 0           // Default

	DOSGuardRefreshDelta        = time.Second * 20
	DOSGuardAllowancePerRefresh = 100
	QueryTrackExpiration        = time.Second * 20
)

// WAPI block.
var (
	// Changing IP & Port must match the ones in the
	// react app, so any changes require an update of
	// https://github.com/crunchypi/wikinodes-app/tree/master
	// After changes, a new build of that app is naturally
	// required.
	WAPIIP   = "localhost"
	WAPIPort = "1234"

	PathToReactApp = "./www/build/"

	ReadTimeout  = time.Duration(time.Second * 5)
	WriteTimeout = time.Duration(time.Second * 5)
)
