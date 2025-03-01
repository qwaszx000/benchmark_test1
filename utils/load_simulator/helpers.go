package main

import (
	"log"
	"sync/atomic"
	"time"
)

// Values are zeroed by default
// https://go.dev/ref/spec#The_zero_value
type TotalResults struct {
	requests_sent uint64
	success_rate  float64

	min_latency   time.Duration
	max_latency   time.Duration
	total_latency time.Duration
}

// This function is called only in main gorutine, so we can write it without atomic operations
func process_data(total_results *TotalResults, data RequestResult) {
	//prepare to calculate new success rate
	current_succeded_requests := float64(total_results.requests_sent) * total_results.success_rate
	if data.ok {
		current_succeded_requests += 1
	}

	//requests counter
	total_results.requests_sent += 1

	//calculate new success rate
	total_results.success_rate = current_succeded_requests / float64(total_results.requests_sent)

	//latency data
	//0 means it's a first iteration and we haven't yet set min_latency
	if data.took < total_results.min_latency || total_results.min_latency == 0 {
		total_results.min_latency = data.took
	}
	if data.took > total_results.max_latency {
		total_results.max_latency = data.took
	}

	total_results.total_latency += data.took
}

func print_results(total_results *TotalResults) {
	log.Printf("Sent requests in total: %d\n", total_results.requests_sent)

	if test_duration != 0 {
		var avg_rps float64 = float64(total_results.requests_sent) / float64(test_duration/time.Second)
		log.Printf("AVG requests per second: %f\n", avg_rps)
	}

	log.Printf("Success rate: %f%%\n\n", total_results.success_rate*100)

	log.Printf("Min latency: %s\n", total_results.min_latency.String())
	log.Printf("Max latency: %s\n", total_results.max_latency.String())

	if total_results.requests_sent != 0 {
		avg_latency := total_results.total_latency / time.Duration(total_results.requests_sent)
		log.Printf("AVG latency: %s\n", avg_latency.String())
	}
}

func correct_workers_sleep(rps_counter *RPSCounter, current_workers int) {
	if target_rps == 0 {
		return
	}

	//Read current data
	sleep_duration := time.Duration(atomic.LoadInt64((*int64)(&rps_counter.sleep_next)))
	last_rps_counter = atomic.LoadUint64(&rps_counter.count)
	var new_sleep time.Duration = 0

	//report debug info
	/*log.Printf("Rps counter was: %d\n", last_rps_counter)
	log.Printf("Sleep duration: %v\n", sleep_duration)*/

	//Correct sleep duration for all workers(to reach target_rps)
	if sleep_duration == 0 && last_rps_counter < target_rps {
		//Even with 0 sleep we can't reach target_rps
		//Now the only thing we can do is nothing
		//But maybe i can optimize code or add creating new workers in the future

		//For test purpose lets create new worker
		/*wg.Add(1)
		go test_request_worker(ctx)
		current_workers += 1
		log.Printf("[EXPERIMENTAL]Can't reach target_rps, creating new worker")*/
		//I'll consider it experimental feature for now
		//But it seems to be quite unstable, so it's better not to use it
	} else if sleep_duration == 0 {
		//First correction, less accurate
		new_sleep = time.Duration(float64(time.Second) * float64(current_workers) / float64(target_rps))
	} else {
		//Assert analog, we expect count to be > 0
		if last_rps_counter == 0 {
			log.Panic("rps counter is zero")
		}

		//More accurate correction
		new_sleep = time.Duration(float64(sleep_duration) / float64(target_rps) * float64(last_rps_counter))
	}

	//zero counter and store new sleep duration
	atomic.StoreInt64((*int64)(&rps_counter.sleep_next), int64(new_sleep))
	atomic.StoreUint64(&rps_counter.count, 0)
}
