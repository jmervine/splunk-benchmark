package search

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jmervine/splunk-benchmark/lib/math"
	"github.com/jmervine/splunking"
)

const (
	search = "/services/search/jobs"
	status = search + "/%s"
)

type splunkStatus struct {
	Entry []struct {
		Content struct {
			IsDone      bool    `json:"isDone"`
			RunDuration float64 `json:"runDuration"`
		} `json:"content"`
	} `json:"entry"`
}

type Runner struct {
	client   splunking.SplunkRequest
	verbose  bool
	vverbose bool
	runs     int

	Query        string
	Threads      int
	Runs         int
	Delay        time.Duration
	Sids         []string
	Results      []map[string]float64
	resultValues [][]float64
}

func NewRunner(host, query string, threads int, runs int, delay float64, verbose, vverbose bool) (*Runner, error) {
	var err error

	if threads < 1 {
		threads = 1
	}

	sr := new(Runner)
	sr.Results = make([]map[string]float64, threads)
	sr.Sids = make([]string, threads)
	sr.resultValues = make([][]float64, threads)

	for i := 0; i < threads; i++ {
		sr.Results[i] = make(map[string]float64)
	}

	sr.Query = query
	sr.Threads = threads
	sr.Runs = runs
	sr.Delay = time.Duration(delay) * time.Second
	sr.verbose = verbose
	sr.vverbose = vverbose

	sr.client, err = splunking.InitURL(host)

	if vverbose {
		fmt.Printf("%#v\n---\n\n", sr)
	}

	return sr, err
}

func (sr *Runner) Do(thread int) error {
	for {
		if sr.verbose || sr.vverbose {
			fmt.Printf("Starting run %d (thread %d)...\n", sr.runs+1, thread)
		}

		err := sr.search(thread)
		if err != nil {
			return err
		}

		for {
			done, err := sr.results(thread)
			if err != nil {
				return err
			}

			if done {
				break
			}
		}

		if sr.verbose || sr.vverbose {
			fmt.Printf("Finished run %d (thread %d).\n", sr.runs+1, thread)
		}

		if sr.Runs > 0 && sr.runs == sr.Runs {
			break
		}

		sr.runs++
	}

	sort.Float64s(sr.resultValues[thread])
	return nil
}

func (sr *Runner) search(thread int) error {
	if sr.vverbose {
		fmt.Printf("  Search Sid ")
	}

	resp, err := sr.client.Post(search, strings.NewReader("search="+sr.Query))
	if err != nil {
		return err
	}

	s := new(struct {
		Sid string `json:"sid"`
	})

	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		return err
	}

	if sr.vverbose {
		fmt.Printf("%s...\n", s.Sid)
	}

	sr.Sids[thread] = s.Sid

	return nil
}

func (sr *Runner) results(thread int) (bool, error) {
	if sr.vverbose {
		fmt.Printf("    Checking results for %s...", sr.Sids[thread])
	}

	resp, err := sr.client.Get(fmt.Sprintf(status, sr.Sids[thread]), nil)
	if err != nil {
		return false, err
	}

	s := new(splunkStatus)
	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		return false, err
	}

	done := s.Entry[0].Content.IsDone

	if sr.vverbose {
		fmt.Printf(" done: %v\n", done)

		if done {
			fmt.Printf("  Finished Sid %s in %f.\n", sr.Sids[thread], s.Entry[0].Content.RunDuration)
		}
	}

	if done {
		sr.Results[thread][sr.Sids[thread]] = s.Entry[0].Content.RunDuration
		sr.resultValues[thread] = append(sr.resultValues[thread], s.Entry[0].Content.RunDuration)
		return true, nil
	}

	return false, nil
}

func (sr *Runner) Avg(thread int) float64 {
	return math.Avg(sr.resultValues[thread])
}

func (sr *Runner) Med(thread int) float64 {
	return math.Med(sr.resultValues[thread])
}

func (sr *Runner) Min(thread int) float64 {
	return math.Min(sr.resultValues[thread])
}

func (sr *Runner) Max(thread int) float64 {
	return math.Max(sr.resultValues[thread])
}

func (sr *Runner) PrettyPrint() {
	sr.PrintBanner()
	for i := 0; i < sr.Threads; i++ {
		sr.PrintResults(i)
	}
	sr.PrintAgg()
	sr.PrintFooter()

	if sr.vverbose {
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("%#v\n---\n\n", sr)
	}
}

func (sr *Runner) PrintBanner() {
	fmt.Printf("\n %-10s | %-10s | %-10s | %-10s | %-10s | %-10s\n",
		"Thread", "Runs", "Average", "Median", "Min", "Max")
	fmt.Println("--------------------------------------------------------------------------------")
}

func (sr *Runner) PrintAgg() {
	v := []float64{}
	for _, t := range sr.resultValues {
		v = append(v, t...)
	}
	sort.Float64s(v)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("     -  aggregate  -     | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
		math.Avg(v), math.Med(v), math.Min(v), math.Max(v))
}

func (sr *Runner) PrintFooter() {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf(" Query: %.70s...\n\n", sr.Query)
}

func (sr *Runner) PrintResults(thread int) {
	fmt.Printf(" %-10d | %-10d | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
		thread, sr.runs, sr.Avg(thread), sr.Med(thread), sr.Min(thread), sr.Max(thread))
}
