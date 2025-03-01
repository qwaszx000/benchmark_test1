package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
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
	count      uint64
	sleep_next time.Duration
}

var (
	request_sender   http.RoundTripper
	request_prepared *http.Request
	err              error
)

func init() {
	request_sender = &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 0,
		MaxConnsPerHost:     0,
		ForceAttemptHTTP2:   false,

		ResponseHeaderTimeout: time.Second * 3,
	}

	request_prepared, err = http.NewRequest("GET", target, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
}

// Sends test request, measures time and checks response body
func test_request() (bool, time.Duration) {
	time_before := time.Now()

	//resp, err := http.Get(target)
	resp, err := request_sender.RoundTrip(request_prepared) //Better optimized way
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

	io.Copy(io.Discard, resp.Body) //Discard data(if we leave data in body - it will cause slowdown)
	resp.Body.Close()              //Body must be closed

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
		atomic.AddUint64(&rps_counter.count, 1)
		sleep_next = time.Duration(atomic.LoadInt64((*int64)(&rps_counter.sleep_next)))

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
