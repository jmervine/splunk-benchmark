package printer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jmervine/splunk-benchmark/lib/search"
)

func Text(r search.Results) {
	fmt.Printf("\n %-10s | %-10s | %-10s | %-10s | %-10s | %-10s\n",
		"Thread", "Runs", "Average", "Median", "Min", "Max")
	fmt.Println("--------------------------------------------------------------------------------")

	for i, t := range r.Thread {
		fmt.Printf(" %-10d | %-10d | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
			i, len(t.Run), t.Average, t.Median, t.Min, t.Max)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("     -  aggregate  -     | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
		r.Average, r.Median, r.Min, r.Max)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf(" Query: %.70s...\n\n", r.Query)
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

	fmt.Println(j.String())
}
