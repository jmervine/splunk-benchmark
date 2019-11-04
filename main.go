package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/jmervine/splunk-benchmark/lib/search"
)

const Version = "0.0.1"

func main() {
	var (
		delay                      float64
		host, query                string
		runs, threads              int
		verbose, vverbose, version bool
	)

	flag.StringVar(&host, "s", "", "Splunk hostname (https://uname:pword@host:port)")
	flag.StringVar(&query, "S", "", "Splunk query")
	flag.IntVar(&runs, "n", 10, "Number of search runs to perform")
	flag.IntVar(&threads, "T", 1, "Number of threads, e.g. 10 runs * 2 threads will run 20 total searches")
	flag.Float64Var(&delay, "d", 0.0, "Delay in seconds between runs")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&vverbose, "vv", false, "Very verbose output")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.Parse()

	if version {
		fmt.Println("splunk-benchmark version " + Version)
		os.Exit(0)
	}

	runner, err := search.NewRunner(host, query, threads, runs, delay, verbose, vverbose)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(runner.Threads)

	for i := 0; i < runner.Threads; i++ {
		go func(t int) {
			defer wg.Done()
			runner.Do(t)
		}(i)
	}

	wg.Wait()

	runner.PrettyPrint()
}