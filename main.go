package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmervine/splunk-benchmark/lib/printer"
	"github.com/jmervine/splunk-benchmark/lib/runner"
	"github.com/urfave/cli"
)

const Version = "0.0.7"

// First will be default
var outputMethods = []string{"text", "json", "jsonsummary"}

var (
	delay                      float64
	host, query                string
	runs, threads              int
	verbose, vverbose, summary bool
	printerFunc                runner.ResultPrinter
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
			Usage: fmt.Sprintf("Output method; %v", outputMethods),
			Value: outputMethods[0],
		},
		cli.StringFlag{
			Name:     "splunk-host,H",
			Usage:    "Splunk hostname; e.g. https://user:pass@splunk.example.com:8089",
			Required: true,
		},
		cli.StringFlag{
			Name:  "query,q",
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
			Name:  "summary,s",
			Usage: "Summarize; show totals",
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
		switch c.String("output") {
		case "text":
			printerFunc = printer.Text
		case "json":
			printerFunc = printer.Json
		case "jsonsummary":
			printerFunc = printer.JsonSummary
		default:
			return errors.New("Unknown output method.")
		}

		if c.Bool("verbose") {
			os.Setenv("VERBOSE", "true")
		}

		if c.Bool("very-verbose") {
			os.Setenv("VERBOSE", "true")
			os.Setenv("VERY_VERBOSE", "true")
		}

		host = c.String("splunk-host")
		query = c.String("query")
		runs = c.Int("runs")
		threads = c.Int("threads")
		delay = c.Float64("delay")
		summary = c.Bool("summary")

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(0)
	}
}

func main() {
	runner, err := runner.NewRunner(host, query, threads, runs, delay)
	if err != nil {
		panic(err)
	}

	// Handle Cleanup
	// TODO: Doesn't current stop running goroutes from runner.Start(), it
	//       probably should.
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := 0
		for range sig {
			s++
			if s == 1 {
				runner.Finalize()
				runner.Results().Print(printerFunc)
				os.Exit(130)
			} else {
				os.Exit(130)
			}
		}
	}()

	runner.Start()
	runner.Finalize()
	runner.Results().Print(printerFunc)
}
