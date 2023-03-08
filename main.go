package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default database
	})

	errs := 0
	max := 1000000
	steps := max / 10000

	for i := 0; i < steps; i++ {
		path := fmt.Sprintf("ns%d", i)
		errs += writeToRedis(client, path, max/steps)

		time.Sleep(200 * time.Millisecond)
		milisecondsPerRequest, errsOnRead := measureAveragePerformance(client, max, steps, i)
		fmt.Printf("Average response time: %f microseconds, errors: %d\n", milisecondsPerRequest, errsOnRead)
	}

	fmt.Printf("Total errors adding %d keys was: %d\n", max, errs)
}
