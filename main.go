package main

import (
	"log"
	"os"

	// "github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/denisenkom/go-mssqldb"

	"hackweek41/benchmark"
	"hackweek41/prewarm"
)



func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments! Please specify either `prewarm` or `benchmark`. ex. `go run . benchmark")
	}

	option := os.Args[1]

	if option == "benchmark" {
		benchmark.BlowingUpDBConnections()
		// benchmark.BlowingUpRedisConnections()
		// benchmark.SingleRoutineQuery()
		// benchmark.MultiThreadedQuery()
		// benchmark.TestConcurrencyLimits()
	} else if (option == "prewarm") {
		prewarm.MainDriver()
	} else {
		log.Fatal("Invalid choice! Must be either `prewarm` or `benchmark`")
	}
}