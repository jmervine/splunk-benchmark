package printer

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/jmervine/splunk-benchmark/lib/runner"
)

func Text(out io.Writer, res runner.Results) {
	var output = log.New(out, "", 0)

	output.Printf("\n %-5s | %-10s | %-10s | %-10s | %s\n", "Runs", "Average", "Median", "Min", "Max")
	output.Println("------------------------------------------------------")
	output.Printf(" %-5d | %-10.3f | %-10.3f | %-10.3f | %.3f\n", len(res.Runs), res.Average, res.Median, res.Min, res.Max)
}

func JsonSummary(out io.Writer, res runner.Results) {
	type results struct {
		Average float64 `json:"average"`
		Median  float64 `json:"median"`
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
		Errors  int     `json:"errors"`
	}

	printJson(out, results{
		Average: res.Average,
		Median:  res.Median,
		Min:     res.Min,
		Max:     res.Max,
		Errors:  res.Errors,
	})
}

func Json(out io.Writer, res runner.Results) {
	printJson(out, res)
}

func printJson(out io.Writer, res interface{}) {
	var output = log.New(out, "", 0)

	data, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	var j bytes.Buffer
	err = json.Indent(&j, data, "", "  ")
	if err != nil {
		panic(err)
	}

	output.Println(j.String())
}
