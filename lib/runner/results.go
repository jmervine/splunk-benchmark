package runner

import (
	"bytes"
	"io"
	"os"
)

type Results struct {
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Errors  int     `json:"errors"`
	Runs    []Run   `json:"runs"`
}

type ResultPrinter func(out io.Writer, res Results)

func (res Results) Print(resultPrinter ResultPrinter) {
	resultPrinter(os.Stdout, res)
}

func (res Results) String(resultPrinter ResultPrinter) string {
	buf := new(bytes.Buffer)
	resultPrinter(buf, res)
	return buf.String()
}
