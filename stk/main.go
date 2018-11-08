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
	waitIntervalSec uint32
	Tests           []*statuscake.Test
	controlInit     bool
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
		waitIntervalSec: 300,
		controlInit:     false,
	}, nil
}

// get/set
func (stk *StkAPI) SetConfigTags(tags string) {
	stk.configTags = tags
}

func (stk *StkAPI) GetTags() string {
	return stk.configTags
}

func (stk *StkAPI) GetTests() []*statuscake.Test {
	return stk.Tests
}

func (stk *StkAPI) SetWaitInterval(sec uint32) {
	stk.waitIntervalSec = sec
}

func (stk *StkAPI) GetWaitInterval() uint32 {
	return stk.waitIntervalSec
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
		if !stk.controlInit {
			log.Println(" Initial API discovery returns the Total of Tests:", len(tests))
			stk.controlInit = true
		}
		time.Sleep(time.Second * time.Duration(stk.waitIntervalSec))
	}
}

func (stk *StkAPI) gatherTestsData() {
	for {
		if len(stk.Tests) <= 0 {
			time.Sleep(time.Second * 10)
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
