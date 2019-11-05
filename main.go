package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/jmervine/splunk-benchmark/lib/printer"
	"github.com/jmervine/splunk-benchmark/lib/search"
)

const Version = "0.0.5"

func main() {
	var (
		delay                      float64
		host, query, output        string
		runs, threads              int
		verbose, vverbose, version bool
	)

	flag.StringVar(&output, "o", "text", "Output method: text or json")
	flag.StringVar(&host, "s", "", "Splunk hostname (https://uname:pword@host:port)")
	flag.StringVar(&query, "S", "search * | head 1", "Splunk query")
	flag.IntVar(&runs, "n", 1, "Number of search runs to perform; 0 runs until SIGINT")
	flag.IntVar(&threads, "T", 1, "Number of threads, e.g. 10 runs * 2 threads will run 20 total searches")
	flag.Float64Var(&delay, "d", 0.0, "Delay in seconds between runs (default 0.0)")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&vverbose, "vv", false, "Very verbose output")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.Parse()

	if threads < 1 {
		threads = 1
	}

	runtime.GOMAXPROCS(threads)

	if version {
		fmt.Println("splunk-benchmark version " + Version)
		os.Exit(0)
	}

	runner, err := search.NewRunner(host, query, threads, runs, delay, verbose, vverbose)
	if err != nil {
		panic(err)
	}

	pp := func() {
		if output == "json" {
			printer.Json(runner.Results())
		} else {
			r := runner.Results()
			printer.Text(r)
		}
	}

	if runs < 1 {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			pp()
			os.Exit(0)
		}()
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func(t int) {
			defer wg.Done()
			runner.Do(t)
		}(i)
	}

	wg.Wait()

	pp()
}
