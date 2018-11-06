package stk

import (
	"log"
	"net/url"
	"time"

	"github.com/DreamItGetIT/statuscake"
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
