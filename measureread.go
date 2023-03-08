package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func measureAveragePerformance(client *redis.Client, max int, steps int, currentMax int) (float64, int32) {
	var wg sync.WaitGroup
	totalMeasurements := max / (steps * 100)
	wg.Add(totalMeasurements)

	mx := sync.Mutex{}

	totalTime := 0 * time.Second
	responseTimeChan := make(chan time.Duration, totalMeasurements)

	// measure response time
	go func() {
		for {
			select {
			case timetaken := <-responseTimeChan:
				mx.Lock()
				totalTime += timetaken
				mx.Unlock()
				wg.Done()
			}
		}
	}()

	errcounter := atomic.Int32{}

	// execute requests
	for i := 0; i < totalMeasurements; i++ {
		go func() {

			start := time.Now()
			path := fmt.Sprintf("ns%d/secret-%d", rand.Intn(currentMax+1), rand.Intn(max/steps))

			resp, err := client.Get(path).Result()
			if err != nil {
				errcounter.Add(1)
				wg.Done()
				return
			}
			_ = resp

			elapsed := time.Since(start)
			responseTimeChan <- elapsed
		}()
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
	}
	wg.Wait()
	microSecondsPerRequest := float64(int(totalTime.Microseconds())) / float64(time.Microsecond)
	errorcount := errcounter.Load()

	return microSecondsPerRequest, errorcount
}
