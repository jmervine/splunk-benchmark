package runner

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/jmervine/splunk-benchmark/lib/util"
	"github.com/jmervine/splunking"
)

type Runner struct {
	client  *splunking.SplunkRequest
	runs    int
	query   string
	threads int
	delay   time.Duration
	mux     sync.Mutex

	Runs []Run
}

func NewRunner(host, query string, threads int, runs int, delay float64) (*Runner, error) {
	if threads < 1 {
		threads = 1
	}

	r := new(Runner)
	r.Runs = []Run{}

	r.query = query
	r.threads = threads
	r.runs = runs
	r.delay = time.Duration(delay) * time.Second

	client, err := splunking.InitURL(host)
	if err != nil {
		return nil, err
	}
	r.client = &client

	util.Vverbosef("runner.NewRunner:\n %#v\n", r)
	return r, nil
}

func (r *Runner) Start() {
	util.Vverbose("func=runner.Start at=start")
	fmt.Printf("Sending %d searches at %d per %v.\n", r.runs, r.threads, r.delay)

	defer func() {
		util.Vverbosef("func=runner.Start runner:\n %#v\n", r)
		util.Vverbose("func=runner.Start at=finish")
	}()

	thread := make(chan struct{}, r.threads)
	for i := 0; i < r.threads; i++ {
		thread <- struct{}{}
	}

	done := make(chan bool)
	fin := make(chan bool)

	go func() {
		for i := 0; i < r.runs; i++ {
			<-done
			thread <- struct{}{}
		}

		fin <- true
	}()

	for i := 0; i < r.runs; i++ {
		<-thread
		go func(t int) {
			r.startThread(t)
			done <- true
		}(i)
	}

	<-fin
}

// TODO: Not sure if this should be limited to "threads"
func (r *Runner) Finalize() {
	fmt.Println("Collecting search results.")
	util.Vverbose("func=runner.Finalize at=start")

	// TODO: Channels might be better to avoid getting racy with
	//       updating Runner. For now I'm using sync.Mutex.
	var wg sync.WaitGroup
	wg.Add(len(r.Runs))

	defer func() {
		wg.Wait()
		util.Vverbosef("func=runner.Finalize runner:\n %#v\n", r)
		util.Vverbose("func=runner.Finalize at=finish")
	}()

	for i := 0; i < len(r.Runs); i++ {
		run := &r.Runs[i]
		if run.Err == nil {
			go func() {
				defer wg.Done()
				run.GetResult()
			}()
		}
	}
}

func (r *Runner) Results() Results {
	util.Vverbose("func=runner.Results at=start")
	var res = Results{}

	defer func() {
		util.Vverbosef("func=runner.Results runner:\n %#v\n", r)
		util.Vverbosef("func=runner.Results results:\n %#v\n", res)
		util.Vverbose("func=runner.Results at=finish")
	}()

	runs := []Run{}
	errors := 0
	durations := []float64{}

	for _, run := range r.Runs {
		if run.Err != nil {
			errors++
			continue
		}

		runs = append(runs, run)
		durations = append(durations, run.Dur)
	}

	sort.Float64s(durations)

	res = Results{
		Runs:    runs,
		Errors:  errors,
		Average: util.Avg(durations),
		Median:  util.Med(durations),
		Min:     util.Min(durations),
		Max:     util.Max(durations),
	}

	return res
}

func (r *Runner) AppendRun(run Run) {
	// Lock Runner for writing.
	r.mux.Lock()
	defer r.mux.Unlock()

	r.Runs = append(r.Runs, run)
}

func (r *Runner) startThread(t int) {
	run := NewRun(r)
	run.Search()
	r.AppendRun(run)

	time.Sleep(r.delay)
}
