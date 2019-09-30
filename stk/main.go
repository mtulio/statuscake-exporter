package stk

import (
	"log"
	"net/url"
	"strconv"
	"time"

	statuscake "github.com/mtulio/statuscake"
)

// StkOptions StatusCake CLI Options
type StkOptions struct {
	Client statuscake.Client
	Tags   string
}

// StkAPI handle the basic API config and last data.
type StkAPI struct {
	client          *statuscake.Client
	configTags      string
	waitIntervalSec uint32
	EnableTests     bool
	EnableTestsSSL  bool
	Tests           []*statuscake.Test
	TestsSSL        []*statuscake.Ssl
	controlInit     bool
	sslFlagsEnabled map[string]bool
}

// type StkTest statuscake.Test

// NewStkAPI create API instance to communicate with StatusCake API.
func NewStkAPI(user string, pass string) (*StkAPI, error) {

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
		EnableTests:     false,
		EnableTestsSSL:  false,
	}, nil
}

// SetConfigTags define Tags on configuration.
func (stk *StkAPI) SetConfigTags(tags string) {
	stk.configTags = tags
}

// GetTags return all tags defined.
func (stk *StkAPI) GetTags() string {
	return stk.configTags
}

// GetTests return all StatusCake Tests from last API discovery.
func (stk *StkAPI) GetTests() []*statuscake.Test {
	return stk.Tests
}

// GetTestsSSL return all StatusCake SSL Tests from last API discovery.
func (stk *StkAPI) GetTestsSSL() []*statuscake.Ssl {
	return stk.TestsSSL
}

// SetWaitInterval define API data scrape wait internval.
func (stk *StkAPI) SetWaitInterval(sec uint32) {
	stk.waitIntervalSec = sec
}

// GetWaitInterval return the wait interval value.
func (stk *StkAPI) GetWaitInterval() uint32 {
	return stk.waitIntervalSec
}

// SetEnableTests set Web Test flag to gather Tests.
func (stk *StkAPI) SetEnableTests(v bool) {
	stk.EnableTests = v
}

// GetEnableTests return the Flag to gather Web Tests.
func (stk *StkAPI) GetEnableTests() bool {
	return stk.EnableTests
}

// SetEnableSSL set SSL flag to gather Tests.
func (stk *StkAPI) SetEnableSSL(v bool) {
	stk.EnableTestsSSL = v
}

// GetEnableSSL return the SSL Flag to gather SSL Tests.
func (stk *StkAPI) GetEnableSSL() bool {
	return stk.EnableTestsSSL
}

// GatherAll retrieves all data for enabled modules.
func (stk *StkAPI) GatherAll() error {
	if stk.EnableTests {
		go stk.gatherTest()
		go stk.gatherTestsData()
	}
	if stk.EnableTestsSSL {
		go stk.gatherTestsSSL()
	}
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
			test := stk.Tests[t]
			filters.Set("TestID", strconv.Itoa(test.TestID))
			perfData, err := stk.client.PerfData().AllWithFilter(filters)
			if err != nil {
				log.Println(err)
			}
			test.PerformanceData = perfData
		}
		time.Sleep(time.Second * time.Duration(stk.waitIntervalSec))
	}
}

func (stk *StkAPI) gatherTestsSSL() {
	sslCli := statuscake.NewSsls(stk.client)
	for {
		ssls, err := sslCli.All()
		if err != nil {
			log.Println(err)
		}
		stk.TestsSSL = ssls
		if !stk.controlInit {
			log.Println(" Initial API discovery returns the Total of SSL Tests:", len(stk.TestsSSL))
		}
		time.Sleep(time.Second * time.Duration(stk.waitIntervalSec))
	}
}

// CheckSSLFlagIsEnabled check if SSL flag is enabled.
func (stk *StkAPI) CheckSSLFlagIsEnabled(fname string) bool {
	if ok := stk.sslFlagsEnabled[fname]; ok {
		return ok
	}
	if ok := stk.sslFlagsEnabled["all"]; ok {
		return ok
	}
	return false
}

// SetSSLFlag Set SSL Flag
func (stk *StkAPI) SetSSLFlag(fname string) {
	if len(stk.sslFlagsEnabled) == 0 {
		stk.sslFlagsEnabled = make(map[string]bool)
	}
	stk.sslFlagsEnabled[fname] = true
}

// GetSSLFlags Set SSL Flag
func (stk *StkAPI) GetSSLFlags() map[string]bool {
	return stk.sslFlagsEnabled
}
