package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	target        string
	target_rps    uint64
	init_workers  uint64
	test_duration time.Duration

	cpuprofile string

	last_rps_counter uint64
)

func init() {
	var duration_string string

	flag.StringVar(&target, "target", "http://127.0.0.1:8080/test_plain", "Target URI")
	flag.Uint64Var(&target_rps, "rps", 0, "Requests per second; 0 == no limit")
	flag.StringVar(&duration_string, "duration", "10s", "Test duration")

	flag.Uint64Var(&init_workers, "workers", 50, "Number of workers") //50 seems to give near-peak performance

	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write CPU profile to")

	flag.Parse()

	test_duration, err = time.ParseDuration(duration_string)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO -- maybe optimize it with profiler
func main() {
	const (
		chan_buffer_per_worker = 10
	)

	var (
		current_workers int = 0

		total_results TotalResults
		rps_counter   RPSCounter
		wg            sync.WaitGroup
		response_chan = make(chan RequestResult, init_workers*chan_buffer_per_worker)
	)

	//Profiling
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}

		defer pprof.StopCPUProfile()
	}

	//Prepare context
	ctx := context.Background()
	ctx = context.WithValue(ctx, CONTEXT_RPS_COUNTER_PTR, &rps_counter)
	ctx = context.WithValue(ctx, CONTEXT_WAITGROUP_PTR, &wg)
	ctx = context.WithValue(ctx, CONTEXT_RESP_CHAN, (chan<- RequestResult)(response_chan)) //convert to send only channel(workers should not be able to read data from it)
	ctx, cancel_all := context.WithCancel(ctx)

	log.Printf("Starting load on %s...\n", target)

	//create workers
	for i := uint64(0); i < init_workers; i++ {
		wg.Add(1)
		go test_request_worker(ctx)
		current_workers += 1
	}

	defer func() {
		cancel_all()

		log.Print("Finishing...\n")
		wg.Wait()

		log.Print("Done")
		print_results(&total_results)

		log.Printf("Total: %d\n", rps_counter.count)
	}()

	//main loop
	timer_stop := time.After(test_duration)
	second_ticker := time.NewTicker(time.Second)

main_loop:
	for {
		select {
		//read data and process it
		case data := <-response_chan:
			process_data(&total_results, data)

		//Each second correct workers to get desired rps(using sleep_next); zero rps counter
		case <-second_ticker.C:
			correct_workers_sleep(&rps_counter, current_workers)

		//timer to stop test
		case <-timer_stop:
			cancel_all()
			second_ticker.Stop()
			break main_loop
		}

		//Adding default case to remove finish_loop through is_stopped bool halves performance
		//So we will let finish_loop stay
	}

finish_loop:
	for {
		select {
		//read data if we have it still
		case data := <-response_chan:
			process_data(&total_results, data)

		//Stop loop if there is no data anymore to process and we've stopped our gorutines
		default:
			break finish_loop
		}
	}
}
