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

func Json(r search.Results) {
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
