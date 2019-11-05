package printer

import (
	"bytes"
	"io"
	"testing"

	"github.com/jmervine/splunk-benchmark/lib/runner"
)

var o = runner.Results{
	Average: 1.2,
	Median:  1.2,
	Min:     1.1,
	Max:     1.3,
	Errors:  1,
	Runs: []runner.Run{
		runner.Run{Sid: "sid-foo-1", Dur: 1.1},
		runner.Run{Sid: "sid-foo-2", Dur: 1.2},
		runner.Run{Sid: "sid-foo-3", Dur: 1.3},
	},
}

func TestText(t *testing.T) {
	var out bytes.Buffer

	Text(io.Writer(&out), o)

	expected := `
 Runs  | Average    | Median     | Min        | Max
------------------------------------------------------
 3     | 1.200      | 1.200      | 1.100      | 1.300
`
	str := out.String()
	if str != expected {
		t.Errorf("Expected json body was not provided.\nExpected:\n`%s`\nActual:\n`%s`", expected, str)
	}
}

func TestJsonIsJson(t *testing.T) {
	var out bytes.Buffer

	Json(io.Writer(&out), o)

	expected := `{
  "average": 1.2,
  "median": 1.2,
  "min": 1.1,
  "max": 1.3,
  "errors": 1,
  "runs": [
    {
      "sid": "sid-foo-1",
      "duration": 1.1
    },
    {
      "sid": "sid-foo-2",
      "duration": 1.2
    },
    {
      "sid": "sid-foo-3",
      "duration": 1.3
    }
  ]
}
`

	str := out.String()
	if str != expected {
		t.Errorf("Expected json body was not provided.\nExpected:\n`%s`\nActual:\n`%s`", expected, str)
	}
}

func TestJsonSummaryIsJson(t *testing.T) {
	var out bytes.Buffer

	JsonSummary(io.Writer(&out), o)

	expected := `{
  "average": 1.2,
  "median": 1.2,
  "min": 1.1,
  "max": 1.3,
  "runs": 3,
  "errors": 1
}
`

	str := out.String()
	if str != expected {
		t.Errorf("Expected json body was not provided.\nExpected:\n`%s`\nActual:\n`%s`", expected, str)
	}
}
