package main

import (
	"context"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	target          string
	target_rps      uint64
	result_filename string

	init_workers uint64

	cpuprofile string
	memprofile string

	last_rps_counter uint64
)

func init() {
	flag.StringVar(&target, "target", "http://127.0.0.1:8080/test_plain", "Target URI")
	flag.Uint64Var(&target_rps, "rps", 100, "Requests per second")
	flag.StringVar(&result_filename, "file", "", "File to append csv results to")

	flag.Uint64Var(&init_workers, "workers", 50, "Number of workers") //50 seems to give near-peak performance

	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write CPU profile to")
	flag.StringVar(&memprofile, "memprofile", "", "Write Memory profile to")

	flag.Parse()
}

func write_resp(csv_writer *csv.Writer, data *RequestResult) {
	if csv_writer == nil {
		return
	}
	err := csv_writer.Write(
		[]string{
			//strconv.FormatUint(target_rps, 10),
			strconv.FormatUint(last_rps_counter, 10),
			strconv.FormatBool(data.ok),
			data.took.String(),
		},
	)

	if err != nil {
		log.Fatalf("CSV writer error: %s\n", err)
	}
}

// TODO -- maybe optimize it with profiler
func main() {
	const (
		chan_buffer_per_worker = 10
	)

	var (
		current_workers int = 0

		rps_counter   RPSCounter
		wg            sync.WaitGroup
		response_chan = make(chan RequestResult, init_workers*chan_buffer_per_worker)

		result_writer *csv.Writer = nil
	)

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

	//Open file and create csv writer
	if result_filename != "" {
		result_file, err := os.OpenFile(result_filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Can't open file %s: %s\n", result_filename, err)
		}

		result_writer = csv.NewWriter(result_file)

		defer func() {
			result_writer.Flush()
			result_file.Close()
		}()
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
	}()

	//main loop
	timer_stop := time.After(time.Second * 10)
	second_ticker := time.NewTicker(time.Second)
main_loop:
	for {
		select {
		//read data and write it to csv file
		case data := <-response_chan:
			//log.Printf("%v %v\n", data.ok, data.took)
			write_resp(result_writer, &data)

		//Each second correct workers to get desired rps(using sleep_next); zero rps counter
		case <-second_ticker.C:

			//Read current data
			sleep_duration := time.Duration(atomic.LoadInt64((*int64)(&rps_counter.sleep_next)))
			last_rps_counter = atomic.LoadUint64(&rps_counter.count)
			var new_sleep time.Duration = 0

			//report debug info
			log.Printf("Rps counter was: %d\n", last_rps_counter)
			log.Printf("Sleep duration: %v\n", sleep_duration)

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
					log.Panic("rsp counter is zero")
				}

				//More accurate correction
				new_sleep = time.Duration(float64(sleep_duration) / float64(target_rps) * float64(last_rps_counter))
			}

			//zero counter and store new sleep duration
			atomic.StoreInt64((*int64)(&rps_counter.sleep_next), int64(new_sleep))
			atomic.StoreUint64(&rps_counter.count, 0)

		//10 seconds timer to stop load
		case <-timer_stop:
			cancel_all()
			second_ticker.Stop()
			break main_loop
		}
	}

	//Process all buffered data
finish_loop:
	for {
		select {
		//read data if we have it still
		case data := <-response_chan:
			//log.Printf("%v %v\n", data.ok, data.took)

			write_resp(result_writer, &data)
		default:
			break finish_loop
		}
	}
}
