package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type ContextKey int

const (
	EXPECTED_BODY                      = "Hello world!"
	CONTEXT_RPS_COUNTER_PTR ContextKey = iota
	CONTEXT_WAITGROUP_PTR   ContextKey = iota
	CONTEXT_RESP_CHAN       ContextKey = iota
)

type RPSCounter struct {
	lock       sync.Mutex
	count      uint64
	sleep_next time.Duration
}

// Sends test request, measures time and checks response body
func test_request() (bool, time.Duration) {
	time_before := time.Now()
	resp, err := http.Get(target)
	took_time := time.Since(time_before)

	if err != nil {
		log.Printf("Get request error: %s\n", err)
		return false, took_time
	}

	var body_buff []byte = make([]byte, len(EXPECTED_BODY))
	n, err := resp.Body.Read(body_buff)
	if err != nil && err != io.EOF {
		log.Printf("Error while reading resp body(read bytes == %d): %s\n", n, err)
		return false, took_time
	}

	if n != len(EXPECTED_BODY) {
		log.Printf("Received %d bytes instead of expected %d\n", n, len(EXPECTED_BODY))
		return false, took_time
	}

	if string(body_buff) != EXPECTED_BODY {
		log.Printf("Received %v resp instead of expected %v\n", string(body_buff), EXPECTED_BODY)
		return false, took_time
	}

	return true, took_time

}

type RequestResult struct {
	ok   bool
	took time.Duration
}

func test_request_worker(ctx context.Context) {
	//Convert types or panic
	var rps_counter *RPSCounter = ctx.Value(CONTEXT_RPS_COUNTER_PTR).(*RPSCounter)
	var wg *sync.WaitGroup = ctx.Value(CONTEXT_WAITGROUP_PTR).(*sync.WaitGroup)
	var resp_chan = ctx.Value(CONTEXT_RESP_CHAN).(chan<- RequestResult) //send only chan

	defer wg.Done()

main_loop:
	for {
		var sleep_next time.Duration = 0
		//request
		ok, took := test_request()

		//inc rps counter
		rps_counter.lock.Lock()
		rps_counter.count += 1
		sleep_next = rps_counter.sleep_next
		rps_counter.lock.Unlock()

		//send resp data or exit on cancel
		select {
		case <-ctx.Done():
			break main_loop
		case resp_chan <- RequestResult{ok, took}:
		}

		//Sleep for duration, controlled by main gorutine
		time.Sleep(sleep_next)
	}
}
