package runner

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmervine/splunk-benchmark/lib/util"
)

const (
	SearchEndpoint = "/services/search/jobs"
	StatusEndpoint = SearchEndpoint + "/%s"
)

type Run struct {
	runner *Runner
	Sid    string  `json:"sid"`
	Dur    float64 `json:"duration"`
	Err    error   `json:"-"`
}

func (r Run) Query() string {
	return r.runner.query
}

func NewRun(runner *Runner) Run {
	run := Run{}
	run.runner = runner
	return run
}

func (run *Run) Search() {
	util.Vverbose("func=run.Search at=start")
	defer func() {
		util.Verbosef("search: %s >", run.Sid)
		util.Vverbosef("func=run.Search run:\n %#v\n", run)
		util.Vverbose("func=run.Search at=finish")
	}()

	resp, err := run.runner.client.Post(SearchEndpoint, strings.NewReader("search="+run.Query()))
	if err != nil {
		run.Err = err
		return
	}

	defer resp.Body.Close()

	res := new(struct {
		Sid string `json:"sid"`
	})

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		run.Err = err
		return
	}

	run.Sid = res.Sid
}

func (run *Run) GetResult() {
	util.Vverbose("func=run.GetResult at=start")
	defer func() {
		util.Verbosef("      < %s: %.3f", run.Sid, run.Dur)
		util.Vverbosef("func=run.GetResult run:\n %#v\n", run)
		util.Vverbose("func=run.GetResult at=finish")
	}()

	if run.Sid == "" {
		run.Err = fmt.Errorf("An SID is required to get results.")
		util.Verbosef("Error: %v", run.Err)
		return
	}

	endpoint := fmt.Sprintf(StatusEndpoint, run.Sid)

	res := new(struct {
		Entry []struct {
			Content struct {
				IsDone      bool    `json:"isDone"`
				RunDuration float64 `json:"runDuration"`
			} `json:"content"`
		} `json:"entry"`
	})

	for {
		resp, err := run.runner.client.Get(endpoint, nil)
		if err != nil {
			run.Err = err
			return
		}

		err = json.NewDecoder(resp.Body).Decode(res)
		if err != nil {
			run.Err = err
			return
		}

		c := res.Entry[0].Content
		util.Vverbosef("func=run.GetResult#loop content:\n %#v\n", c)
		if c.IsDone {
			run.Dur = c.RunDuration
			break
		}

		time.Sleep(time.Second)
	}
}
