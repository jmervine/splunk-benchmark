package search

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

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

type SplunkRun struct {
	client   splunking.SplunkRequest
	verbose  bool
	vverbose bool

	Query        string
	Runs         int
	Delay        time.Duration
	Sid          string
	Results      map[string]float64
	resultValues []float64
}

func NewRun(host, query string, runs int, delay float64, verbose, vverbose bool) (*SplunkRun, error) {
	var err error

	sr := new(SplunkRun)

	sr.Query = query
	sr.Runs = runs
	sr.Delay = time.Duration(delay) * time.Second
	sr.Results = make(map[string]float64)
	sr.verbose = verbose
	sr.vverbose = vverbose

	sr.client, err = splunking.InitURL(host)

	if vverbose {
		fmt.Printf("%#v\n---\n\n", sr)
	}

	return sr, err
}

func (sr *SplunkRun) Do() error {
	for i := 0; i < sr.Runs; i++ {
		if sr.verbose || sr.vverbose {
			fmt.Printf("Starting run %d...\n", i+1)
		}

		err := sr.search()
		if err != nil {
			return err
		}

		for {
			done, err := sr.results()
			if err != nil {
				return err
			}

			if done {
				break
			}
		}

		if sr.verbose || sr.vverbose {
			fmt.Printf("Finished run %d.\n", i+1)
		}
	}

	sort.Float64s(sr.resultValues)
	return nil
}

func (sr *SplunkRun) search() error {
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

	sr.Sid = s.Sid

	return nil
}

func (sr *SplunkRun) results() (bool, error) {
	if sr.vverbose {
		fmt.Printf("    Checking results for %s...", sr.Sid)
	}

	resp, err := sr.client.Get(fmt.Sprintf(status, sr.Sid), nil)
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
			fmt.Printf("  Finished Sid %s in %f.\n", sr.Sid, s.Entry[0].Content.RunDuration)
		}
	}

	if done {
		sr.Results[sr.Sid] = s.Entry[0].Content.RunDuration
		sr.resultValues = append(sr.resultValues, s.Entry[0].Content.RunDuration)
		return true, nil
	}

	return false, nil
}

func (sr *SplunkRun) Avg() float64 {
	t := float64(0)
	for _, n := range sr.resultValues {
		t = t + n
	}

	return t / float64(len(sr.resultValues))
}

func (sr *SplunkRun) Med() float64 {
	v := sr.resultValues

	nn := len(v) / 2
	if nn%2 != 0 {
		return v[nn]
	}

	return (v[nn+1] + v[nn]) / 2
}

func (sr *SplunkRun) Min() float64 {
	return sr.resultValues[0]
}

func (sr *SplunkRun) Max() float64 {
	return sr.resultValues[len(sr.resultValues)-1]
}

func (sr *SplunkRun) PrettyPrint() {
	fmt.Printf("\n %-12s | %-12s | %-12s | %-12s | %-12s\n",
		"Runs", "Average", "Median", "Min", "Max")
	fmt.Println("--------------------------------------------------------------------------------")
	sr.PrintResults()
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf(" Query: %.70s...\n\n", sr.Query)
}

func (sr *SplunkRun) PrintResults() {
	fmt.Printf(" %-12d | %-12.4f | %-12.4f | %-12.4f | %-12.4f\n",
		sr.Runs, sr.Avg(), sr.Med(), sr.Min(), sr.Max())
}
