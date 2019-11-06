package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/jmervine/splunk-benchmark/lib/runner"
)

func Csv(out io.Writer, res runner.Results) {
	var output = log.New(out, "", 0)
	output.Println("sid,duration,error")
	for _, run := range res.Runs {
		erro := ""
		if run.Err != nil {
			erro = fmt.Sprintf("%v", run.Err)
		}

		output.Printf("%s,%.3f,%v\n", run.Sid, run.Dur, erro)
	}
}

func Text(out io.Writer, res runner.Results) {
	var output = log.New(out, "", 0)

	output.Printf("\n %-5s | %-8s | %-8s | %-8s | %-8s | %s\n",
		"Runs", "Average", "Median", "Min", "Max", "Errors")

	output.Println("--------------------------------------------------------------")
	output.Printf(" %-5d | %-8.3f | %-8.3f | %-8.3f | %-8.3f | %d\n",
		len(res.Runs), res.Average, res.Median, res.Min, res.Max, res.Errors)
}

func JsonSummary(out io.Writer, res runner.Results) {
	type results struct {
		Average float64 `json:"average"`
		Median  float64 `json:"median"`
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
		Runs    int     `json:"runs"`
		Errors  int     `json:"errors"`
	}

	printJson(out, results{
		Average: res.Average,
		Median:  res.Median,
		Min:     res.Min,
		Max:     res.Max,
		Errors:  res.Errors,
		Runs:    len(res.Runs),
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
