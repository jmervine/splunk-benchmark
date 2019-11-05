package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/jmervine/splunk-benchmark/lib/printer"
	"github.com/jmervine/splunk-benchmark/lib/search"
	"github.com/urfave/cli"
)

const Version = "0.0.6"

var (
	delay                      float64
	host, query, output        string
	runs, threads              int
	verbose, vverbose, version bool
)

func init() {
	oldHelpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		oldHelpPrinter(w, templ, data)
		os.Exit(0)
	}

	app := cli.NewApp()
	app.Name = "splunk-benchmark"
	app.Usage = "search load generation"
	app.Description = "Simple search load generation tool for Splunk"
	app.UsageText = "splunk-benchmark [args...]"
	app.Version = Version
	app.Author = "Joshua Mervine"
	app.Email = "joshua@mervine.net"
	app.EnableBashCompletion = true
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output,o",
			Usage: "Output method; text, json",
			Value: "text",
		},
		cli.StringFlag{
			Name:     "splunk-host,s",
			Usage:    "Splunk hostname; e.g. https://user:pass@splunk.example.com:8089",
			Required: true,
		},
		cli.StringFlag{
			Name:  "query,S",
			Usage: "Splunk search query",
			Value: "search * | head 1",
		},
		cli.IntFlag{
			Name:  "runs,r",
			Usage: "Number of search runs to perform; 0 runs until SIGINT",
			Value: 1,
		},
		cli.IntFlag{
			Name:  "threads,T",
			Usage: "Number of threads, e.g. 10 runs * 2 threads will run 20 total searches",
			Value: 1,
		},
		cli.Float64Flag{
			Name:  "delay,d",
			Usage: "Delay in seconds between runs",
			Value: 0.0,
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Verbose output",
		},
		cli.BoolFlag{
			Name:  "very-verbose",
			Usage: "Very verbose output",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("version") {
			fmt.Printf("splunk-benchmark version %s\n", Version)
			os.Exit(0)
		}

		output = c.String("output")
		host = c.String("splunk-host")
		query = c.String("query")
		runs = c.Int("runs")
		threads = c.Int("threads")
		delay = c.Float64("delay")
		verbose = c.Bool("verbose")
		vverbose = c.Bool("very-verbose")

		if threads < 1 {
			threads = 1
		}
		runtime.GOMAXPROCS(threads)

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(0)
	}
}

func main() {
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
