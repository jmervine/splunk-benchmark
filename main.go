package main

import (
	"flag"

	"github.com/jmervine/splunk-benchmark/lib/search"
)

func main() {
	var (
		splunkHost  string
		splunkQuery string
		iters       int
		verbose     bool
		vverbose    bool
	)

	flag.StringVar(&splunkHost, "s", "", "Splunk hostname")
	flag.StringVar(&splunkQuery, "S", "", "Splunk query")
	flag.IntVar(&iters, "n", 10, "Number of times to perform search")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&vverbose, "vv", false, "Very verbose output")
	flag.Parse()

	runner, err := search.NewRun(splunkHost, splunkQuery, iters, verbose, vverbose)
	if err != nil {
		panic(err)
	}

	runner.Do()
	runner.PrettyPrint()
}
