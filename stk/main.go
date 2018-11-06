package stk

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/mtulio/statuscake"
)

type StkOptions struct {
	Client statuscake.Client
	Tags   string
}

type StkAPI struct {
	client          *statuscake.Client
	configTags      string
	waitIntervalSec uint8
	ClientTests     []*statuscake.Test
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
	return stk.ClientTests
}

func (stk *StkAPI) SetWaitInterval(sec uint8) {
	stk.waitIntervalSec = sec
}

// gather functions
func (stk *StkAPI) GatherAll() error {
	go stk.gatherTest()
	stk.gatherPerfData()
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

		stk.ClientTests = tests
		time.Sleep(time.Second * time.Duration(stk.waitIntervalSec))
	}
}

func (stk *StkAPI) gatherPerfData() {
	testId := 1314813
	for {
		// v := url.Values{}
		// if stk.configTags != "" {
		// 	v.Set("tags", stk.configTags)
		// }
		perfData, err := stk.client.PerfData().All(testId)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("%v", perfData)

		// stk.ClientTests = tests
		time.Sleep(time.Second * time.Duration(stk.waitIntervalSec))
	}
}
