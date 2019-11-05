package printer

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/jmervine/splunk-benchmark/lib/search"
)

var o = search.Results{
	Query:   "search foo",
	Average: 1.0,
	Median:  1.0,
	Min:     1.0,
	Max:     1.0,
	Thread: []search.ResultThread{
		search.ResultThread{
			Average: 1.0,
			Median:  1.0,
			Min:     1.0,
			Max:     1.0,
		},
	},
}

func TestJsonIsJson(t *testing.T) {
	var out bytes.Buffer
	logger.SetOutput(&out)

	Json(o)

	s := new(search.Results)
	err := json.NewDecoder(&out).Decode(s)

	if err != nil {
		t.Errorf("Expected nil, got: %#v", err)
	}

	// Select a few basic values to check for equality.
	if s.Query != o.Query || s.Thread[0].Max != o.Thread[0].Max {
		t.Error("Expected 'o' to be eq to 's'")
	}
}

func TestJsonSummaryIsJson(t *testing.T) {
	var out bytes.Buffer
	logger.SetOutput(&out)

	JsonSummary(o)

	s := new(search.Results)
	err := json.NewDecoder(&out).Decode(s)

	if err != nil {
		t.Errorf("Expected nil, got: %#v", err)
	}

	// Select a few basic values to check for equality.
	if s.Query != o.Query || s.Average != o.Average || len(s.Thread) != 0 {
		t.Error("Expected 'o' to be eq to 's'")
	}
}
