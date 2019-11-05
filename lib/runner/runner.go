package runner

import (
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
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
	util.Verbosef("Sending %d searches at %d per %v.\n", r.runs, r.threads, r.delay)
	util.Vverbose("func=runner.Start at=start")

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

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		var i int
		for {
			if i == r.runs {
				break
			}
			<-done
			thread <- struct{}{}
			i++
		}

		fin <- true
	}()

	var i int
	for {
		if i == r.runs {
			break
		}
		select {
		case <-sig:
			return
		case <-thread:
			go func(t int) {
				r.startThread(t)
				done <- true
			}(i)
		}
		i++
	}

	<-fin
}

// TODO: Not sure if this should be limited to "threads"
func (r *Runner) Finalize() {
	util.Verbose("Collecting search results.")
	util.Vverbose("func=runner.Finalize at=start")

	// TODO: Channels might be better to avoid getting racy with
	//       updating Runner. For now I'm using sync.Mutex.
	var wg sync.WaitGroup

	defer func() {
		wg.Wait()
		util.Vverbosef("func=runner.Finalize runner:\n %#v\n", r)
		util.Vverbose("func=runner.Finalize at=finish")
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	for i := 0; i < len(r.Runs); i++ {
		wg.Add(1)
		select {
		case <-sig:
			return
		default:
			run := &r.Runs[i]
			if run.Err == nil {
				go func() {
					defer wg.Done()
					run.GetResult()
				}()
			}
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
