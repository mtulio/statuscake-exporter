package stk

import (
	"log"
	"net/url"
	"strconv"
	"time"

	statuscake "github.com/mtulio/statuscake-exporter/statuscacke"
)

type StkOptions struct {
	Client statuscake.Client
	Tags   string
}

type StkAPI struct {
	client          *statuscake.Client
	configTags      string
	waitIntervalSec uint8
	Tests           []*statuscake.Test
}

type StkTest statuscake.Test

func NewStkAPI(user string, pass string) (*StkAPI, error) {

	// connect to the StatusCake API
	c, err := statuscake.New(
		statuscake.Auth{
			Username: user,
			Apikey:   pass,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &StkAPI{
		client:          c,
		waitIntervalSec: 30,
	}, nil
}

// get/set
func (stk *StkAPI) SetConfigTags(tags string) {
	stk.configTags = tags
}

func (stk *StkAPI) GetTests() []*statuscake.Test {
	return stk.Tests
}

func (stk *StkAPI) SetWaitInterval(sec uint8) {
	stk.waitIntervalSec = sec
}

// gather functions
func (stk *StkAPI) GatherAll() error {
	go stk.gatherTest()
	go stk.gatherTestsData()
	return nil
}

func (stk *StkAPI) gatherTest() {
	for {
		v := url.Values{}
		if stk.configTags != "" {
			v.Set("tags", stk.configTags)
		}
		tests, err := stk.client.Tests().AllWithFilter(v)
		if err != nil {
			log.Println(err)
		}
		stk.Tests = tests
		time.Sleep(time.Second * time.Duration(stk.waitIntervalSec))
	}
}

func (stk *StkAPI) gatherTestsData() {
	for {
		if len(stk.Tests) <= 0 {
			time.Sleep(10 * time.Duration(stk.waitIntervalSec))
			continue
		}
		filters := url.Values{}

		filters.Set("Fields", "performance,status,location,time")
		filters.Set("Limit", strconv.Itoa(10))

		for t := range stk.Tests {
			filters.Set("TestID", strconv.Itoa(stk.Tests[t].TestID))
			perfData, err := stk.client.PerfData().AllWithFilter(filters)
			if err != nil {
				log.Println(err)
			}
			stk.Tests[t].PerformanceData = perfData
		}

		time.Sleep(time.Second * time.Duration(stk.waitIntervalSec))
	}
}
