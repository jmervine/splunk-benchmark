package main

import (
	"flag"

	"github.com/jmervine/splunk-benchmark/lib/search"
)

func main() {
	var (
		splunkHost  string
		splunkQuery string
		runs        int
		delay       float64
		verbose     bool
		vverbose    bool
	)

	flag.StringVar(&splunkHost, "s", "", "Splunk hostname (https://uname:pword@host:port)")
	flag.StringVar(&splunkQuery, "S", "", "Splunk query")
	flag.IntVar(&runs, "n", 10, "Number of search runs to perform")
	flag.Float64Var(&delay, "d", 0.0, "Delay in seconds between runs")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&vverbose, "vv", false, "Very verbose output")
	flag.Parse()

	runner, err := search.NewRun(splunkHost, splunkQuery, runs, delay, verbose, vverbose)
	if err != nil {
		panic(err)
	}

	runner.Do()
	runner.PrettyPrint()
}
