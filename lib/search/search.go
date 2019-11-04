package search

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

var logger = log.New(os.Stdout, "", 0)

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
	runs         int
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
		logger.Printf("%#v\n---\n\n", sr)
	}

	return sr, err
}

func (sr *Runner) Do(thread int) error {
	logPrefix := fmt.Sprintf("thread=%d fn=Do", thread)

	if sr.verbose || sr.vverbose {
		logger.Printf("%s at=start", logPrefix)
	}

	for {
		ran := len(sr.resultValues[thread])
		run := ran + 1

		if sr.verbose || sr.vverbose {
			logger.Printf("%s run=%d at=start#loop", logPrefix, run)
		}

		err := sr.search(thread)
		if err != nil {
			return err
		}

		for {
			done, err := sr.getResults(thread)
			if err != nil {
				return err
			}

			if done {
				break
			}
		}

		if sr.verbose || sr.vverbose {
			logger.Printf("%s run=%d at=finish#loop", logPrefix, run)
		}

		if sr.runs > 0 && run == sr.runs {
			break
		}
	}

	if sr.verbose || sr.vverbose {
		logger.Printf("%s at=finish", logPrefix)
	}

	sort.Float64s(sr.resultValues[thread])
	return nil
}

func (sr *Runner) search(thread int) error {
	run := len(sr.resultValues[thread]) + 1

	logPrefix := fmt.Sprintf("thread=%d run=%d fn=search", thread, run)

	if sr.vverbose {
		logger.Printf("%s at=start", logPrefix)
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

	sr.sids[thread] = s.Sid

	if sr.vverbose {
		logger.Printf("%s sid=%s at=finish", logPrefix, s.Sid)
	}

	return nil
}

func (sr *Runner) getResults(thread int) (bool, error) {
	// short sleep before checking for results
	time.Sleep(time.Second)

	run := len(sr.resultValues[thread]) + 1

	logPrefix := fmt.Sprintf("thread=%d run=%d fn=getResults sid=%s ", thread, run, sr.sids[thread])

	if sr.vverbose {
		logger.Printf("%s at=start", logPrefix)
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
		if done {
			logger.Printf("%s isDone=%v duration=%f at=finish", logPrefix, done, s.Entry[0].Content.RunDuration)
		} else {
			logger.Printf("%s isDone=%v at=finish", logPrefix, done)
		}
	}

	if done {
		sr.resultValues[thread] = append(sr.resultValues[thread], s.Entry[0].Content.RunDuration)
		return true, nil
	}

	return false, nil
}

func (sr *Runner) Results() Results {
	v := []float64{}
	for _, t := range sr.resultValues {
		v = append(v, t...)
	}
	sort.Float64s(v)

	o := Results{
		Query:   sr.query,
		Average: math.Avg(v),
		Median:  math.Med(v),
		Min:     math.Min(v),
		Max:     math.Max(v),
		Thread:  []ResultThread{},
	}

	for _, r := range sr.resultValues {
		sort.Float64s(r)

		o.Thread = append(o.Thread, ResultThread{
			Average: math.Avg(r),
			Median:  math.Med(r),
			Min:     math.Min(r),
			Max:     math.Max(r),
			Run:     r,
		})
	}

	return o
}
