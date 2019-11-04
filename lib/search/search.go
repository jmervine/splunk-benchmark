package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

type ResultThread struct {
	Average float64   `json:"average"`
	Median  float64   `json:"median"`
	Min     float64   `json:"min"`
	Max     float64   `json:"max"`
	Run     []float64 `json:"run"`
}

type Results struct {
	Query   string         `json:"query"`
	Average float64        `json:"average"`
	Median  float64        `json:"median"`
	Min     float64        `json:"min"`
	Max     float64        `json:"max"`
	Thread  []ResultThread `json:"thread"`
}

type Runner struct {
	client       splunking.SplunkRequest
	verbose      bool
	vverbose     bool
	ran, runs    int
	query        string
	threads      int
	delay        time.Duration
	sids         []string
	resultValues [][]float64
	//Results      Results
}

func NewRunner(host, query string, threads int, runs int, delay float64, verbose, vverbose bool) (*Runner, error) {
	var err error

	sr := new(Runner)
	sr.sids = make([]string, threads)
	sr.resultValues = make([][]float64, threads)

	sr.query = query
	sr.threads = threads
	sr.runs = runs
	sr.delay = time.Duration(delay) * time.Second
	sr.verbose = verbose
	sr.vverbose = vverbose

	sr.client, err = splunking.InitURL(host)

	if vverbose {
		log.Printf("%#v\n---\n\n", sr)
	}

	return sr, err
}

func (sr *Runner) Do(thread int) error {
	for {
		if sr.verbose || sr.vverbose {
			log.Printf("Starting run %d (thread %d)...\n", sr.ran+1, thread)
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
			log.Printf("Finished run %d (thread %d).\n", sr.ran+1, thread)
		}

		if sr.runs > 0 && sr.ran == sr.runs {
			break
		}

		sr.ran++
	}

	sort.Float64s(sr.resultValues[thread])
	return nil
}

func (sr *Runner) search(thread int) error {
	if sr.vverbose {
		log.Printf("  Search Sid ")
	}

	resp, err := sr.client.Post(search, strings.NewReader("search="+sr.query))
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
		log.Printf("%s...\n", s.Sid)
	}

	sr.sids[thread] = s.Sid

	return nil
}

func (sr *Runner) results(thread int) (bool, error) {
	if sr.vverbose {
		log.Printf("    Checking results for %s...", sr.sids[thread])
	}

	resp, err := sr.client.Get(fmt.Sprintf(status, sr.sids[thread]), nil)
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
		log.Printf(" done: %v\n", done)

		if done {
			log.Printf("  Finished Sid %s in %f.\n", sr.sids[thread], s.Entry[0].Content.RunDuration)
		}
	}

	if done {
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

func (sr *Runner) JsonPrint() {
	o := Results{}

	o.Query = sr.query
	o.Average, o.Median, o.Min, o.Max = sr.agg()

	for _, r := range sr.resultValues {
		v := ResultThread{}
		v.Average = math.Avg(r)
		v.Median = math.Med(r)
		v.Min = math.Min(r)
		v.Max = math.Max(r)
		v.Run = r

		o.Thread = append(o.Thread, v)
	}

	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	var j bytes.Buffer
	err = json.Indent(&j, data, "", "  ")
	if err != nil {
		panic(err)
	}

	log.Println(j.String())
}

func (sr *Runner) PrettyPrint() {
	sr.PrintBanner()
	for i := 0; i < sr.threads; i++ {
		sr.PrintResults(i)
	}
	sr.PrintAgg()
	sr.PrintFooter()

	if sr.vverbose {
		log.Println("--------------------------------------------------------------------------------")
		log.Printf("%#v\n---\n\n", sr)
	}
}

func (sr *Runner) PrintBanner() {
	log.Printf("\n %-10s | %-10s | %-10s | %-10s | %-10s | %-10s\n",
		"Thread", "Runs", "Average", "Median", "Min", "Max")
	log.Println("--------------------------------------------------------------------------------")
}

func (sr *Runner) agg() (avg, med, min, max float64) {
	v := []float64{}
	for _, t := range sr.resultValues {
		v = append(v, t...)
	}
	sort.Float64s(v)

	return math.Avg(v), math.Med(v), math.Min(v), math.Max(v)
}

func (sr *Runner) PrintAgg() {
	avg, med, min, max := sr.agg()

	log.Println("--------------------------------------------------------------------------------")
	log.Printf("     -  aggregate  -     | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
		avg, med, min, max)
}

func (sr *Runner) PrintFooter() {
	log.Println("--------------------------------------------------------------------------------")
	log.Printf(" Query: %.70s...\n\n", sr.query)
}

func (sr *Runner) PrintResults(thread int) {
	v := sr.resultValues[thread]
	log.Printf(" %-10d | %-10d | %-10.4f | %-10.4f | %-10.4f | %-10.4f\n",
		thread, sr.ran, math.Avg(v), math.Med(v), math.Min(v), math.Max(v))
}
