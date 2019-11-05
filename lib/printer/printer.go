package printer

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/jmervine/splunk-benchmark/lib/search"
)

var logger = log.New(os.Stdout, "", 0)

func Text(r search.Results) {
	logger.Printf("\n %-10s | %-10s | %-10s | %-10s | %-10s | %-10s\n",
		"Thread", "Runs", "Average", "Median", "Min", "Max")
	logger.Println("--------------------------------------------------------------------------------")

	for i, t := range r.Thread {
		logger.Printf(" %-10d | %-10d | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
			i, len(t.Run), t.Average, t.Median, t.Min, t.Max)
	}

	logger.Println("--------------------------------------------------------------------------------")
	logger.Printf("     -  aggregate  -     | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
		r.Average, r.Median, r.Min, r.Max)

	logger.Println("--------------------------------------------------------------------------------")
	logger.Printf(" Query: %.70s...\n\n", r.Query)
}

func TextSummary(r search.Results) {
	logger.Printf("\n %-10s | %-10s | %-10s | %-10s\n", "Average", "Median", "Min", "Max")
	logger.Println("--------------------------------------------------------------------------------")
	logger.Printf(" %-10.4f | %-10.4f | %-10.4f | %-10.4f\n", r.Average, r.Median, r.Min, r.Max)
	logger.Println("--------------------------------------------------------------------------------")
	logger.Printf(" Query: %.70s...\n\n", r.Query)
}

func JsonSummary(r search.Results) {
	type results struct {
		Query   string  `json:"query"`
		Average float64 `json:"average"`
		Median  float64 `json:"median"`
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
	}

	d := results{
		Query:   r.Query,
		Average: r.Average,
		Median:  r.Median,
		Min:     r.Min,
		Max:     r.Max,
	}
	printJson(d)
}

func Json(r search.Results) {
	printJson(r)
}

func printJson(r interface{}) {
	data, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	var j bytes.Buffer
	err = json.Indent(&j, data, "", "  ")
	if err != nil {
		panic(err)
	}

	logger.Println(j.String())
}
